package app

import (
	"context"

	"golang.org/x/sync/errgroup"

	"github/alexveli1/packageanalyzer/internal/config"
	"github/alexveli1/packageanalyzer/internal/domain"
	"github/alexveli1/packageanalyzer/internal/repository"
	"github/alexveli1/packageanalyzer/internal/service"
	"github/alexveli1/packageanalyzer/internal/transport/httpv1"
	"github/alexveli1/packageanalyzer/internal/usecase"
	"github/alexveli1/packageanalyzer/pkg/mylog"
)

// Run - application starter
// initializes configuration, logger, repositories, HTTP client, service layer, and usecase object
// sends high level instructions to Usecase layer
func Run() {
	mylog.SugarLogger = mylog.InitLogger(domain.LogTypeStdOut, domain.LogFileName)
	newConfig := config.NewConfig()
	mylog.SugarLogger.Infof("starting with config %v", newConfig)
	newRepositories := repository.NewRepositories()
	newClient := httpv1.NewClient(newConfig)
	newServices := service.NewServices(newRepositories, newClient, newConfig)
	newUsecase := usecase.NewUsecase(newServices)

	grp, ctx1 := errgroup.WithContext(context.Background())
	grp.Go(func() error {
		return newUsecase.GetPackages(ctx1, domain.Sisyphus)
	})
	grp.Go(func() error {
		return newUsecase.GetPackages(ctx1, domain.P10)
	})
	if err := grp.Wait(); err != nil {
		return
	}
	ctx := context.Background()
	if newConfig.Scope != domain.ScopeReleases {
		newUsecase.GetUniquePackages(ctx, domain.Sisyphus, domain.P10)
		newUsecase.GetUniquePackages(ctx, domain.P10, domain.Sisyphus)
	}
	if newConfig.Scope != domain.ScopeDiff {
		newUsecase.GetHigherReleases(ctx, domain.Sisyphus, domain.P10)
		newUsecase.GetHigherReleases(ctx, domain.P10, domain.Sisyphus)
	}
	if newConfig.VerifyResult {
		grp1, ctx2 := errgroup.WithContext(context.Background())
		grp1.Go(func() error {
			return newUsecase.GetVerificationData(ctx2, domain.Sisyphus, domain.P10)
		})
		grp1.Go(func() error {
			return newUsecase.GetVerificationData(ctx, domain.P10, domain.Sisyphus)
		})
		if err := grp1.Wait(); err != nil {
			return
		}
		grp2, ctx3 := errgroup.WithContext(context.Background())
		grp2.Go(func() error {
			return newUsecase.VerifyResult(ctx3, domain.Sisyphus)
		})
		grp2.Go(func() error {
			return newUsecase.VerifyResult(ctx, domain.P10)
		})
		if err := grp2.Wait(); err != nil {
			return
		}
	}
	newUsecase.PrintResult(ctx)
}

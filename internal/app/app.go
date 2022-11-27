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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	mylog.SugarLogger = mylog.InitLogger(domain.LogTypeStdOut, domain.LogFileName)
	newConfig := config.NewConfig()
	mylog.SugarLogger.Infof("starting with config %v", newConfig)
	newRepositories := repository.NewRepositories()
	newClient := httpv1.NewClient(newConfig)
	newServices := service.NewServices(newRepositories, newClient, newConfig)
	newUsecase := usecase.NewUsecase(newServices)

	grp, ctx1 := errgroup.WithContext(ctx)
	grp.Go(func() error {
		return newUsecase.GetPackages(ctx1, domain.Sisyphus)
	})
	grp.Go(func() error {
		return newUsecase.GetPackages(ctx1, domain.P10)
	})
	if err := grp.Wait(); err != nil {
		return
	}
	ctx2 := context.Background()
	if newConfig.Scope != domain.ScopeReleases {
		grp1, c2 := errgroup.WithContext(ctx2)
		grp1.Go(func() error {
			return newUsecase.GetUniquePackages(c2, domain.Sisyphus, domain.P10)
		})
		grp1.Go(func() error {
			return newUsecase.GetUniquePackages(c2, domain.P10, domain.Sisyphus)
		})
		if err := grp1.Wait(); err != nil {
			return
		}
	}
	ctx3 := context.Background()
	if newConfig.Scope != domain.ScopeDiff {
		grp2, c3 := errgroup.WithContext(ctx3)
		grp2.Go(func() error {
			return newUsecase.GetHigherReleases(c3, domain.Sisyphus, domain.P10)
		})
		grp2.Go(func() error {
			return newUsecase.GetHigherReleases(c3, domain.P10, domain.Sisyphus)
		})
		if err := grp2.Wait(); err != nil {
			return
		}
	}
	newUsecase.PrintResult(ctx)
}

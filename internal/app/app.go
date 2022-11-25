package app

import (
	"context"

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

	newUsecase.GetPackages(ctx, domain.Sisyphus)
	newUsecase.GetPackages(ctx, domain.P10)

	if newConfig.Scope != domain.ScopeReleases {
		newUsecase.GetUniquePackages(ctx, domain.Sisyphus, domain.P10)
		newUsecase.GetUniquePackages(ctx, domain.P10, domain.Sisyphus)
	}
	if newConfig.Scope != domain.ScopeDiff {
		newUsecase.GetHigherReleases(ctx, domain.Sisyphus, domain.P10)
		newUsecase.GetHigherReleases(ctx, domain.P10, domain.Sisyphus)
	}
	newUsecase.PrintResult(ctx)
}

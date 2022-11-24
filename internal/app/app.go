package app

import (
	"github/alexveli1/packageanalyzer/internal/config"
	"github/alexveli1/packageanalyzer/internal/domain"
	"github/alexveli1/packageanalyzer/internal/repository"
	"github/alexveli1/packageanalyzer/internal/service"
	"github/alexveli1/packageanalyzer/internal/transport/httpv1"
	"github/alexveli1/packageanalyzer/internal/usecase"
	"github/alexveli1/packageanalyzer/pkg/mylog"
)

func Run() {
	mylog.SugarLogger = mylog.InitLogger(domain.LogTypeStdOut, domain.LogFileName)
	newConfig := config.NewConfig()
	newRepositories := repository.NewRepositories()
	newClient := httpv1.NewClient(newConfig)
	newServices := service.NewServices(newRepositories, newClient, newConfig)
	newUsecase := usecase.NewUsecase(newServices)
	newUsecase.GetXOR()
}

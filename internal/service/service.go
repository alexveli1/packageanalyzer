package service

import (
	"github/alexveli1/packageanalyzer/internal/config"
	"github/alexveli1/packageanalyzer/internal/domain"
	"github/alexveli1/packageanalyzer/internal/repository"
	"github/alexveli1/packageanalyzer/internal/transport/httpv1"
)

type Analyzer interface {
	CompareXOR(b1, b2 string) ([]domain.Binpack, error)
	ExistsInBoth(b1, b2 string) ([]domain.Binpack, error)
	CompareReleases([]domain.Binpack) ([]domain.Comparepack, error)
}

type Services struct {
	Analyzer
}

func NewServices(repo *repository.Repositories, transporter httpv1.ITransporter, cfg *config.Config) *Services {
	return &Services{
		Analyzer: NewAnalyzerService(repo, &transporter, cfg),
	}
}

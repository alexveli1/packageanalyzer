package service

import (
	"github/alexveli1/packageanalyzer/internal/config"
	"github/alexveli1/packageanalyzer/internal/domain"
	"github/alexveli1/packageanalyzer/internal/repository"
	"github/alexveli1/packageanalyzer/internal/transport/httpv1"
)

type AnalyzerService struct {
	repo   *repository.Repositories
	client *httpv1.ITransporter
}

func NewAnalyzerService(repo *repository.Repositories, transporter *httpv1.ITransporter, cfg *config.Config) *AnalyzerService {

	return &AnalyzerService{
		repo:   repo,
		client: transporter,
	}
}

func (as *AnalyzerService) CompareXOR(b1, b2 string) ([]domain.Binpack, error) {
	return nil, nil
}
func (as *AnalyzerService) ExistsInBoth(b1, b2 string) ([]domain.Binpack, error) {
	return nil, nil
}
func (as *AnalyzerService) CompareReleases([]domain.Binpack) ([]domain.Comparepack, error) {
	return nil, nil
}

package service

import (
	"context"

	"github/alexveli1/packageanalyzer/internal/config"
	"github/alexveli1/packageanalyzer/internal/domain"
	"github/alexveli1/packageanalyzer/internal/repository"
	"github/alexveli1/packageanalyzer/internal/transport/httpv1"
)

type Analyzer interface {
	PackagesFromBranch1(ctx context.Context, branch1 string, branch2 string) (map[string][]domain.Binpack, map[string][]domain.Binpack)
	Branch1Higher(ctx context.Context, branch1 string, branch2 string) map[string][]domain.Binpack
	GetPacks(ctx context.Context, branch string) error
}

type Services struct {
	Analyzer
}

func NewServices(repo *repository.Repositories, transporter httpv1.ITransporter, cfg *config.Config) *Services {
	return &Services{
		Analyzer: NewAnalyzerService(repo, transporter, cfg),
	}
}

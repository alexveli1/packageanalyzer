// Package service requests clients to provide data, requests repo to save data,
// processes data and provides results to usecase layer
package service

import (
	"context"

	"github/alexveli1/packageanalyzer/internal/config"
	"github/alexveli1/packageanalyzer/internal/domain"
	"github/alexveli1/packageanalyzer/internal/repository"
	"github/alexveli1/packageanalyzer/internal/transport/httpv1"
)

// Analyzer used for mocking service layer logic for tests
type Analyzer interface {
	GetUnique(ctx context.Context, branch1 string, branch2 string) (domain.Result, error)
	GetHigher(ctx context.Context, branch1 string, branch2 string) (domain.Result, error)
	GetPacks(ctx context.Context, branch string) error
}

// Services consolidated object for hosting any current or future services
type Services struct {
	Analyzer
}

func NewServices(repo *repository.Repositories, transporter httpv1.ITransporter, cfg *config.Config) *Services {
	return &Services{
		Analyzer: NewAnalyzerService(repo, transporter, cfg),
	}
}

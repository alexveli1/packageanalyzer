package repository

import (
	"context"

	"github/alexveli1/packageanalyzer/internal/domain"
)

type Analyzer interface {
	SavePacks(ctx context.Context, branch string, packs map[string][]domain.Binpack) error
	GetAllPacks(ctx context.Context, branch string) (map[string][]domain.Binpack, error)
	GetPackByName(ctx context.Context, branch string, packName string) ([]domain.Binpack, bool)
}

type Repositories struct {
	Analyzer
}

func NewRepositories() *Repositories {
	return &Repositories{
		Analyzer: NewAnalyzerRepo(),
	}
}

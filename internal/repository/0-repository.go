package repository

import (
	"context"

	"github/alexveli1/packageanalyzer/internal/domain"
)

type Analyzer interface {
	SaveAll(ctx context.Context, packs []domain.Binpack) error
	GetPacksByBranch(ctx context.Context) ([]domain.Binpack, error)
	GetPackByName(ctx context.Context, pack string) (domain.Binpack, error)
}

type Repositories struct {
	Sisyphus Analyzer
	P10      Analyzer
}

func NewRepositories() *Repositories {
	return &Repositories{
		Sisyphus: NewSisyphusRepo(),
		P10:      NewP10Repo(),
	}
}

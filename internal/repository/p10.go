package repository

import (
	"context"

	"github/alexveli1/packageanalyzer/internal/domain"
	"github/alexveli1/packageanalyzer/pkg/storage"
)

type P10Repo struct {
	branch *storage.MapDB
}

func NewP10Repo() *P10Repo {
	return &P10Repo{branch: storage.NewMap()}
}
func (p *P10Repo) SaveAll(ctx context.Context, packs []domain.Binpack) error {
	return nil
}
func (p *P10Repo) GetPacksByBranch(ctx context.Context) ([]domain.Binpack, error) {
	return nil, nil
}
func (p *P10Repo) GetPackByName(ctx context.Context, pack string) (domain.Binpack, error) {
	return domain.Binpack{}, nil
}

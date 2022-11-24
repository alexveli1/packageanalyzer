package repository

import (
	"context"

	"github/alexveli1/packageanalyzer/internal/domain"
	"github/alexveli1/packageanalyzer/pkg/storage"
)

type SisyphusRepo struct {
	branch *storage.MapDB
}

func NewSisyphusRepo() *SisyphusRepo {
	return &SisyphusRepo{
		branch: storage.NewMap(),
	}
}

func (s *SisyphusRepo) SaveAll(ctx context.Context, packs []domain.Binpack) error {
	return nil
}
func (s *SisyphusRepo) GetPacksByBranch(ctx context.Context) ([]domain.Binpack, error) {
	return nil, nil
}
func (s *SisyphusRepo) GetPackByName(ctx context.Context, pack string) (domain.Binpack, error) {
	return domain.Binpack{}, nil
}

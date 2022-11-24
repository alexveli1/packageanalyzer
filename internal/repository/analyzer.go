package repository

import (
	"context"

	"github/alexveli1/packageanalyzer/internal/domain"
	"github/alexveli1/packageanalyzer/pkg/mylog"
	"github/alexveli1/packageanalyzer/pkg/storage"
)

type AnalyzerRepo struct {
	sisyphus *storage.MapDB
	p10      *storage.MapDB
}

func NewAnalyzerRepo() *AnalyzerRepo {
	return &AnalyzerRepo{
		sisyphus: storage.NewMap(),
		p10:      storage.NewMap(),
	}
}

func (s *AnalyzerRepo) SavePacks(ctx context.Context, branch string, packs map[string][]domain.Binpack) error {
	if err := ctx.Err(); err != nil {

		return err
	}
	b := s.setBranch(branch)
	b.Lock.Lock()
	defer b.Lock.Unlock()
	b.DB = packs

	return nil
}
func (s *AnalyzerRepo) GetAllPacks(ctx context.Context, branch string) (map[string][]domain.Binpack, error) {
	if err := ctx.Err(); err != nil {

		return nil, err
	}
	b := s.setBranch(branch)
	b.Lock.Lock()
	defer b.Lock.Unlock()
	return b.DB, nil
}
func (s *AnalyzerRepo) GetPackByArchAndName(ctx context.Context, branch string, arch string, packName string) (domain.Binpack, bool) {
	if err := ctx.Err(); err != nil {

		return domain.Binpack{}, false
	}
	b := s.setBranch(branch)
	b.Lock.Lock()
	defer b.Lock.Unlock()
	v, ok := b.DB[packName]
	if ok {
		for i := 0; i < len(v); i++ {
			if v[i].Arch == arch {

				return v[i], true
			}
		}
	}
	return domain.Binpack{}, false
}

func (s *AnalyzerRepo) setBranch(branch string) *storage.MapDB {
	switch branch {
	case domain.P10:

		return s.p10
	case domain.Sisyphus:

		return s.sisyphus
	default:

		mylog.SugarLogger.Warnf("incorrect branch name %s", branch)
	}
	return nil
}

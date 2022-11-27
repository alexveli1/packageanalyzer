package repository

import (
	"context"

	"github/alexveli1/packageanalyzer/internal/domain"
	"github/alexveli1/packageanalyzer/pkg/mylog"
	"github/alexveli1/packageanalyzer/pkg/storage"
)

// AnalyzerRepo implementation of Analyzer repo for mocking repo layer or potential using several data storage options
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

// SavePacks stores data got be service layer from HTTP client
func (s *AnalyzerRepo) SavePacks(ctx context.Context, branch string, packs domain.Branch) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	b := s.setBranch(branch)
	b.Lock.Lock()
	defer b.Lock.Unlock()
	b.DB = packs

	return nil
}

// GetBranchPacks returns full list of packages for single branch
func (s *AnalyzerRepo) GetBranchPacks(ctx context.Context, branch string) (domain.Branch, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	b := s.setBranch(branch)
	b.Lock.Lock()
	defer b.Lock.Unlock()

	return b.DB, nil
}

// GetPacksByArchAndName provides slice of domain.Binpack
// considering there might be several packages with same architecure,
// but different releases and versions in ALT repository
func (s *AnalyzerRepo) GetPacksByArchAndName(ctx context.Context, branch string, arch string, packName string) ([]domain.Binpack, bool, error) {
	if err := ctx.Err(); err != nil {
		return nil, false, err
	}
	packs := make([]domain.Binpack, 0)
	b := s.setBranch(branch)
	b.Lock.Lock()
	defer b.Lock.Unlock()
	v, ok := b.DB[packName]
	if !ok {
		return nil, false, nil
	}
	for i := 0; i < len(v); i++ {
		if v[i].Arch == arch {
			packs = append(packs, v[i]) // in case we might have several packages with same name for same arch but different versions
		}
	}
	if len(packs) > 0 {
		return packs, true, nil
	}

	return nil, false, nil
}

// setBranch selects repo branch for operations
// implemented to avoid duplicated functions/repos and simplification of data storage
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

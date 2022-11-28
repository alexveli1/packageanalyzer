package repository

import (
	"context"
	"sync"

	"github/alexveli1/packageanalyzer/internal/domain"
	"github/alexveli1/packageanalyzer/pkg/mylog"
)

type BranchStore struct {
	DB   domain.Branch
	Lock *sync.RWMutex
}
type CompareStore struct {
	DB   domain.CompareBranch
	Lock *sync.RWMutex
}
type SourceStore struct {
	DB   domain.BranchSources
	Lock *sync.RWMutex
}

// AnalyzerRepo implementation of Analyzer repo for mocking repo layer or potential using several data storage options
type AnalyzerRepo struct {
	sisyphus    *BranchStore
	p10         *BranchStore
	compare     *CompareStore
	sourceStore *SourceStore
}

func NewBranchStore() *BranchStore {
	return &BranchStore{
		DB:   make(domain.Branch),
		Lock: &sync.RWMutex{},
	}
}
func NewCompareStore() *CompareStore {
	return &CompareStore{
		DB:   make(domain.CompareBranch),
		Lock: &sync.RWMutex{},
	}
}
func NewSourceStore() *SourceStore {
	return &SourceStore{
		DB:   make(domain.BranchSources),
		Lock: &sync.RWMutex{},
	}
}

func NewAnalyzerRepo() *AnalyzerRepo {
	return &AnalyzerRepo{
		sisyphus:    NewBranchStore(),
		p10:         NewBranchStore(),
		compare:     NewCompareStore(),
		sourceStore: NewSourceStore(),
	}
}

// SavePacks stores data got be service layer from HTTP client
func (r *AnalyzerRepo) SavePacks(ctx context.Context, branch string, packs domain.Branch, sources domain.Sources) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	b := r.setBranch(branch)
	b.Lock.Lock()
	defer b.Lock.Unlock()
	b.DB = packs
	r.sourceStore.DB[branch] = sources

	return nil
}

func (r *AnalyzerRepo) SaveComparison(ctx context.Context, verResults *domain.CompareBranch) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	r.compare.Lock.Lock()
	defer r.compare.Lock.Unlock()
	for method, branches := range *verResults {
		for branch, sources := range branches {
			if r.compare.DB[method] == nil {
				b := make(domain.Compare)
				b[branch] = sources
				r.compare.DB[method] = b
			} else {
				r.compare.DB[method][branch] = sources
			}
		}
	}

	return nil
}
func (r *AnalyzerRepo) GetMethodComparison(ctx context.Context, branch, method string) ([]string, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	r.compare.Lock.Lock()
	defer r.compare.Lock.Unlock()

	return r.compare.DB[method][branch], nil
}

// GetBranchPacks returns full list of packages for single branch
func (r *AnalyzerRepo) GetBranchPacks(ctx context.Context, branch string) (domain.Branch, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	b := r.setBranch(branch)
	b.Lock.RLock()
	defer b.Lock.RUnlock()

	return b.DB, nil
}

func (r *AnalyzerRepo) GetSource(ctx context.Context, branch string, pkg domain.Binpack) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}
	r.sourceStore.Lock.RLock()
	defer r.sourceStore.Lock.RUnlock()
	_, ok := r.sourceStore.DB[branch][pkg.Source]
	if !ok {
		return false, nil
	}

	return true, nil
}

// PacksByArchAndNameExist provides slice of domain.Binpack
// considering there might be several packages with same architecure,
// but different releases and versions in ALT repository
func (r *AnalyzerRepo) PacksByArchAndNameExist(ctx context.Context, branch string, pkg domain.Binpack) ([]domain.Binpack, bool, error) {
	if err := ctx.Err(); err != nil {
		return nil, false, err
	}
	packs := make([]domain.Binpack, 0)
	b := r.setBranch(branch)
	b.Lock.RLock()
	defer b.Lock.RUnlock()
	v, ok := b.DB[pkg.Name]
	if !ok {
		return nil, true, nil
	}
	for i := 0; i < len(v); i++ {
		if v[i].Arch == pkg.Arch && v[i].Source == pkg.Source {
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
func (r *AnalyzerRepo) setBranch(branch string) *BranchStore {
	switch branch {
	case domain.P10:

		return r.p10
	case domain.Sisyphus:

		return r.sisyphus
	default:

		mylog.SugarLogger.Warnf("incorrect branch name %r", branch)
	}

	return nil
}

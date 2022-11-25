package service

import (
	"context"
	"fmt"
	"sort"

	version "github.com/knqyf263/go-rpm-version"

	"github/alexveli1/packageanalyzer/internal/config"
	"github/alexveli1/packageanalyzer/internal/domain"
	"github/alexveli1/packageanalyzer/internal/repository"
	"github/alexveli1/packageanalyzer/internal/transport/httpv1"
)

type AnalyzerService struct {
	repo   *repository.Repositories
	client httpv1.ITransporter
}

func NewAnalyzerService(repo *repository.Repositories, transporter httpv1.ITransporter, cfg *config.Config) *AnalyzerService {
	return &AnalyzerService{
		repo:   repo,
		client: transporter,
	}
}

func (as *AnalyzerService) GetPacks(ctx context.Context, branch string) error {
	p, err := as.client.GetRepo(ctx, branch)
	if err != nil {
		return err
	}
	sort.Slice(p, func(i, j int) bool {
		return p[i].Name < p[j].Name && p[i].Version < p[j].Version
	})
	packs := make(map[string][]domain.Binpack)
	for i := 0; i < len(p); i++ {
		packs[p[i].Name] = append(packs[p[i].Name], p[i])
	}
	err = as.repo.SavePacks(ctx, branch, packs)
	if err != nil {
		return err
	}
	return nil
}
func (as *AnalyzerService) PackagesFromBranch1(ctx context.Context, branch1 string, branch2 string) (map[string][]domain.Binpack, map[string][]domain.Binpack) {
	if err := ctx.Err(); err != nil {
		return nil, nil
	}
	chB1 := make(chan map[string][]domain.Binpack)
	chB2 := make(chan map[string][]domain.Binpack)
	go as.uniquePacks(ctx, chB1, branch1, branch2)
	go as.uniquePacks(ctx, chB2, branch2, branch1)

	return <-chB1, <-chB2
}
func (as *AnalyzerService) Branch1Higher(ctx context.Context, branch1 string, branch2 string) map[string][]domain.Binpack {
	if err := ctx.Err(); err != nil {
		return nil
	}
	packs, _ := as.repo.GetAllPacks(ctx, branch1)
	branch1Higher := make(map[string][]domain.Binpack, 0)
	for pkgName, p1archs := range packs {
		for i := 0; i < len(p1archs); i++ {
			pkgBranch2, exists := as.repo.GetPackByArchAndName(ctx, branch2, p1archs[i].Arch, pkgName)
			if exists {
				p1VersionHigher(p1archs[i], pkgBranch2)
				branch1Higher[p1archs[i].Arch] = append(branch1Higher[p1archs[i].Arch], p1archs[i])
			}
		}
	}
	return branch1Higher
}
func (as *AnalyzerService) uniquePacks(ctx context.Context, ch chan map[string][]domain.Binpack, branch1, branch2 string) {
	if err := ctx.Err(); err != nil {
		return
	}
	packs, _ := as.repo.GetAllPacks(ctx, branch1)
	only := make(map[string][]domain.Binpack, 0)
	for k, v := range packs {
		for i := 0; i < len(v); i++ {
			if _, exists := as.repo.GetPackByArchAndName(ctx, branch2, v[i].Arch, k); !exists {
				only[v[i].Arch] = append(only[v[i].Arch], v[i])
			}
		}
	}
	ch <- only
}

func p1VersionHigher(ver1, ver2 domain.Binpack) bool {
	v1 := version.NewVersion(fmt.Sprint(ver1.Epoch) + ":" + ver1.Release + "-" + ver1.Version)
	v2 := version.NewVersion(fmt.Sprint(ver2.Epoch) + ":" + ver2.Release + "-" + ver2.Version)

	if v1.GreaterThan(v2) {
		return true
	}

	return false
}

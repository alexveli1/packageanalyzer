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

// AnalyzerService combines access to repositories and access to HTTP client
// provides results to usecase layer for presentation (e.g printing to stdout or files)
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

// GetPacks triggers request to HTTP Client and saves received data into repository
func (as *AnalyzerService) GetPacks(ctx context.Context, branch string) error {
	p, err := as.client.GetRepo(ctx, branch)
	if err != nil {
		return nil
	}
	sort.Slice(p, func(i, j int) bool {
		return p[i].Name < p[j].Name && p[i].Version < p[j].Version
	})
	packs := make(domain.Branch)
	for i := 0; i < len(p); i++ {
		packs[p[i].Name] = append(packs[p[i].Name], p[i])
	}
	err = as.repo.SavePacks(ctx, branch, packs)
	if err != nil {
		return nil
	}

	return nil
}

// GetUnique finds whether package existing in branch1 is missing in branch2 and adds it to return result
// only package name is taken into account - rpm compare is not used (epochs, releases, versions)
func (as *AnalyzerService) GetUnique(ctx context.Context, branch1 string, branch2 string) (domain.Result, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	packs, err := as.repo.GetBranchPacks(ctx, branch1)
	if err != nil {
		return nil, err
	}
	tmp, only := make(map[string]string), make(domain.Branch)
	for k, v := range packs {
		for i := 0; i < len(v); i++ {
			_, exists, err := as.repo.GetPacksByArchAndName(ctx, branch2, v[i].Arch, k)
			if err != nil {
				return nil, err
			}
			if !exists {
				if _, ok := tmp[v[i].Name]; !ok {
					only[v[i].Arch] = append(only[v[i].Arch], v[i])
					tmp[v[i].Name] = v[i].Name
				}
			}
		}
	}

	return convertToResult(only, branch1, domain.MethodUnique), nil
}

// GetHigher compares package versions in branch1 and branch2 and returns result to usecase layer
// in case no package exists in branch2 package in branch1 considered to have higher version and added into resulting set
func (as *AnalyzerService) GetHigher(ctx context.Context, branch1 string, branch2 string) (domain.Result, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	packs, err := as.repo.GetBranchPacks(ctx, branch1)
	if err != nil {
		return nil, err
	}
	r := make(domain.Branch)
	for pName, pb1 := range packs {
		for i := 0; i < len(pb1); i++ {
			p2, exists, err := as.repo.GetPacksByArchAndName(ctx, branch2, pb1[i].Arch, pName)
			if err != nil {
				return nil, err
			}
			if exists {
				var higher = false
				for j := 0; j < len(p2); j++ {
					if p1VersionHigher(pb1[i], p2[j]) {
						higher = true
					}
					if higher {
						r[pb1[i].Arch] = append(r[pb1[i].Arch], pb1[i])
					}
				}
			} else {
				r[pb1[i].Arch] = append(r[pb1[i].Arch], pb1[i])
			}
		}
	}

	return convertToResult(r, branch1, domain.MethodHigher), nil
}

// p1VersionHigher compares ALT package versions using rpm logic
func p1VersionHigher(ver1, ver2 domain.Binpack) bool {
	v1 := version.NewVersion(fmt.Sprint(ver1.Epoch) + ":" + ver1.Release + "-" + ver1.Version)
	v2 := version.NewVersion(fmt.Sprint(ver2.Epoch) + ":" + ver2.Release + "-" + ver2.Version)

	return v1.GreaterThan(v2)
}

// convertToResult converts domain.Branch structure into domain.Result for returning to usecase layer
func convertToResult(archPkgs domain.Branch, branchName, methodName string) domain.Result {
	b := make(domain.Branch)
	m := make(domain.Method)
	r := make(domain.Result)
	for k, v := range archPkgs {
		b[branchName] = v
		if r[k] != nil {
			m = r[k]
		}
		m[methodName] = b
		r[k] = m
	}

	return r
}

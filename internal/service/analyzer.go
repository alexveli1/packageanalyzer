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

// StorePacks triggers request to HTTP Client and saves received data into repository
func (as *AnalyzerService) StorePacks(ctx context.Context, branch string) error {
	p, err := as.client.GetRepo(ctx, branch)
	if err != nil {
		return nil
	}
	packs := make(domain.Branch)
	sources := make(domain.Sources)
	for i := 0; i < len(p); i++ {
		packs[p[i].Name] = append(packs[p[i].Name], p[i])
		if !contains(sources[p[i].Source], p[i].Name) {
			sources[p[i].Source] = append(sources[p[i].Source], p[i].Name)
		}
	}
	err = as.repo.SavePacks(ctx, branch, packs, sources)
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
	only := make(domain.Branch)
	for _, v := range packs {
		for i := 0; i < len(v); i++ {
			exists, err := as.repo.GetSource(ctx, branch2, v[i]) // unique if package.store doesn't exist in the other branch
			if err != nil {
				return nil, err
			}
			if !exists {
				only[v[i].Arch] = append(only[v[i].Arch], v[i])
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
	for _, pb1 := range packs {
		for i := 0; i < len(pb1); i++ {
			p2, exists, err := as.repo.PacksByArchAndNameExist(ctx, branch2, pb1[i])
			if err != nil {
				return nil, err
			}
			if pb1[i].Name == "binutils-arm-linux-gnueabihf" {
				l := pb1[i].Name
				_ = l
			}
			if exists {
			forFoundPack:
				for j := 0; j < len(p2); j++ {
					if p1VersionHigher(pb1[i], p2[j]) {
						r[pb1[i].Arch] = append(r[pb1[i].Arch], pb1[i])

						break forFoundPack
					}
				}
			} else {
				// r[pb1[i].Arch] = append(r[pb1[i].Arch], pb1[i])
			}
		}
	}

	return convertToResult(r, branch1, domain.MethodHigher), nil
}

func (as *AnalyzerService) GetVerificationInfo(ctx context.Context, branch1, branch2 string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	p, err := as.client.GetOfficialDiff(ctx, branch1, branch2)
	if err != nil {
		return err
	}
	compareBranch, err := convertToCompareStore(ctx, p)
	if err != nil {
		return err
	}
	err = as.repo.SaveComparison(ctx, compareBranch)
	if err != nil {
		return err
	}
	return nil
}

func (as *AnalyzerService) VerifyByMethod(ctx context.Context, branch string, method string, myResult domain.Result) ([]string, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	b, err := as.repo.GetMethodComparison(ctx, branch, method)
	if err != nil {
		return nil, err
	}
	failures := make([]string, 0)
	for _, methods := range myResult { // not taking into account architecture
		rb := methods[method][branch]
		sort.Slice(rb, func(i, j int) bool {
			return rb[i].Name < rb[j].Name
		})
		for j := 0; j < len(rb); j++ {
			if !contains(b, rb[j].Source) && rb[j].Source != "" {
				if !contains(failures, rb[j].Source) {
					failures = append(failures, rb[j].Source)
				}
			}
		}
	}

	return failures, nil
}
func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

// p1VersionHigher compares ALT package versions using rpm logic
func p1VersionHigher(ver1, ver2 domain.Binpack) bool {
	v1 := version.NewVersion(fmt.Sprint(ver1.Epoch) + ":" + ver1.Release + "-" + ver1.Version)
	v2 := version.NewVersion(fmt.Sprint(ver2.Epoch) + ":" + ver2.Release + "-" + ver2.Version)

	return v1.GreaterThan(v2)
}

// convertToResult converts domain.Branch structure into domain.Result for returning to usecase layer
func convertToResult(archPkgs domain.Branch, branchName, methodName string) domain.Result {
	r := make(domain.Result)
	for arch, packs := range archPkgs {
		b := make(domain.Branch)
		m := make(domain.Method)
		b[branchName] = packs
		if r[arch] != nil {
			m = r[arch]
		}
		m[methodName] = b
		r[arch] = m
	}

	return r
}
func convertToCompareStore(ctx context.Context, compResult *domain.VerificationInfo) (*domain.CompareBranch, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	cb := make(domain.CompareBranch)
	unique := make(domain.Compare)
	higher := make(domain.Compare)
	pks := compResult.Packages
	for i := 0; i < len(pks); i++ {
		if pks[i].Package2.Name == "" {
			unique[pks[i].Pkgset1] = append(unique[pks[i].Pkgset1], pks[i].Package1.Name)

			continue
		}
		v1 := domain.Binpack{
			Version: pks[i].Package1.Version,
			Release: pks[i].Package1.Release,
		}
		v2 := domain.Binpack{
			Version: pks[i].Package2.Version,
			Release: pks[i].Package2.Release,
		}
		if p1VersionHigher(v1, v2) {
			higher[pks[i].Pkgset1] = append(higher[pks[i].Pkgset1], pks[i].Package1.Name)
		}
	}
	c := unique[pks[0].Pkgset1]
	sort.Slice(c, func(i, j int) bool {
		return c[i] < c[j]
	})
	cb[domain.MethodUnique] = unique
	cb[domain.MethodHigher] = higher
	return &cb, nil
}

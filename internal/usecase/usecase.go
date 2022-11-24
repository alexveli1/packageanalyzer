package usecase

import (
	"context"
	"fmt"

	"github/alexveli1/packageanalyzer/internal/domain"
	"github/alexveli1/packageanalyzer/internal/service"
)

type Usecase struct {
	*service.Services
	result []domain.ResultsOutput
}

func NewUsecase(services *service.Services) *Usecase {
	return &Usecase{
		Services: services,
		result:   []domain.ResultsOutput{},
	}
}

func (u *Usecase) GetPacks(ctx context.Context, branch string) error {
	return u.Analyzer.GetPacks(ctx, branch)
}

func (u *Usecase) UniqueBranchPackages(ctx context.Context, branch1 string, branch2 string) {
	b1Packs, b2Packs := u.Analyzer.PackagesFromBranch1(ctx, branch1, branch2)
	for k, v := range b1Packs {
		u.result = append(u.result, domain.ResultsOutput{
			Branch:        branch1,
			Method:        domain.MethodUnique,
			Arch:          k,
			PackagesCount: len(v),
		})
	}
	for k, v := range b2Packs {
		u.result = append(u.result, domain.ResultsOutput{
			Branch:        branch2,
			Method:        domain.MethodUnique,
			Arch:          k,
			PackagesCount: len(v),
		})
	}
}

// GetHigherReleases collect and save info on packages which ver+rel is higher than in another branch
func (u *Usecase) GetHigherReleases(ctx context.Context, branch1 string, branch2 string) {
	packs, failures := u.Services.Analyzer.Branch1Higher(ctx, branch1, branch2)
	for k, v := range packs {
		u.result = append(u.result, domain.ResultsOutput{
			Branch:        branch1,
			Method:        domain.MethodHigher,
			Arch:          k,
			PackagesCount: len(v),
		})
	}
	if len(failures) > 0 {
		u.result = append(u.result, domain.ResultsOutput{
			Branch:        branch1,
			Method:        domain.MethodFailure,
			Arch:          "",
			PackagesCount: len(failures),
		})
	}
}

func (u *Usecase) PrintResult() {
	for i := 0; i < len(u.result); i++ {
		fmt.Printf("%v\n", u.result[i])
	}
}

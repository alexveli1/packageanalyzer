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
	/*	fmt.Printf("*********Unique packages in branch %s\n", branch1)
	 */for k, v := range b1Packs {
		u.result = append(u.result, domain.ResultsOutput{
			Branch:        branch1,
			Method:        domain.MethodUnique,
			Arch:          k,
			PackagesCount: len(v),
		})
		/*		fmt.Printf("%s : %d \n", k, len(v))
		 */ // fmt.Printf("%v\n", v)

	}
	/*	fmt.Printf("*********Unique packages in branch %s\n", branch1)
	 */for k, v := range b2Packs {
		u.result = append(u.result, domain.ResultsOutput{
			Branch:        branch2,
			Method:        domain.MethodUnique,
			Arch:          k,
			PackagesCount: len(v),
		})
		/*		fmt.Printf("%s : %d \n", k, len(v))*/
		// fmt.Printf("%v\n", v)

	}
	/*	fmt.Printf("\nPackages existing in %s repository only %d \n%v", branch1, len(b1Packs), b1Packs)
		fmt.Printf("\nPackages existing in %s repository only %d \n%v", branch2, len(b2Packs), b2Packs)*/
}

// GetHigherReleases collect and save info on packages which ver+rel is higher than in another branch
func (u *Usecase) GetHigherReleases(ctx context.Context, branch1 string, branch2 string) {
	packs, failures := u.Services.Analyzer.Branch1Higher(ctx, branch1, branch2)
	// fmt.Printf("*********Higher releases in %s\n", branch1)
	for k, v := range packs {
		u.result = append(u.result, domain.ResultsOutput{
			Branch:        branch1,
			Method:        domain.MethodHigher,
			Arch:          k,
			PackagesCount: len(v),
		})
		/*		fmt.Printf("%s - packages count %d\n", k, len(v))
				// fmt.Printf("%v\n", v)*/
	}
	if len(failures) > 0 {
		u.result = append(u.result, domain.ResultsOutput{
			Branch:        branch1,
			Method:        domain.MethodFailure,
			Arch:          "",
			PackagesCount: len(failures),
		})

		/*fmt.Printf("failed to compare package release or versions for %d packages\n", len(failures))*/
	}
	/*for k, v := range failures {
		fmt.Printf("%s: ---- %s\n", k, v)
	}*/
}

func (u *Usecase) PrintResult() {
	fmt.Printf("%v", u.result)
}

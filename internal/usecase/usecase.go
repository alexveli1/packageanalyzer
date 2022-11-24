package usecase

import (
	"context"
	"fmt"

	"github/alexveli1/packageanalyzer/internal/service"
)

type Usecase struct {
	*service.Services
}

func NewUsecase(services *service.Services) *Usecase {
	return &Usecase{
		Services: services,
	}
}

func (u *Usecase) GetPacks(ctx context.Context, branch string) error {
	return u.Analyzer.GetPacks(ctx, branch)
}

func (u *Usecase) UniqueBranchPackages(ctx context.Context, branch1 string, branch2 string) {
	b1Packs, b2Packs := u.Analyzer.PackagesFromBranch1(ctx, branch1, branch2)
	fmt.Printf("\nPackages existing in %s repository only %d \n", branch1, len(b1Packs))
	fmt.Printf("\nPackages existing in %s repository only %d \n", branch2, len(b2Packs))
	/*fmt.Printf("\nPackages existing in %s repository only %d \n%v", branch1, len(b1Packs), b1Packs)
	fmt.Printf("\nPackages existing in %s repository only %d \n%v", branch2, len(b2Packs), b2Packs)*/
}

// GetHigherReleases collect and save info on packages which ver+rel is higher than in another branch
func (u *Usecase) GetHigherReleases(ctx context.Context, branch1 string, branch2 string) {
	packs := u.Services.Analyzer.Branch1Higher(ctx, branch1, branch2)
	fmt.Printf("\nPackages with higher releases in %s than in %s is %d\n", branch1, branch2, len(packs))
	//fmt.Printf("%v", packs)
}

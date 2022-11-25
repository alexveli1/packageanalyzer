package usecase

import (
	"context"
	"encoding/json"
	"os"

	"github/alexveli1/packageanalyzer/internal/domain"
	"github/alexveli1/packageanalyzer/internal/service"
	"github/alexveli1/packageanalyzer/pkg/mylog"
)

type Usecase struct {
	*service.Services
	result map[string]map[string]map[string][]domain.Binpack
}

func NewUsecase(services *service.Services) *Usecase {
	return &Usecase{
		Services: services,
		result:   make(map[string]map[string]map[string][]domain.Binpack),
	}
}

func (u *Usecase) GetPacks(ctx context.Context, branch string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	mylog.SugarLogger.Infof("Getting packages for branch %s", branch)
	return u.Analyzer.GetPacks(ctx, branch)
}

func (u *Usecase) UniqueBranchPackages(ctx context.Context, branch1 string, branch2 string) {
	if err := ctx.Err(); err != nil {
		return
	}
	mylog.SugarLogger.Infof("Analyzing unique packages for branch %s", branch1)
	b1Packs, b2Packs := u.Analyzer.PackagesFromBranch1(ctx, branch1, branch2)
	branch := make(map[string][]domain.Binpack)
	method := make(map[string]map[string][]domain.Binpack)
	for k, v := range b1Packs {
		branch[branch1] = v
		method["unique"] = branch
		u.result[k] = method
	}
	mylog.SugarLogger.Infof("Analyzing unique packages for branch %s", branch2)
	for k, v := range b2Packs {
		branch = u.result[k]["unique"]
		branch[branch2] = v
		u.result[k]["unique"] = branch
	}
}

// GetHigherReleases collect and save info on packages which ver+rel is higher than in another branch
func (u *Usecase) GetHigherReleases(ctx context.Context, branch1 string, branch2 string) {
	if err := ctx.Err(); err != nil {
		return
	}
	mylog.SugarLogger.Infof("Analyzing release differences for branch %s", branch1)
	packs := u.Services.Analyzer.Branch1Higher(ctx, branch1, branch2)
	branch := make(map[string][]domain.Binpack)
	for k, v := range packs {
		branch[branch1] = v
		method := u.result[k]
		method["higher"] = branch
		u.result[k] = method
	}
}

func (u *Usecase) PrintResult() {
	for k, v := range u.result {
		mylog.SugarLogger.Infof("Writing file %s.json", k)
		data, _ := json.Marshal(v)
		err := os.WriteFile(k+".json", data, 0666)
		if err != nil {
			mylog.SugarLogger.Warnf("cannot write json file")
		}
	}
}

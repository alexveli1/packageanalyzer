// Package usecase implements api level operations with top level parameters
// sends requests to service layer and prints out results
package usecase

import (
	"context"
	"encoding/json"
	"os"

	"github/alexveli1/packageanalyzer/internal/domain"
	"github/alexveli1/packageanalyzer/internal/service"
	"github/alexveli1/packageanalyzer/pkg/mylog"
)

// Usecase - api level object for receiving instructions from top most layer of app
type Usecase struct {
	*service.Services
	result domain.Result
}

func NewUsecase(services *service.Services) *Usecase {
	return &Usecase{
		Services: services,
		result:   make(domain.Result),
	}
}

// GetPackages requests service layer to collect and save ALT repository data
func (u *Usecase) GetPackages(ctx context.Context, branch string) {
	if err := ctx.Err(); err != nil {
		mylog.SugarLogger.Infof("Getting packages for branch %s failed: %v", branch, err)

		return
	}
	mylog.SugarLogger.Infof("Getting packages for branch %s", branch)
	err := u.Analyzer.GetPacks(ctx, branch)
	if err != nil {
		mylog.SugarLogger.Infof("Getting packages for branch %s failed: %v", branch, err)
	}
}

// GetUniquePackages requests service layer to provide domain.Result for packages unique for branch1
// and stores results into Usecase object for later printing of all results
func (u *Usecase) GetUniquePackages(ctx context.Context, branch1 string, branch2 string) {
	if err := ctx.Err(); err != nil {
		mylog.SugarLogger.Infof("Getting unique packages for branch %s failed: %v", branch1, err)

		return
	}
	mylog.SugarLogger.Infof("Analyzing unique packages for branch %s", branch1)
	packs, err := u.Analyzer.GetUnique(ctx, branch1, branch2)
	if err != nil {
		mylog.SugarLogger.Warnf("cannot get unique packages from branch %s: %v", branch1, err)

		return
	}
	u.appendResult(packs)
}

// GetHigherReleases collect and save info on packages which ver+rel is higher than in another branch
// and stores results into Usecase object for later printing of all results
func (u *Usecase) GetHigherReleases(ctx context.Context, branch1 string, branch2 string) {
	if err := ctx.Err(); err != nil {
		mylog.SugarLogger.Infof("Getting higher releases for branch %s failed: %v", branch1, err)

		return
	}
	mylog.SugarLogger.Infof("Analyzing release differences for branch %s", branch1)
	packs, err := u.Services.Analyzer.GetHigher(ctx, branch1, branch2)
	if err != nil {
		mylog.SugarLogger.Warnf("cannot get higher release packages from branch %s: %v", branch1, err)

		return
	}
	u.appendResult(packs)
}

// PrintResult manages format of output - file, stdout, etc
// logic might be extended with flags and other config tools
func (u *Usecase) PrintResult(ctx context.Context) {
	if err := ctx.Err(); err != nil {
		mylog.SugarLogger.Infof("Result printing failed: %v", err)

		return
	}
	for k, v := range u.result {
		mylog.SugarLogger.Infof("Writing file %s.json", k)
		data, _ := json.Marshal(v)
		err := os.WriteFile(k+".json", data, 0666)
		if err != nil {
			mylog.SugarLogger.Warnf("cannot write json file: %v", err)

			return
		}
	}
}

// appendResult stores domain.Result of processing single branch with single method (unique or higher)
// received from service layer
func (u *Usecase) appendResult(packs domain.Result) {
	for k, v := range packs { // arch
		if u.result[k] == nil {
			u.result[k] = v

			continue
		}
		for k1, v1 := range v { // method
			if u.result[k][k1] == nil {
				u.result[k][k1] = v1

				continue
			}
			for k2, v2 := range v1 { // branch
				u.result[k][k1][k2] = v2
			}
		}
	}
}

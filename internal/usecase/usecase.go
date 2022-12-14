// Package usecase implements api level operations with top level parameters
// sends requests to service layer and prints out results
package usecase

import (
	"context"
	"encoding/json"
	"os"
	"sync"

	"github/alexveli1/packageanalyzer/internal/domain"
	"github/alexveli1/packageanalyzer/internal/service"
	"github/alexveli1/packageanalyzer/pkg/mylog"
)

// Usecase - api level object for receiving instructions from top most layer of app
type Usecase struct {
	*service.Services
	result domain.Result
	lock   *sync.RWMutex
}

func NewUsecase(services *service.Services) *Usecase {
	return &Usecase{
		Services: services,
		result:   make(domain.Result),
		lock:     &sync.RWMutex{},
	}
}

// GetPackages requests service layer to collect and save ALT repository data
func (u *Usecase) GetPackages(ctx context.Context, branch string) error {
	mylog.SugarLogger.Infof("Started collecting packages info for branch %s", branch)
	if err := ctx.Err(); err != nil {
		mylog.SugarLogger.Warnf("Getting packages for branch %s failed: %v", branch, err)

		return err
	}
	err := u.Analyzer.StorePacks(ctx, branch)
	if err != nil {
		mylog.SugarLogger.Warnf("Getting packages for branch %s failed: %v", branch, err)

		return err
	}

	return nil
}

// GetUniquePackages requests service layer to provide domain.Result for packages unique for branch1
// and stores results into Usecase object for later printing of all results
func (u *Usecase) GetUniquePackages(ctx context.Context, branch1 string, branch2 string) error {
	if err := ctx.Err(); err != nil {
		mylog.SugarLogger.Warnf("Getting unique packages for branch %s failed: %v", branch1, err)

		return err
	}
	mylog.SugarLogger.Infof("Analyzing unique packages for branch %s", branch1)
	packs, err := u.Analyzer.GetUnique(ctx, branch1, branch2)
	if err != nil {
		mylog.SugarLogger.Warnf("cannot get unique packages from branch %s: %v", branch1, err)

		return err
	}
	u.appendResult(packs)

	return nil
}

// GetHigherReleases collect and save info on packages which ver+rel is higher than in another branch
// and stores results into Usecase object for later printing of all results
func (u *Usecase) GetHigherReleases(ctx context.Context, branch1 string, branch2 string) error {
	if err := ctx.Err(); err != nil {
		mylog.SugarLogger.Warnf("Getting higher releases for branch %s failed: %v", branch1, err)

		return err
	}
	mylog.SugarLogger.Infof("Analyzing release differences for branch %s", branch1)
	packs, err := u.Services.Analyzer.GetHigher(ctx, branch1, branch2)
	if err != nil {
		mylog.SugarLogger.Warnf("cannot get higher release packages from branch %s: %v", branch1, err)

		return err
	}
	u.appendResult(packs)

	return nil
}

func (u *Usecase) GetVerificationData(ctx context.Context, branch1, branch2 string) error {
	err := u.Services.Analyzer.GetVerificationInfo(ctx, branch1, branch2)
	if err != nil {
		mylog.SugarLogger.Warnf("cannot get verification info for branches %s, %s: %v", branch1, branch2, err)

		return err
	}
	return nil
}

func (u *Usecase) VerifyResult(ctx context.Context, branch string) error {
	failures, err := u.Services.Analyzer.VerifyByMethod(ctx, branch, domain.MethodUnique, u.result)
	if err != nil {
		mylog.SugarLogger.Warnf("cannot verify results for branch %s: %v", branch, err)

		return err
	}
	if len(failures) > 0 {
		mylog.SugarLogger.Warnf("Unique packages not in official data in branch %s:\n%v\n", branch, failures)

		return nil
	}
	mylog.SugarLogger.Infof("Uniqueness verification successful for branch %s", branch)
	failures, err = u.Services.Analyzer.VerifyByMethod(ctx, branch, domain.MethodHigher, u.result)
	if err != nil {
		mylog.SugarLogger.Warnf("cannot verify results for branch %s: %v", branch, err)

		return err
	}
	if len(failures) > 0 {
		mylog.SugarLogger.Warnf("Higher releases not in official data in branch  %s:\n%v\n", branch, failures)

		return nil
	}
	mylog.SugarLogger.Infof("Release difference verification was successful for branch %s", branch)

	return nil
}

// PrintResult manages format of output - file, stdout, etc
// logic might be extended with flags and other config tools
func (u *Usecase) PrintResult(ctx context.Context) {
	if err := ctx.Err(); err != nil {
		mylog.SugarLogger.Warnf("Result printing failed: %v", err)

		return
	}
	data, _ := json.Marshal(u.result)
	_ = os.WriteFile("result.json", data, 0666)
}

// appendResult stores domain.Result of processing single branch with single method (unique or higher)
// received from service layer
func (u *Usecase) appendResult(result domain.Result) {
	u.lock.Lock()
	defer u.lock.Unlock()
	for arch, methods := range result { // arch
		if u.result[arch] == nil {
			u.result[arch] = methods

			continue
		}
		for method, branches := range methods { // method
			if u.result[arch][method] == nil {
				u.result[arch][method] = branches

				continue
			}
			for branch, packs := range branches { // branch
				u.result[arch][method][branch] = packs
			}
		}
	}
}

package usecase

import "github/alexveli1/packageanalyzer/internal/service"

type Usecase struct {
	*service.Services
}

func NewUsecase(services *service.Services) *Usecase {
	return &Usecase{
		Services: services,
	}
}

// GetXOR collect and output info on packs existing in 1 and missing in 2nd branch
func (u *Usecase) GetXOR() {

}

// GetHigherReleases collect and save info on packages which ver+rel is higher than in another branch
func (u *Usecase) GetHigherReleases() {

}

func GetSisyphusBranch() {

}

func GetP10Branch() {

}

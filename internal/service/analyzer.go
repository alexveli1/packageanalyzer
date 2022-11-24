package service

import (
	"context"
	"sort"

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
func (as *AnalyzerService) PackagesFromBranch1(ctx context.Context, branch1 string, branch2 string) ([]string, []string) {
	if err := ctx.Err(); err != nil {

		return nil, nil
	}
	chB1 := make(chan []string)
	chB2 := make(chan []string)
	go as.uniquePacks(ctx, chB1, branch1, branch2)
	go as.uniquePacks(ctx, chB2, branch2, branch1)

	return <-chB1, <-chB2
}
func (as *AnalyzerService) Branch1Higher(ctx context.Context, branch1 string, branch2 string) []string {
	packs, _ := as.repo.GetAllPacks(ctx, branch1)
	branch1Higher := make([]string, 0)
	for k, p1s := range packs {
		if p2s, exists := as.repo.GetPackByName(ctx, branch2, k); exists {
			if version1IsGreater(p1s[len(p1s)-1].Version, p2s[len(p2s)-1].Version) {
				branch1Higher = append(branch1Higher, k+"\n")
			}
		}
	}
	return branch1Higher
}
func (as *AnalyzerService) uniquePacks(ctx context.Context, ch chan []string, branch1, branch2 string) {
	packs, _ := as.repo.GetAllPacks(ctx, branch1)
	only := make([]string, 0)
	for k := range packs {
		if _, exists := as.repo.GetPackByName(ctx, branch2, k); !exists {
			only = append(only, k+"\n")
		}
	}
	ch <- only
}
func version1IsGreater(ver1, ver2 string) bool {
	return ver1 > ver2
}

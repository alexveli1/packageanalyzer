package service

import (
	"context"
	"errors"
	"sort"
	"strconv"
	"strings"

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
func (as *AnalyzerService) PackagesFromBranch1(ctx context.Context, branch1 string, branch2 string) (map[string][]string, map[string][]string) {
	if err := ctx.Err(); err != nil {
		return nil, nil
	}
	chB1 := make(chan map[string][]string)
	chB2 := make(chan map[string][]string)
	go as.uniquePacks(ctx, chB1, branch1, branch2)
	go as.uniquePacks(ctx, chB2, branch2, branch1)

	return <-chB1, <-chB2
}
func (as *AnalyzerService) Branch1Higher(ctx context.Context, branch1 string, branch2 string) (map[string][]string, map[string]string) {
	packs, _ := as.repo.GetAllPacks(ctx, branch1)
	branch1Higher := make(map[string][]string, 0)
	failures := make(map[string]string)
	for pkgName, p1archs := range packs {
		for i := 0; i < len(p1archs); i++ {
			pkgBranch2, exists := as.repo.GetPackByArchAndName(ctx, branch2, p1archs[i].Arch, pkgName)
			if exists {
				p1Rel, p2Rel := p1archs[i].Release, pkgBranch2.Release
				p1Ver, p2Ver := p1archs[i].Version, pkgBranch2.Version
				yes, err := version1IsGreater(p1Rel, p2Rel, p1Ver, p2Ver)
				if err != nil {
					if !stringMatch(p1Rel, p2Rel, p1Ver, p2Ver) {
						failures[pkgName+p1Rel+p1Ver] = "rel :" + p1Rel + " ver " + p1Ver + " and rel " + p2Rel + " ver " + p2Ver
					}
				}
				if yes {
					branch1Higher[p1archs[i].Arch] = append(branch1Higher[p1archs[i].Arch], pkgName+p1archs[i].Version+">"+pkgBranch2.Version+"\n")
				}
			}
		}
	}
	return branch1Higher, failures
}
func (as *AnalyzerService) uniquePacks(ctx context.Context, ch chan map[string][]string, branch1, branch2 string) {
	packs, _ := as.repo.GetAllPacks(ctx, branch1)
	only := make(map[string][]string, 0)
	for k, v := range packs {
		for i := 0; i < len(v); i++ {
			if _, exists := as.repo.GetPackByArchAndName(ctx, branch2, v[i].Arch, k); !exists {
				only[v[i].Arch] = append(only[v[i].Arch], k+"\n")
			}
		}
	}
	ch <- only
}
func version1IsGreater(rel1, rel2, ver1, ver2 string) (bool, error) {
	r1 := prepForSplit(strings.ReplaceAll(rel1, "alt", ""))
	r2 := prepForSplit(strings.ReplaceAll(rel2, "alt", ""))

	rs1 := strings.Split(r1, ".")
	rs2 := strings.Split(r2, ".")
	yes, err := firstSplitGreater(rs1, rs2)
	if err != nil {
		if errors.Is(err, domain.ErrSecondHigher) {
			return false, nil
		}
		return false, err
	}
	if yes {
		// fmt.Printf("%s+%s > %s+%s\n", rel1, ver1, rel2, ver2)

		return true, nil
	}
	v1 := strings.Split(prepForSplit(ver1), ".")
	v2 := strings.Split(prepForSplit(ver2), ".")
	yes, err = firstSplitGreater(v1, v2)
	if err != nil {
		if errors.Is(err, domain.ErrSecondHigher) {
			return false, nil
		}
		return false, err
	}
	if yes {
		// fmt.Printf("%s+%s > %s+%s\n", rel1, ver1, rel2, ver2)

		return true, nil
	}
	return false, nil
}

func firstSplitGreater(v1, v2 []string) (bool, error) {
	for i := 0; i < len(v1); i++ {
		vs1, err := strconv.ParseInt(v1[i], 10, 64)
		if err != nil {
			// mylog.SugarLogger.Warnf("cannot convert %s: %v", v1[i], err)

			return domain.ErrCompare, err
		}
		if len(v2) <= i {
			return domain.ErrCompare, nil
		}
		vs2, err := strconv.ParseInt(v2[i], 10, 64)
		if err != nil {
			// mylog.SugarLogger.Warnf("cannot convert %s: %v", v1[i], err)

			return domain.ErrCompare, err
		}
		if vs1 > vs2 {
			return domain.FirstHigher, nil
		}
		if vs1 < vs2 {
			return domain.SecondHigher, domain.ErrSecondHigher
		}
	}

	return domain.SecondHigher, nil
}

func prepForSplit(s string) string {
	s = strings.ReplaceAll(s, "_", ".")
	s = strings.ReplaceAll(s, "jpp", "")
	s = strings.ReplaceAll(s, "qa", "")
	s = strings.ReplaceAll(s, "svn", "")
	s = strings.ReplaceAll(s, ".git.", ".")
	s = strings.ReplaceAll(s, "git", "")
	s = strings.ReplaceAll(s, "r", "")

	return s
}

func stringMatch(p1Rel, p2Rel, p1Ver, p2Ver string) bool {
	return p1Rel == p2Rel && p1Ver == p2Ver
}

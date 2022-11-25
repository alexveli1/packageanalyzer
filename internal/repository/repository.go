// Package repository provides access to storage for service layer
package repository

import (
	"context"

	"github/alexveli1/packageanalyzer/internal/domain"
)

// Analyzer provides services with access to raw data collected from web api
type Analyzer interface {
	SavePacks(ctx context.Context, branch string, packs domain.Branch) error
	GetBranchPacks(ctx context.Context, branch string) (domain.Branch, error)
	GetPacksByArchAndName(ctx context.Context, branch string, arch string, packName string) ([]domain.Binpack, bool, error)
}

// Repositories consolidated object for providing repositories for seriveces
// though excessive for current task scope but used for extensibility
type Repositories struct {
	Analyzer
}

func NewRepositories() *Repositories {
	return &Repositories{
		Analyzer: NewAnalyzerRepo(),
	}
}

// Package storage provider of data objects for server repositories
package storage

import (
	"sync"

	"github/alexveli1/packageanalyzer/internal/domain"
)

// MapDB place for storing raw data collected from ALT' web api
type MapDB struct {
	DB   domain.Branch
	Lock *sync.RWMutex
}

func NewMap() *MapDB {
	return &MapDB{
		DB:   make(domain.Branch),
		Lock: &sync.RWMutex{},
	}
}

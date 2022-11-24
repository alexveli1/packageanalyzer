// Package storage provider of data objects for server repositories
package storage

import (
	"sync"

	"github/alexveli1/packageanalyzer/internal/domain"
)

type MapDB struct {
	DB   map[string][]domain.Binpack
	Lock *sync.RWMutex
}

func NewMap() *MapDB {
	return &MapDB{
		DB:   map[string][]domain.Binpack{},
		Lock: &sync.RWMutex{},
	}
}

// Package storage provider of data objects for server repositories
package storage

import (
	"sync"

	"github/alexveli1/packageanalyzer/internal/domain"
	"github/alexveli1/packageanalyzer/pkg/mylog"
)

type MapDB struct {
	DB   map[string]domain.Binpack
	Lock *sync.RWMutex
}

func NewMap() *MapDB {
	mylog.SugarLogger.Infof("repositories will use map")

	return &MapDB{
		DB:   map[string]domain.Binpack{},
		Lock: &sync.RWMutex{},
	}
}

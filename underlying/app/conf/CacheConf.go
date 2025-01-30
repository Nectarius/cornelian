package conf

import (
	"github.com/dgraph-io/ristretto"
)

type CacheConf struct {
	Cache *ristretto.Cache
}

func NewCacheConf() (*CacheConf, error) {

	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,
		MaxCost:     100000000,
		BufferItems: 1e7,
	})

	return &CacheConf{
		Cache: cache,
	}, err
}

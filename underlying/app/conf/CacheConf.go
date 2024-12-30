package conf

import (
	"fmt"

	"github.com/dgraph-io/ristretto"
)

type CacheConf struct {
	Cache *ristretto.Cache
}

func NewCacheConf() *CacheConf {

	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,
		MaxCost:     1000000,
		BufferItems: 1e7,
	})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	return &CacheConf{
		Cache: cache,
	}
}

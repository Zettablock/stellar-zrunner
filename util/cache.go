package util

import (
	cache "github.com/Code-Hex/go-generics-cache"
	"github.com/Code-Hex/go-generics-cache/policy/lfu"
)

const (
	CacheCapacity = 1000
)

var (
	PhoenixPoolsCache			*cache.Cache[string, int]
	BlendLendingPoolsCache		*cache.Cache[string, int]
	XycLoansLendingPoolsCache	*cache.Cache[string, int]
	LoadedDb					bool
)

func init() {
	PhoenixPoolsCache = cache.New(cache.AsLFU[string, int](lfu.WithCapacity(CacheCapacity)))
	BlendLendingPoolsCache = cache.New(cache.AsLFU[string, int](lfu.WithCapacity(CacheCapacity)))
	XycLoansLendingPoolsCache = cache.New(cache.AsLFU[string, int](lfu.WithCapacity(CacheCapacity)))
}

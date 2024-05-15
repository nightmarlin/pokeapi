package cache_test

import (
	"context"
	"sync"
	"testing"

	"github.com/nightmarlin/pokeapi/cache"
	"github.com/nightmarlin/pokeapi/cache/cachetest"
)

// wrappableCache is a simple get/put cache implementation suitable for wrapping
// by cache.Wrapper.
type wrappableCache struct {
	mux   sync.RWMutex
	cap   int
	store map[string]any
}

func (wc *wrappableCache) get(_ context.Context, url string) (any, bool) {
	defer wc.mux.RUnlock()
	wc.mux.RLock()
	v, ok := wc.store[url]
	return v, ok
}

func (wc *wrappableCache) put(_ context.Context, url string, value any) {
	defer wc.mux.Unlock()
	wc.mux.Lock()

	wc.store[url] = value
	for len(wc.store) > wc.cap {
		// evict an arbitrary key
		for k := range wc.store {
			delete(wc.store, k)
			break
		}
	}
}

func TestWrapper(t *testing.T) {
	t.Parallel()

	t.Run(
		"cache implementation",
		func(t *testing.T) {
			t.Parallel()
			cachetest.TestCache(
				t,
				func(size int) *cache.Wrapper {
					wc := wrappableCache{cap: size, store: make(map[string]any)}
					return cache.NewWrapper(wc.get, wc.put)
				},
			)
		},
	)
}

package cache_test

import (
	"testing"

	"github.com/nightmarlin/pokeapi/cache"
	"github.com/nightmarlin/pokeapi/cache/cachetest"
)

func TestLRU(t *testing.T) {
	t.Parallel()

	t.Run(
		"cache implementation",
		func(t *testing.T) {
			t.Parallel()
			cachetest.TestCache(t, cache.NewLRU)
		},
	)
}

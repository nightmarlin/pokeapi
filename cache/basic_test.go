package cache_test

import (
	"testing"

	"github.com/nightmarlin/pokeapi/cache"
	"github.com/nightmarlin/pokeapi/cache/cachetest"
)

func TestBasic(t *testing.T) {
	t.Parallel()

	t.Run(
		"cache implementation",
		func(t *testing.T) {
			t.Parallel()
			cachetest.TestCache(t, cache.NewBasic)
		},
	)

	t.Run(
		"evicts oldest cached value",
		func(t *testing.T) {
			t.Parallel()

			const (
				resourceA = "https://pokeapi.co/api/v2/spinda"
				valueA    = "small bear"
				resourceB = "https://pokeapi.co/api/v2/lillipup"
				valueB    = "small dog"
				resourceC = "https://pokeapi.co/api/v2/pidove"
				valueC    = "small bird"
			)

			c := cache.NewBasic(2)
			c.Lookup(resourceA).Hydrate(valueA)
			c.Lookup(resourceB).Hydrate(valueB)
			c.Lookup(resourceC).Hydrate(valueC)

			lookup := c.Lookup(resourceA)
			defer lookup.Close()

			if _, ok := lookup.Value(); ok {
				t.Errorf("want resource %q to be evicted, but it wasn't", resourceA)
			}
		},
	)
}

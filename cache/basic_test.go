package cache_test

import (
	"context"
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
			ctx := context.Background()

			c := cache.NewBasic(2)
			c.Lookup(ctx, resourceA).Hydrate(ctx, valueA)
			c.Lookup(ctx, resourceB).Hydrate(ctx, valueB)
			c.Lookup(ctx, resourceC).Hydrate(ctx, valueC)

			lookup := c.Lookup(ctx, resourceA)
			defer lookup.Close(ctx)

			if _, ok := lookup.Value(ctx); ok {
				t.Errorf("want resource %q to be evicted, but it wasn't", resourceA)
			}
		},
	)
}

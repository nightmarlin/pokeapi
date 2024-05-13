package cache_test

import (
	"context"
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

	t.Run(
		"evicts least-recently-used item",
		func(t *testing.T) {
			ctx := context.Background()
			c := cache.NewLRU(3)

			c.Lookup(ctx, "https://pokeapi.co/api/v2/wooper").
				Hydrate(ctx, "a ground-type pokemon")

			c.Lookup(ctx, "https://pokeapi.co/api/v2/dragalge").
				Hydrate(ctx, "a poison-type pokemon")

			c.Lookup(ctx, "https://pokeapi.co/api/v2/wooper").
				Hydrate(ctx, "a water-type pokemon") // wooper now more recent than dragalge

			c.Lookup(ctx, "https://pokeapi.co/api/v2/miltank").
				Hydrate(ctx, "a normal-type pokemon") // cache now full

			c.Lookup(ctx, "https://pokeapi.co/api/v2/necrozma").
				Hydrate(ctx, "a psychic-type pokemon") // dragalge should be evicted

			lookup := c.Lookup(ctx, "https://pokeapi.co/api/v2/dragalge")
			defer lookup.Close(ctx)

			if v, ok := lookup.Value(ctx); ok {
				t.Errorf("wanted lookup to return (nil, false); got (%v, %T)", v, ok)
			}
		},
	)
}

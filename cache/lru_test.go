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

	t.Run(
		"evicts least-recently-used item",
		func(t *testing.T) {
			c := cache.NewLRU(3)

			c.Lookup("https://pokeapi.co/api/v2/wooper").Hydrate("a ground-type pokemon")
			c.Lookup("https://pokeapi.co/api/v2/dragalge").Hydrate("a poison-type pokemon")
			c.Lookup("https://pokeapi.co/api/v2/wooper").Hydrate("a water-type pokemon")     // wooper now more recent than dragalge
			c.Lookup("https://pokeapi.co/api/v2/miltank").Hydrate("a normal-type pokemon")   // cache now full
			c.Lookup("https://pokeapi.co/api/v2/necrozma").Hydrate("a psychic-type pokemon") // dragalge should be evicted

			lookup := c.Lookup("https://pokeapi.co/api/v2/dragalge")
			defer lookup.Close()

			if v, ok := lookup.Value(); ok {
				t.Errorf("wanted lookup to return (nil, false); got (%v, %T)", v, ok)
			}
		},
	)
}

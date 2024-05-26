package cache_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/nightmarlin/pokeapi"
	"github.com/nightmarlin/pokeapi/cache"
	"github.com/nightmarlin/pokeapi/cache/cachetest"
)

func TestLRU(t *testing.T) {
	t.Parallel()

	t.Run(
		"cache implementation",
		func(t *testing.T) {
			t.Parallel()
			cachetest.TestCache(
				t,
				func(size int) pokeapi.Cache { return cache.NewLRU(&cache.LRUOpts{Size: size}) },
			)
		},
	)

	var (
		loader = func(s string) pokeapi.CacheLoader {
			return func(context.Context) (any, error) { return s, nil }
		}
		missCheckLoader = func() (pokeapi.CacheLoader, func() bool) {
			var missed bool
			return func(context.Context) (any, error) {
					missed = true
					return nil, nil
				},
				func() bool { return missed }
		}
	)

	t.Run(
		"evicts least-recently-used item",
		func(t *testing.T) {
			t.Parallel()

			var (
				ctx = context.Background()
				c   = cache.NewLRU(&cache.LRUOpts{Size: 3})
			)

			_, _ = c.Lookup(
				ctx,
				"https://pokeapi.co/api/v2/pokemon/wooper",
				loader("a ground-type pokemon"),
			)

			_, _ = c.Lookup(
				ctx,
				"https://pokeapi.co/api/v2/pokemon/dragalge",
				loader("a poison-type pokemon"),
			)

			_, _ = c.Lookup(
				ctx,
				"https://pokeapi.co/api/v2/pokemon/wooper",
				loader("a water-type pokemon"),
			) // wooper now more recent than dragalge

			_, _ = c.Lookup(
				ctx,
				"https://pokeapi.co/api/v2/pokemon/miltank",
				loader("a normal-type pokemon"),
			) // cache now full

			_, _ = c.Lookup(
				ctx,
				"https://pokeapi.co/api/v2/pokemon/necrozma",
				loader("a psychic-type pokemon"),
			) // dragalge should be evicted

			l, missed := missCheckLoader()
			v, _ := c.Lookup(ctx, "https://pokeapi.co/api/v2/pokemon/dragalge", l)

			if !missed() {
				t.Errorf("wanted lookup to miss; got %v", v)
			}
		},
	)

	t.Run(
		"evicts old items from the cache once ttl expires",
		func(t *testing.T) {
			t.Parallel()

			var (
				ctx          = context.Background()
				zeroDuration = time.Duration(0)
				clockTimeMux = sync.RWMutex{}
				clockTime    = time.Now()
				c            = cache.NewLRU(
					&cache.LRUOpts{
						Size: 2,
						Clock: func() time.Time {
							defer clockTimeMux.RUnlock()
							clockTimeMux.RLock()
							return clockTime
						},
						TTL:         time.Minute,
						ExpiryDelay: &zeroDuration,
					},
				)

				// the expire() function is called as a background goroutine at the end
				// of a Lookup() call. we have no way to wait for it here, so we rely on
				// time.Sleep's goroutine-friendliness & a little bit of luck.
				sleepForExpiryGoRoutine = func() {
					for range 5 {
						time.Sleep(time.Millisecond)
					}
				}
			)

			_, _ = c.Lookup(
				ctx,
				"https://pokeapi.co/api/v2/pokemon/copperajah",
				loader("a steel-type pokemon"),
			)

			sleepForExpiryGoRoutine()

			clockTimeMux.Lock()
			clockTime = clockTime.Add(time.Hour) // next call to Lookup should evict previous entry
			clockTimeMux.Unlock()

			_, _ = c.Lookup(
				ctx,
				"https://pokeapi.co/api/v2/pokemon/mamoswine",
				loader("an ice-type pokemon"),
			)

			sleepForExpiryGoRoutine()

			l, m := missCheckLoader()
			_, _ = c.Lookup(ctx, "https://pokeapi.co/api/v2/pokemon/copperajah", l)

			if !m() {
				t.Errorf("wanted lookup to miss; but it didn't")
			}
		},
	)

	t.Run(
		"skips matching urls when required",
		func(t *testing.T) {
			t.Parallel()

			var (
				ctx = context.Background()
				c   = cache.NewLRU(
					&cache.LRUOpts{
						Size: 1,
						SkipURL: func(url string) bool {
							return url == "https://pokeapi.co/api/v2/pokemon?offset=20&limit=20"
						},
					},
				)

				loadOnMissCallCount int
				loadOnMiss          = func(context.Context) (any, error) {
					loadOnMissCallCount += 1
					return pokeapi.Pokemon{NamedIdentifier: pokeapi.NamedIdentifier{Name: "magearna"}}, nil
				}
			)

			_, _ = c.Lookup(
				ctx,
				"https://pokeapi.co/api/v2/pokemon?offset=20&limit=20",
				loadOnMiss,
			) // should skip cache

			_, _ = c.Lookup(
				ctx,
				"https://pokeapi.co/api/v2/pokemon/magearna",
				loadOnMiss,
			) // should be loaded into cache

			_, _ = c.Lookup(
				ctx,
				"https://pokeapi.co/api/v2/pokemon?offset=20&limit=20",
				loadOnMiss,
			) // should also skip cache

			_, _ = c.Lookup(
				ctx,
				"https://pokeapi.co/api/v2/pokemon/magearna",
				loadOnMiss,
			) // should be fetched from cache

			if loadOnMissCallCount != 3 {
				t.Errorf("want loadOnMiss to be called 3 times; got %d", loadOnMissCallCount)
			}
		},
	)
}

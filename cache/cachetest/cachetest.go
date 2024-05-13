// Package cachetest is a test suite for [pokeapi.Cache] implementations. Simply
// call [cachetest.TestCache] from your cache's test file.
//
// As it imports package testing, cachetest should not be used in normal
// application code.
package cachetest

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/nightmarlin/pokeapi"
)

type NewCacheFn[T pokeapi.Cache] func(size int) T

func TestCache[C pokeapi.Cache](t *testing.T, newCache NewCacheFn[C]) {
	t.Run(
		"lookup on an empty cache misses",
		func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			cache := newCache(1)

			lookup := cache.Lookup(ctx, "https://pokapi.co/api/v2/pokemon/gumshoos")
			defer lookup.Close(ctx)

			if v, ok := lookup.Value(ctx); ok || v != nil {
				t.Errorf("want lookup.Value() to return (nil, false); got (%v, %T)", v, ok)
			}
		},
	)

	t.Run(
		"lookup after hydration returns value",
		func(t *testing.T) {
			t.Parallel()

			const (
				resource = "https://pokapi.co/api/v2/pokemon/yamask"
				value    = "a ghost-type pokemon"
			)

			ctx := context.Background()
			cache := newCache(1)

			cache.Lookup(ctx, resource).Hydrate(ctx, value)

			lookup := cache.Lookup(ctx, resource)
			defer lookup.Close(ctx)

			if v, ok := lookup.Value(ctx); !ok || value != v {
				t.Errorf(`want lookup.Value() to return (%q, true); got (%v, %T)`, value, value, ok)
			}
		},
	)

	t.Run(
		"multiple hydration of same lookup returns first value",
		func(t *testing.T) {
			t.Parallel()
			const (
				resource = "https://pokapi.co/api/v2/pokemon/stunfisk"
				valueA   = "an electric-type pokemon"
				valueB   = "a ground-type pokemon"
			)

			ctx := context.Background()
			cache := newCache(1)

			lookup := cache.Lookup(ctx, resource)
			lookup.Hydrate(ctx, valueA) // only the first Hydrate() call should have an effect
			lookup.Hydrate(ctx, valueB)
			lookup.Close(ctx)

			lookup = cache.Lookup(ctx, resource)
			defer lookup.Close(ctx)

			if got, ok := lookup.Value(ctx); !ok || valueA != got {
				t.Errorf(`want lookup.Value() to return (%q, true); got (%v, %T)`, valueA, got, ok)
			}
		},
	)

	t.Run(
		"hydration of lookup after close has no effect",
		func(t *testing.T) {
			t.Parallel()

			const (
				resource = "https://pokapi.co/api/v2/pokemon/surskit"
				value    = "a water-type pokemon"
			)

			ctx := context.Background()
			cache := newCache(1)

			lookup := cache.Lookup(ctx, resource)
			lookup.Close(ctx)
			lookup.Hydrate(ctx, value) // hydrate after close should have no effect

			lookup = cache.Lookup(ctx, resource) // this lookup should miss
			defer lookup.Close(ctx)

			if got, ok := lookup.Value(ctx); ok {
				t.Errorf(`want lookup.Value() to return (nil, false); got (%v, %T)`, got, ok)
			}
		},
	)

	t.Run(
		"for a cache of size N, hydration of N resources means no lookups for those resources miss",
		func(t *testing.T) {
			t.Parallel()

			resource := func(i int) string { return fmt.Sprintf("resource-%d", i) }

			for _, entry := range []struct{ N int }{
				{N: 1},
				{N: 10},
				{N: 100},
			} {
				t.Run(
					fmt.Sprintf("%d resources", entry.N),
					func(t *testing.T) {
						t.Parallel()

						ctx := context.Background()
						cache := newCache(entry.N)

						for i := range entry.N {
							cache.Lookup(ctx, resource(i)).Hydrate(ctx, i)
						}

						var misses int
						for i := range entry.N {
							lookup := cache.Lookup(ctx, resource(i))
							if _, ok := lookup.Value(ctx); !ok {
								misses++
							}
							lookup.Close(ctx)
						}

						if misses != 0 {
							t.Errorf("want 0 cache misses; got %d", misses)
						}
					},
				)
			}
		},
	)

	t.Run(
		"lookup after multiple hydration returns most recent value",
		func(t *testing.T) {
			t.Parallel()

			const (
				resource    = "https://pokapi.co/api/v2/pokemon/aegislash"
				firstValue  = "a ghost-type pokemon"
				secondValue = "a steel-type pokemon"
			)

			ctx := context.Background()
			cache := newCache(1)

			cache.Lookup(ctx, resource).Hydrate(ctx, firstValue)
			cache.Lookup(ctx, resource).Hydrate(ctx, secondValue)

			lookup := cache.Lookup(ctx, resource)
			defer lookup.Close(ctx)

			if v, ok := lookup.Value(ctx); !ok || v != secondValue {
				t.Errorf(`want lookup.Value() to return (%q, true); got (%v, %T)`, secondValue, v, ok)
			}
		},
	)

	t.Run(
		"cached value is unchanged by other lookups",
		func(t *testing.T) {
			t.Parallel()

			const (
				resourceA = "https://pokapi.co/api/v2/pokemon/serperior"
				valueA    = "a grass-type pokemon"
				resourceB = "https://pokapi.co/api/v2/pokemon/dreepy"
				valueB    = "a dragon-type pokemon"
			)

			ctx := context.Background()
			cache := newCache(2)

			cache.Lookup(ctx, resourceA).Hydrate(ctx, valueA)
			cache.Lookup(ctx, resourceB).Hydrate(ctx, valueB)

			lookup := cache.Lookup(ctx, resourceA)
			defer lookup.Close(ctx)

			if v, ok := lookup.Value(ctx); !ok || v != valueA {
				t.Errorf(`want lookup.Value() to return (%q, true); got (%v, %T)`, valueA, v, ok)
			}
		},
	)

	// the eviction strategy is not strongly defined, so this is how we can test
	// that cache evictions are occurring
	t.Run(
		"for a cache of size N, hydrating M (M > N) values means M-N cache lookups for the same resources miss",
		func(t *testing.T) {
			t.Parallel()

			const dummyValue = "a ???-type pokemon"
			resourceName := func(i int) string {
				return fmt.Sprintf("https://pokapi.co/api/v2/pokemon/%d", i)
			}

			for _, entry := range []struct{ N, M int }{
				{N: 1, M: 2},
				{N: 10, M: 20},
				{N: 5, M: 15},
				{N: 50, M: 150},
			} {
				m, n, wantMisses := entry.M, entry.N, entry.M-entry.N

				t.Run(
					fmt.Sprintf(
						"filling %d values in a cache of size %d causes %d misses",
						m, n, wantMisses,
					),
					func(t *testing.T) {
						t.Parallel()

						ctx := context.Background()
						cache := newCache(entry.N)

						for i := range m {
							cache.Lookup(ctx, resourceName(i)).Hydrate(ctx, dummyValue)
						}

						var gotMisses int
						for i := range m {
							lookup := cache.Lookup(ctx, resourceName(i))
							if _, ok := lookup.Value(ctx); !ok {
								gotMisses++
							}
							lookup.Close(ctx)
						}

						if wantMisses != gotMisses {
							t.Errorf("want %d cache misses; got %d", wantMisses, gotMisses)
						}
					},
				)
			}
		},
	)

	// the cache may be used by multiple concurrent goroutines
	// note: code that could fail this test may not always. be sure to enable the
	// race detector and run the tests with -test.count.
	t.Run(
		"cache supports concurrent lookups on the same resource - only the first cache lookup should miss",
		func(t *testing.T) {
			t.Parallel()
			const (
				resource = "https://pokapi.co/api/v2/pokemon/tinkaton"
				value    = "a fairy-type pokemon"
			)

			ctx := context.Background()
			c := newCache(1)

			var (
				wg        sync.WaitGroup
				missesMux sync.Mutex
				misses    int
			)

			for range 10 {
				wg.Add(1)

				go func() {
					defer wg.Done()

					lookup := c.Lookup(ctx, resource)
					defer lookup.Close(ctx)

					if _, ok := lookup.Value(ctx); !ok {
						lookup.Hydrate(ctx, value)

						missesMux.Lock()
						misses++
						missesMux.Unlock()
					}
				}()
			}

			wg.Wait()

			if misses != 1 {
				t.Errorf("wanted only the first cache lookup to miss; got %d misses", misses)
			}
		},
	)
}

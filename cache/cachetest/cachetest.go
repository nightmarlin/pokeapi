// Package cachetest is a test suite for [pokeapi.Cache] implementations. Simply
// call [cachetest.TestCache] from your cache's test file.
//
// As it imports package testing, cachetest should not be used in normal
// application code.
package cachetest

import (
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

			cache := newCache(1)

			lookup := cache.Lookup("https://pokapi.co/api/v2/pokemon/gumshoos")
			defer lookup.Close()

			if v, ok := lookup.Value(); ok || v != nil {
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

			cache := newCache(1)

			cache.Lookup(resource).Hydrate(value)

			lookup := cache.Lookup(resource)
			defer lookup.Close()

			if v, ok := lookup.Value(); !ok || value != v {
				t.Errorf(`want lookup.Value() to return (%q, true); got (%v, %T)`, value, value, ok)
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

						cache := newCache(entry.N)

						for i := range entry.N {
							cache.Lookup(resource(i)).Hydrate(i)
						}

						var misses int
						for i := range entry.N {
							lookup := cache.Lookup(resource(i))
							if _, ok := lookup.Value(); !ok {
								misses++
							}
							lookup.Close()
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
			cache := newCache(1)

			cache.Lookup(resource).Hydrate(firstValue)
			cache.Lookup(resource).Hydrate(secondValue)

			lookup := cache.Lookup(resource)
			defer lookup.Close()

			if v, ok := lookup.Value(); !ok || v != secondValue {
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

			cache := newCache(2)

			cache.Lookup(resourceA).Hydrate(valueA)
			cache.Lookup(resourceB).Hydrate(valueB)

			lookup := cache.Lookup(resourceA)
			defer lookup.Close()

			if v, ok := lookup.Value(); !ok || v != valueA {
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

						cache := newCache(entry.N)

						for i := range m {
							cache.Lookup(resourceName(i)).Hydrate(dummyValue)
						}

						var gotMisses int
						for i := range m {
							lookup := cache.Lookup(resourceName(i))
							if _, ok := lookup.Value(); !ok {
								gotMisses++
							}
							lookup.Close()
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
	t.Run(
		"cache supports concurrent lookups on the same resource - only the first cache lookup should miss",
		func(t *testing.T) {
			t.Parallel()
			const (
				resource = "https://pokapi.co/api/v2/pokemon/tinkaton"
				value    = "a fairy-type pokemon"
			)

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

					lookup := c.Lookup(resource)
					defer lookup.Close()

					if _, ok := lookup.Value(); !ok {
						lookup.Hydrate(value)

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

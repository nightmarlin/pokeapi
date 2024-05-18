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

// mockLoader returns a pokeapi.CacheLoader and a function that returns how many
// times that pokeapi.CacheLoader has been called.
func mockLoader(v any, err error) (_ pokeapi.CacheLoader, callCount func() int) {
	var (
		ccMux       sync.RWMutex
		callCounter int
	)

	return func(context.Context) (any, error) {
			defer ccMux.Unlock()
			ccMux.Lock()
			callCounter += 1
			return v, err
		},
		func() int {
			defer ccMux.RUnlock()
			ccMux.RLock()
			return callCounter
		}
}

func TestCache[C pokeapi.Cache](t *testing.T, newCache NewCacheFn[C]) {
	t.Run(
		"lookup on an empty cache misses",
		func(t *testing.T) {
			t.Parallel()

			const (
				value = "a normal type pokemon"
			)

			var (
				ctx      = context.Background()
				cache    = newCache(1)
				ml, cc   = mockLoader(value, nil)
				res, err = cache.Lookup(ctx, "https://pokeapi.co/api/v2/pokemon/gumshoos", ml)
			)

			if err != nil {
				t.Errorf("want no error; got %v", err)
			}
			if calls := cc(); calls != 1 {
				t.Errorf("want loader to be called 1 time; got %d times", calls)
			}
			if res != value {
				t.Errorf("want returned value %v; got %v", value, res)
			}
		},
	)

	t.Run(
		"lookup after load returns value without another load",
		func(t *testing.T) {
			t.Parallel()

			const (
				resource = "https://pokeapi.co/api/v2/pokemon/yamask"
				value    = "a ghost-type pokemon"
			)

			var (
				ctx    = context.Background()
				cache  = newCache(1)
				ml, cc = mockLoader(value, nil)

				_, _     = cache.Lookup(ctx, resource, ml)
				res, err = cache.Lookup(ctx, resource, ml)
			)

			if err != nil {
				t.Errorf("want no error; got %v", err)
			}
			if calls := cc(); calls != 1 {
				t.Errorf("want loader to be called 1 time; got %d times", calls)
			}
			if res != value {
				t.Errorf("want returned value %v; got %v", value, res)
			}
		},
	)

	t.Run(
		"cached value is unchanged by other lookups",
		func(t *testing.T) {
			t.Parallel()

			const (
				resourceA = "https://pokeapi.co/api/v2/pokemon/serperior"
				valueA    = "a grass-type pokemon"
				resourceB = "https://pokeapi.co/api/v2/pokemon/dreepy"
				valueB    = "a dragon-type pokemon"
			)

			var (
				ctx   = context.Background()
				cache = newCache(2)

				mlA, ccA = mockLoader(valueA, nil)
				mlB, ccB = mockLoader(valueB, nil)
			)

			_, _ = cache.Lookup(ctx, resourceA, mlA)
			_, _ = cache.Lookup(ctx, resourceB, mlB)

			res, _ := cache.Lookup(ctx, resourceA, mlA)

			if res != valueA {
				t.Errorf(`want returned value for resource-a to be %q; got %v`, valueA, res)
			}
			if callsA := ccA(); callsA != 1 {
				t.Errorf("want resource-a loader to be called 1 time; got %d times", callsA)
			}
			if callsB := ccB(); callsB != 1 {
				t.Errorf("want resource-b loader to be called 1 time; got %d times", callsB)
			}
		},
	)

	t.Run(
		"for a cache of size N, lookup of N resources means no subsequent lookups for those resources miss",
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
							ml, _ := mockLoader(i, nil)
							_, _ = cache.Lookup(ctx, resource(i), ml)
						}

						var misses int
						for i := range entry.N {
							ml, cc := mockLoader(i, nil)
							_, _ = cache.Lookup(ctx, resource(i), ml)

							misses += cc()
						}

						if misses != 0 {
							t.Errorf("want 0 cache misses; got %d", misses)
						}
					},
				)
			}
		},
	)

	// the eviction strategy is not strongly defined, so this is how we can test
	// that cache evictions are correctly occurring, and that caches don't cache
	// error responses.
	t.Run(
		"for a cache of size N, hydrating M (M > N) values means M-N cache lookups for the same resources miss",
		func(t *testing.T) {
			t.Parallel()

			const dummyValue = "a ???-type pokemon"
			resourceName := func(i int) string {
				return fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%d", i)
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

						var (
							ctx   = context.Background()
							cache = newCache(entry.N)
							ml, _ = mockLoader(dummyValue, nil)
						)

						for i := range m {
							_, _ = cache.Lookup(ctx, resourceName(i), ml)
						}

						// reset loader: fail on load to prevent more values being cached
						ml, cc := mockLoader(nil, fmt.Errorf("brocken"))

						for i := range m {
							_, _ = cache.Lookup(ctx, resourceName(i), ml)
						}

						if gotMisses := cc(); wantMisses != gotMisses {
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
		"cache supports concurrent lookups on the same resource - only the one cache lookup should miss",
		func(t *testing.T) {
			t.Parallel()
			const (
				resource = "https://pokeapi.co/api/v2/pokemon/tinkaton"
				value    = "a fairy-type pokemon"
			)

			var (
				ctx = context.Background()
				c   = newCache(1)

				wg     sync.WaitGroup
				ml, cc = mockLoader(value, nil)
			)

			for range 100 {
				wg.Add(1)

				go func() {
					defer wg.Done()
					_, _ = c.Lookup(ctx, resource, ml)
				}()
			}

			wg.Wait()

			if misses := cc(); misses != 1 {
				t.Errorf("wanted only the first cache lookup to miss; got %d misses", misses)
			}
		},
	)
}

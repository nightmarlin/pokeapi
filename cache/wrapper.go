package cache

import (
	"context"
	"sync"

	"github.com/nightmarlin/pokeapi"
)

// The Wrapper cache wraps a standard get/put cache and converts it to the
// transactional pokeapi.Cache interface. It's useful if you want to make use of
// the pokeapi.Client's concurrency guarantees.
type Wrapper struct {
	getFn func(ctx context.Context, url string) (any, bool)
	putFn func(ctx context.Context, url string, value any)

	ongoing sync.Map
}

// NewWrapper accepts the Get and Put (or equivalent) method references of the
// cache it wraps and returns a pokeapi.Cache that loads and stores values
// from/to that cache.
//
// Example usage:
//
//	r := NewRedisCache(redisConn, defaultTTL)
//	c := pokeapi.Client(
//		&pokeapi.NewClientOpts{
//			Cache: cache.NewWrapper(r.Get, r.Put),
//		}
//	)
//
// It's worth noting that the url passed will be the raw url, query parameters
// included. Advanced cache implementations may choose to make use of this.
func NewWrapper(
	getFn func(ctx context.Context, url string) (any, bool),
	putFn func(ctx context.Context, url string, value any),
) *Wrapper {
	return &Wrapper{getFn: getFn, putFn: putFn}
}

func (w *Wrapper) Lookup(ctx context.Context, url string) pokeapi.CacheLookup {
	for {
		// spin: attempt to acquire lock on individual url.
		// todo: reduce spins?

		_, urlIsLocked := w.ongoing.LoadOrStore(url, struct{}{})
		if !urlIsLocked {
			break
		}
	}

	v, hasValue := w.getFn(ctx, url)
	return &wrapperCacheLookup{
		value:     v,
		hasValue:  hasValue,
		putFn:     func(ctx context.Context, value any) { w.putFn(ctx, url, value) },
		cleanupFn: func() { w.ongoing.Delete(url) },
	}
}

type wrapperCacheLookup struct {
	mux  sync.RWMutex
	once sync.Once

	value    any
	hasValue bool

	putFn     func(ctx context.Context, value any)
	cleanupFn func()
}

func (w *wrapperCacheLookup) Value(_ context.Context) (_ any, ok bool) {
	defer w.mux.RUnlock()
	w.mux.RLock()
	return w.value, w.hasValue
}

func (w *wrapperCacheLookup) cleanup() {
	w.cleanupFn()
	w.value = nil
	w.hasValue = false
	w.putFn = nil
	w.cleanupFn = nil
}

func (w *wrapperCacheLookup) Hydrate(ctx context.Context, resource any) {
	defer w.mux.Unlock()
	w.mux.Lock()
	w.once.Do(
		func() {
			w.putFn(ctx, resource)
			w.cleanup()
		},
	)
}

func (w *wrapperCacheLookup) Close(_ context.Context) {
	defer w.mux.Unlock()
	w.mux.Lock()
	w.once.Do(w.cleanup)
}

package cache

import (
	"context"

	"golang.org/x/sync/singleflight"

	"github.com/nightmarlin/pokeapi"
)

// The Wrapper cache wraps a standard get/put cache and converts it to the
// transactional pokeapi.Cache interface. It's useful if you want to make use of
// the pokeapi.Client's concurrency guarantees.
type Wrapper struct {
	getFn func(ctx context.Context, url string) (any, bool)
	putFn func(ctx context.Context, url string, value any)

	ongoing singleflight.Group
}

// NewWrapper accepts the Get and Put (or equivalent) method references of the
// cache it wraps and returns a pokeapi.Cache that loads and stores values
// from/to that cache.
//
// Example usage:
//
//	r := NewRedisCache(redisConn, defaultTTL)
//	c := pokeapi.Client(
//		&pokeapi.ClientOpts{
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

func (w *Wrapper) Lookup(
	ctx context.Context,
	url string,
	loadOnMiss pokeapi.CacheLoader,
) (any, error) {
	res, err, _ := w.ongoing.Do(
		url,
		func() (any, error) {
			if v, ok := w.getFn(ctx, url); ok {
				return v, nil
			}

			v, err := loadOnMiss(ctx)
			if err != nil {
				return nil, err
			}

			w.putFn(ctx, url, v)
			return v, nil
		},
	)

	if err != nil {
		return nil, err
	}
	return res, nil
}

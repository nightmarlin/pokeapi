package iterator

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/nightmarlin/pokeapi"
)

type nextPageFn[R pokeapi.GettableAPIResource[T], T any] func(
	context.Context,
	*pokeapi.Client,
) (*pokeapi.Page[R, T], error)

// An Iterator allows for quick and easy iteration through the elements of every
// page of a resource list.
type Iterator[R pokeapi.GettableAPIResource[T], T any] struct {
	mux sync.Mutex

	ctx       context.Context
	cancelCTX context.CancelFunc
	client    *pokeapi.Client

	results     []R
	nextReadIDX int
	nextPageFn  nextPageFn[R, T]

	closeOnce sync.Once
}

// New creates a new Iterator for the provided pokeapi.ResourceName.
// Iterator.Stop should always be called on the returned Iterator once you are
// finished with it.
func New[R pokeapi.GettableAPIResource[T], T any](
	ctx context.Context,
	client *pokeapi.Client,
	resourceName pokeapi.ResourceName[R, T],
) *Iterator[R, T] {
	ctx, cancel := context.WithCancel(ctx)

	return &Iterator[R, T]{
		ctx:       ctx,
		cancelCTX: cancel,
		client:    client,
		nextPageFn: func(ctx context.Context, c *pokeapi.Client) (*pokeapi.Page[R, T], error) {
			return resourceName.List(ctx, c, nil)
		},
	}
}

func (i *Iterator[R, T]) Next() (*T, error) {
	defer i.mux.Unlock()
	i.mux.Lock()

	if i.nextPageFn == nil {
		return nil, pokeapi.ErrListExhausted
	}

	if len(i.results) <= i.nextReadIDX {
		// fetch next page

		np, err := i.nextPageFn(i.ctx, i.client)
		if err != nil {
			if errors.Is(err, pokeapi.ErrListExhausted) {
				i.close()
				return nil, pokeapi.ErrListExhausted
			}
			return nil, fmt.Errorf("fetching next page: %w", err)
		}

		i.nextPageFn = np.GetNext
		i.results = np.Results
		i.nextReadIDX = 0
	}

	if len(i.results) == 0 {
		i.close()
		return nil, pokeapi.ErrListExhausted
	}

	r, err := i.results[i.nextReadIDX].Get(i.ctx, i.client)
	if err != nil {
		return nil, fmt.Errorf("fetching next resource: %w", err)
	}
	i.nextReadIDX += 1
	return r, nil
}

func (i *Iterator[R, T]) close() {
	i.closeOnce.Do(
		func() {
			i.cancelCTX()
			i.nextPageFn = nil
			i.results = nil
			i.nextReadIDX = 0
			i.client = nil
			i.ctx = nil
			i.cancelCTX = nil
		},
	)
}

// Stop cleans up resources used by the Iterator.
func (i *Iterator[R, T]) Stop() {
	defer i.mux.Unlock()
	i.mux.Lock()
	i.close()
}

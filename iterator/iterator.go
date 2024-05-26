// Package iterator provides an Iterator type that manages page iteration for
// you - fetching the following pokeapi.GettableAPIResource and (if necessary)
// the next pokeapi.Page of resources.
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
// pokeapi.Page of a resource list.
type Iterator[R pokeapi.GettableAPIResource[T], T any] struct {
	mux sync.Mutex

	client *pokeapi.Client

	results     []R
	nextReadIDX int
	nextPageFn  nextPageFn[R, T]

	closeOnce sync.Once
}

// New creates a new Iterator for the provided pokeapi.ResourceName.
//
// Iterator.Stop should always be called on the returned Iterator once you are
// finished with it.
func New[R pokeapi.GettableAPIResource[T], T any](
	client *pokeapi.Client,
	resourceName pokeapi.ResourceName[R, T],
) *Iterator[R, T] {
	return &Iterator[R, T]{
		client: client,
		nextPageFn: func(ctx context.Context, c *pokeapi.Client) (*pokeapi.Page[R, T], error) {
			return resourceName.List(ctx, c, nil)
		},
	}
}

// NewFromPage creates a new Iterator that starts at the provided pokeapi.Page.
// Iteration will continue with pages of the size of the original pokeapi.Page.
//
// Iterator.Stop should always be called on the returned Iterator once you are
// finished with it.
func NewFromPage[R pokeapi.GettableAPIResource[T], T any](
	client *pokeapi.Client,
	page pokeapi.Page[R, T],
) *Iterator[R, T] {
	return &Iterator[R, T]{
		client:     client,
		results:    page.Results,
		nextPageFn: page.GetNext,
	}
}

// Next fetches the next page if the current one is empty or has been exhausted.
// It then fetches the next available pokeapi.GettableAPIResource and returns
// its value (or any error that may have occurred).
//
// Calling Next after Stop will return pokeapi.ErrListExhausted.
func (i *Iterator[R, T]) Next(ctx context.Context) (*T, error) {
	defer i.mux.Unlock()
	i.mux.Lock()

	if i.nextPageFn == nil {
		return nil, pokeapi.ErrListExhausted
	}

	if len(i.results) <= i.nextReadIDX {
		// fetch next page

		np, err := i.nextPageFn(ctx, i.client)
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

	r, err := i.results[i.nextReadIDX].Get(ctx, i.client)
	if err != nil {
		return nil, fmt.Errorf("fetching next resource: %w", err)
	}
	i.nextReadIDX += 1
	return r, nil
}

func (i *Iterator[R, T]) close() {
	i.closeOnce.Do(
		func() {
			i.nextPageFn = nil
			i.results = nil
			i.nextReadIDX = 0
			i.client = nil
		},
	)
}

// Stop cleans up resources used by the Iterator.
func (i *Iterator[R, T]) Stop() {
	defer i.mux.Unlock()
	i.mux.Lock()
	i.close()
}

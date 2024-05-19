//go:build goexperiment.rangefunc

package iterator

import (
	"context"
	"errors"
	"iter"

	"github.com/nightmarlin/pokeapi"
)

func NewSeq[R pokeapi.GettableAPIResource[T], T any](
	ctx context.Context,
	client *pokeapi.Client,
	resource pokeapi.ResourceName[R, T],
) iter.Seq2[*T, error] {
	i := New(ctx, client, resource)
	defer i.Stop()

	return func(yield func(*T, error) bool) {
		v, err := i.Next()
		if errors.Is(err, pokeapi.ErrListExhausted) {
			return
		}

		if !yield(v, err) {
			return
		}
	}
}

package iterator_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nightmarlin/pokeapi"
	"github.com/nightmarlin/pokeapi/iterator"
)

func ptr[T any](in T) *T { return &in }

func TestIterator(t *testing.T) {
	t.Parallel()

	stub := func(t *testing.T) (handle http.HandlerFunc, add func(path string, v any)) {
		m := make(map[string]any)

		return func(w http.ResponseWriter, r *http.Request) {
			lookupStr := r.URL.String()
			v, ok := m[lookupStr]
			if !ok {
				t.Logf("%s requested, but it couldn't be found", lookupStr)
				http.NotFound(w, r)
				return
			}

			vb, err := json.Marshal(v)
			if err != nil {
				t.Fatalf("failed to marshal response for %q (%v): %v", lookupStr, v, err)
			}

			if _, err := w.Write(vb); err != nil {
				t.Fatalf("failed to write response for %q: %v", lookupStr, err)
			}
		}, func(path string, v any) { m[path] = v }
	}

	t.Run(
		"converts not found on List call to ErrListExhausted",
		func(t *testing.T) {
			t.Parallel()

			var (
				ctx        = context.Background()
				handler, _ = stub(t)
				ts         = httptest.NewServer(handler)
				c          = pokeapi.NewClient(
					&pokeapi.ClientOpts{HTTPClient: ts.Client(), PokeAPIRoot: ts.URL},
				)
				i = iterator.New(ctx, c, pokeapi.PokemonResource)
			)
			t.Cleanup(ts.Close)
			t.Cleanup(i.Stop)

			v, err := i.Next()
			if !errors.Is(err, pokeapi.ErrListExhausted) || v != nil {
				t.Errorf("want (nil, ErrListExhausted); got (%v, %v)", v, err)
			}
		},
	)

	t.Run(
		"propagates normal errors",
		func(t *testing.T) {
			t.Parallel()

			var (
				ctx = context.Background()
				ts  = httptest.NewServer(
					http.HandlerFunc(
						func(w http.ResponseWriter, _ *http.Request) {
							http.Error(w, "brocken", http.StatusInternalServerError)
						},
					),
				)
				c = pokeapi.NewClient(&pokeapi.ClientOpts{HTTPClient: ts.Client(), PokeAPIRoot: ts.URL})
				i = iterator.New(ctx, c, pokeapi.PokemonResource)
			)
			t.Cleanup(ts.Close)
			t.Cleanup(i.Stop)

			v, err := i.Next()
			if !errors.Is(err, pokeapi.HTTPError{Code: http.StatusInternalServerError}) || v != nil {
				t.Errorf("want (nil, ErrListExhausted); got (%v, %v)", v, err)
			}
		},
	)

	t.Run(
		"propagates not found errors for individual resources (unlikely)",
		func(t *testing.T) {
			t.Parallel()

			var (
				ctx          = context.Background()
				handler, add = stub(t)
				ts           = httptest.NewServer(handler)
				c            = pokeapi.NewClient(
					&pokeapi.ClientOpts{HTTPClient: ts.Client(), PokeAPIRoot: ts.URL},
				)
				i = iterator.New(ctx, c, pokeapi.PokemonResource)
			)
			t.Cleanup(ts.Close)
			t.Cleanup(i.Stop)

			add(
				"/pokemon/",
				pokeapi.Page[pokeapi.NamedAPIResource[pokeapi.Pokemon], pokeapi.Pokemon]{
					Count:    1,
					Next:     nil,
					Previous: nil,
					Results: []pokeapi.NamedAPIResource[pokeapi.Pokemon]{
						{
							APIResource: pokeapi.APIResource[pokeapi.Pokemon]{
								URL: fmt.Sprintf("%s/pokemon/missingno", ts.URL),
							},
							Name: "missingno",
						},
					},
				},
			)

			v, err := i.Next()
			if !errors.Is(err, pokeapi.ErrNotFound) || v != nil {
				t.Errorf(
					"want get '/pokemon/missingno' to return (nil, ErrListExhausted); got (%v, %v)",
					v, err,
				)
			}

		},
	)

	t.Run(
		"correctly iterates across multiple pages",
		func(t *testing.T) {
			t.Parallel()

			var (
				ctx          = context.Background()
				handler, add = stub(t)
				ts           = httptest.NewServer(handler)
				c            = pokeapi.NewClient(
					&pokeapi.ClientOpts{HTTPClient: ts.Client(), PokeAPIRoot: ts.URL},
				)
				i = iterator.New(ctx, c, pokeapi.PokemonResource)
			)
			t.Cleanup(ts.Close)
			t.Cleanup(i.Stop)

			add(
				"/pokemon/",
				pokeapi.Page[pokeapi.NamedAPIResource[pokeapi.Pokemon], pokeapi.Pokemon]{
					Count:    2,
					Next:     ptr(fmt.Sprintf("%s/pokemon/?offset=1", ts.URL)),
					Previous: nil,
					Results: []pokeapi.NamedAPIResource[pokeapi.Pokemon]{
						{
							APIResource: pokeapi.APIResource[pokeapi.Pokemon]{
								URL: fmt.Sprintf("%s/pokemon/1", ts.URL),
							},
							Name: "bulbasaur",
						},
					},
				},
			)
			add(
				"/pokemon/?offset=1",
				pokeapi.Page[pokeapi.NamedAPIResource[pokeapi.Pokemon], pokeapi.Pokemon]{
					Count:    2,
					Next:     nil,
					Previous: nil,
					Results: []pokeapi.NamedAPIResource[pokeapi.Pokemon]{
						{
							APIResource: pokeapi.APIResource[pokeapi.Pokemon]{
								URL: fmt.Sprintf("%s/pokemon/2", ts.URL),
							},
							Name: "ivysaur",
						},
					},
				},
			)
			add(
				"/pokemon/1",
				pokeapi.Pokemon{
					NamedIdentifier: pokeapi.NamedIdentifier{
						Identifier: pokeapi.Identifier{ID: 1},
						Name:       "bulbasaur",
					},
				},
			)
			add(
				"/pokemon/2",
				pokeapi.Pokemon{
					NamedIdentifier: pokeapi.NamedIdentifier{
						Identifier: pokeapi.Identifier{ID: 2},
						Name:       "ivysaur",
					},
				},
			)

			v, err := i.Next()
			if v == nil || v.ID != 1 || err != nil {
				t.Errorf(
					"want get '/pokemon/1' to return (ID:1, nil); got (%v, %v)",
					v, err,
				)
			}

			v, err = i.Next()
			if v == nil || v.ID != 2 || err != nil {
				t.Errorf(
					"want get '/pokemon/2' to return (ID:2, nil); got (%v, %v)",
					v, err,
				)
			}

			v, err = i.Next()
			if !errors.Is(err, pokeapi.ErrListExhausted) || v != nil {
				t.Errorf("want list to be exhausted; got (%v, %v)", v, err)
			}
		},
	)

	t.Run(
		"calling Next() after Stop() returns ErrListExhausted, even if more elements were left",
		func(t *testing.T) {
			t.Parallel()

			var (
				ctx          = context.Background()
				handler, add = stub(t)
				ts           = httptest.NewServer(handler)
				c            = pokeapi.NewClient(
					&pokeapi.ClientOpts{HTTPClient: ts.Client(), PokeAPIRoot: ts.URL},
				)
				i = iterator.New(ctx, c, pokeapi.PokemonResource)
			)
			t.Cleanup(ts.Close)

			add(
				"/pokemon/",
				pokeapi.Page[pokeapi.NamedAPIResource[pokeapi.Pokemon], pokeapi.Pokemon]{
					Count:    2,
					Next:     nil,
					Previous: nil,
					Results: []pokeapi.NamedAPIResource[pokeapi.Pokemon]{
						{
							APIResource: pokeapi.APIResource[pokeapi.Pokemon]{
								URL: fmt.Sprintf("%s/pokemon/1", ts.URL),
							},
							Name: "bulbasaur",
						},
						{
							APIResource: pokeapi.APIResource[pokeapi.Pokemon]{
								URL: fmt.Sprintf("%s/pokemon/2", ts.URL),
							},
							Name: "ivysaur",
						},
					},
				},
			)
			add(
				"/pokemon/1",
				pokeapi.Pokemon{
					NamedIdentifier: pokeapi.NamedIdentifier{
						Identifier: pokeapi.Identifier{ID: 1},
						Name:       "bulbasaur",
					},
				},
			)

			v, err := i.Next()
			if v == nil || v.ID != 1 || err != nil {
				t.Errorf(
					"want get '/pokemon/1' to return (ID:1, nil); got (%v, %v)",
					v, err,
				)
			}

			i.Stop()

			v, err = i.Next()
			if !errors.Is(err, pokeapi.ErrListExhausted) || v != nil {
				t.Errorf("want (nil, ErrListExhausted); got (%v, %v)", v, err)
			}
		},
	)
}

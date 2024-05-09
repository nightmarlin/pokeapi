// Package cachetest is a test suite for [pokeapi.Cache] implementations. Simply
// call [cachetest.TestCache] from your cache's test file.
//
// As it imports package testing, cachetest should not be used in normal
// application code.
package cachetest

import (
	"testing"

	"github.com/nightmarlin/pokeapi"
)

type NewCacheFn[T pokeapi.Cache] func(size int) T

func TestCache[C pokeapi.Cache](t *testing.T, newCache NewCacheFn[C]) {
	t.Run(
		"",
		func(t *testing.T) {
			_ = newCache(1) // todo
		},
	)
}

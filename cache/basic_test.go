package cache_test

import (
	"testing"

	"github.com/nightmarlin/pokeapi/cache"
	"github.com/nightmarlin/pokeapi/cache/cachetest"
)

func TestBasic(t *testing.T) {
	cachetest.TestCache(t, cache.NewBasic)
}

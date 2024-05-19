package pokeapi_test

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/nightmarlin/pokeapi"
	"github.com/nightmarlin/pokeapi/cache"
)

func ExampleClient() {
	var (
		ctx, cancel = signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
		printf      = func(f string, a ...any) { _, _ = fmt.Fprintf(os.Stderr, f, a...) }
	)
	defer cancel()

	c := pokeapi.NewClient(&pokeapi.ClientOpts{Cache: cache.NewLRU(nil)})

	berries, err := c.ListBerries(ctx, nil)
	if err != nil {
		printf("failed to list all the berries: %v", err)
		return
	}

	printf("there are %d known berries", berries.Count)

	firstBerry, err := berries.Results[0].Get(ctx, c)
	if err != nil {
		printf("failed to fetch the first berry: %v", err)
		return
	}

	printf("the first berry is the %s berry", firstBerry.Name)

	reFetchedBerry, err := c.GetBerry(ctx, firstBerry.Ident()) // should hit the cache
	if err != nil {
		printf("failed to re-fetch the first berry: %v", err)
		return
	}

	printf("no really, it's the %s berry", reFetchedBerry.Name)

	reReFetchedBerry, err := c.GetBerry(ctx, firstBerry.Ident()) // definitely hits the cache!
	if err != nil {
		printf("failed to re-fetch the first berry: %v", err)
		return
	}

	printf("i am 100%% certain it's the %s berry", reReFetchedBerry.Name)

	item, err := firstBerry.Item.Get(ctx, c)
	if err != nil {
		printf("failed to fetch the item corresponding to the first berry: %v", err)
		return
	}

	printf("it typically costs %d$poké", item.Cost)

	category, err := item.Category.Get(ctx, c)
	if err != nil {
		printf("failed to fetch the item category for the first berry: %v", err)
		return
	}

	printf("and it goes in the %s pocket", category.Pocket.Name)
}

// Command the-first-berry provides a quick example of how the pokeapi.Client
// is created and used, and the benefit of using a pokeapi.Cache to cache
// results!
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	"github.com/nightmarlin/pokeapi"
	"github.com/nightmarlin/pokeapi/cache"
)

func logErr(log *slog.Logger) func(error, string) {
	return func(err error, msg string) { log.Error(msg, slog.String("error", err.Error())) }
}
func logInfof(log *slog.Logger) func(string, ...any) {
	return func(format string, args ...any) { log.Info(fmt.Sprintf(format, args...)) }
}

func main() {
	var (
		ctx, cancel = signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
		log         = slog.New(slog.NewTextHandler(os.Stderr, nil))
		e, i        = logErr(log), logInfof(log)
	)

	defer cancel()

	c := pokeapi.NewClient(&pokeapi.ClientOpts{Cache: cache.NewLRU(nil)})

	berries, err := c.ListBerries(ctx, nil)
	if err != nil {
		e(err, "failed to list all the berries")
		return
	}

	i("there are %d known berries", berries.Count)

	firstBerry, err := berries.Results[0].Get(ctx, c)
	if err != nil {
		e(err, "failed to fetch the first berry")
		return
	}

	i("the first berry is the %s berry", firstBerry.Name)

	reFetchedBerry, err := c.GetBerry(ctx, firstBerry.Ident()) // should hit the cache
	if err != nil {
		e(err, "failed to re-fetch the first berry")
		return
	}

	i("no really, it's the %s berry", reFetchedBerry.Name)

	reReFetchedBerry, err := c.GetBerry(ctx, firstBerry.Ident()) // definitely hits the cache!
	if err != nil {
		e(err, "failed to re-fetch the first berry")
		return
	}

	i("i am 100%% certain it's the %s berry", reReFetchedBerry.Name)

	item, err := firstBerry.Item.Get(ctx, c)
	if err != nil {
		e(err, "failed to fetch the item corresponding to the first berry")
		return
	}

	i("it typically costs %d$poké", item.Cost)

	category, err := item.Category.Get(ctx, c)
	if err != nil {
		e(err, "failed to fetch the item category for the first berry")
		return
	}

	i("and it goes in the %s pocket", category.Pocket.Name)
}

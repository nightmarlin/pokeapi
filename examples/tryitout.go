package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/nightmarlin/pokeapi"
	"github.com/nightmarlin/pokeapi/cache"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	c := pokeapi.NewClient(&pokeapi.NewClientOpts{Cache: cache.NewBasic(50)})

	berries, err := c.ListBerries(ctx, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("there are", berries.Count, "known berries")

	firstBerry, err := berries.Results[0].Get(ctx, c)
	if err != nil {
		panic(err)
	}

	fmt.Println("the first berry is the", firstBerry.Name, "berry")

	reFetchedBerry, err := c.GetBerry(ctx, firstBerry.Ident())
	if err != nil {
		panic(err)
	}

	fmt.Println("no really, it's the", reFetchedBerry.Name, "berry")

	reReFetchedBerry, err := c.GetBerry(ctx, firstBerry.Ident())
	if err != nil {
		panic(err)
	}

	fmt.Println("i am 100% certain it's the", reReFetchedBerry.Name, "berry")

	item, err := firstBerry.Item.Get(ctx, c)
	if err != nil {
		panic(err)
	}

	fmt.Printf("it typically costs %d$poké\n", item.Cost)

	category, err := item.Category.Get(ctx, c)
	if err != nil {
		panic(err)
	}

	fmt.Println("and it goes in the", category.Pocket.Name, "pocket")
}

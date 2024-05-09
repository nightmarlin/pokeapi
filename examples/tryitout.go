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

	item, err := firstBerry.Item.Get(ctx, c)
	if err != nil {
		panic(err)
	}

	fmt.Println("and it goes in the", item.Category.Name, "pocket")
}

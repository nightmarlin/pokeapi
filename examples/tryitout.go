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

	fmt.Println("retrieved", berries.Count, "results")

	firstBerry, err := berries.Results[0].Get(ctx, c)
	if err != nil {
		panic(err)
	}

	fmt.Println("the first berry is:", firstBerry.Name)

	alsoTheFirstBerry, err := c.GetBerry(ctx, firstBerry.Ident())
	if err != nil {
		panic(err)
	}

	fmt.Println("no, really, it's", alsoTheFirstBerry.Name)
}

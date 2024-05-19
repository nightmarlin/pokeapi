// Command iterator demonstrates how to use the [iterator.Iterator] with a
// [pokeapi.Client] and [pokeapi.ResourceName] to easily iterate through the
// full list of resources.
//
// When built using GOEXPERIMENT=rangefunc, iterator also provides the
// [iterator.NewSeq] function, which simplifies this usage even further. When
// the rangefunc experiment is officially added to go, this command will be
// updated to use it.
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"text/template"

	"github.com/nightmarlin/pokeapi"
	"github.com/nightmarlin/pokeapi/cache"
	"github.com/nightmarlin/pokeapi/iterator"
)

func main() {
	var (
		ctx, cancel = signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
		c           = pokeapi.NewClient(&pokeapi.ClientOpts{Cache: cache.NewLRU(nil)})

		dexIter  = iterator.New(ctx, c, pokeapi.PokedexResource)
		dexCount = 0
	)
	defer cancel()
	defer dexIter.Stop()

	fmt.Println("let's learn about pokedexes...")

	for {
		dex, err := dexIter.Next()
		if err != nil {
			if errors.Is(err, pokeapi.ErrListExhausted) {
				break
			}
			_, _ = fmt.Fprintf(os.Stderr, "failed to get next pokedex: %v", err)
			return
		}

		if err := detailsTemplate.Execute(os.Stdout, dexToTemplateArgs(dex)); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "failed to write details for dex %v: %v", dex.Name, err)
			return
		}

		dexCount++
	}

	fmt.Println("and that's all", dexCount, "pokedexes!")
}

func dexToTemplateArgs(dex *pokeapi.Pokedex) templateArgs {
	rc := "no specific"
	if dex.Region != nil {
		rc = fmt.Sprintf("the %s", dex.Region.Name)
	}

	return templateArgs{
		Name:         dex.Name,
		RegionCopy:   rc,
		EntryCount:   len(dex.PokemonEntries),
		IsMainSeries: dex.IsMainSeries,
		FirstEntry:   dex.PokemonEntries[0].PokemonSpecies.Name,
		LastEntry:    dex.PokemonEntries[len(dex.PokemonEntries)-1].PokemonSpecies.Name,
	}
}

type templateArgs struct {
	Name         string
	RegionCopy   string
	EntryCount   int
	IsMainSeries bool

	FirstEntry string
	LastEntry  string
}

var detailsTemplate = template.Must(template.New("dex-details").Parse(details))

const details = `the {{.Name}} pokédex applies to {{.RegionCopy}} region
	it contains {{.EntryCount}} entries, and is{{ if not .IsMainSeries}}'t {{end}} present in the main series games
	its first entry is {{.FirstEntry}} & its last entry is {{.LastEntry}}
`

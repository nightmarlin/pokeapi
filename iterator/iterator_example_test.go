package iterator_test

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

func ExampleIterator() {
	var (
		ctx, cancel = signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
		c           = pokeapi.NewClient(&pokeapi.ClientOpts{Cache: cache.NewLRU(nil)})

		dexIter  = iterator.New(c, pokeapi.PokedexResource)
		dexCount = 0
	)
	defer cancel()
	defer dexIter.Stop()

	fmt.Println("let's learn about pokedexes...")

	for {
		dex, err := dexIter.Next(ctx)
		if err != nil {
			if errors.Is(err, pokeapi.ErrListExhausted) {
				break
			}
			_, _ = fmt.Fprintf(os.Stderr, "failed to get next pokedex: %v", err)
			return
		}

		if err := detailsTemplate.Execute(os.Stdout, dex); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "failed to write details for dex %v: %v", dex.Name, err)
			return
		}

		dexCount++
	}

	fmt.Println("and that's all", dexCount, "pokedexes!")
}

var detailsTemplate = template.Must(
	template.
		New("dex-details").
		Funcs(template.FuncMap{"subOne": func(i int) int { return i - 1 }}).
		Parse(
			`the {{.Name}} pokédex applies to {{if (eq .Region nil)}}no specific{{else}}the {{.Region.Name}}{{end}} region:
	- there are {{len .PokemonEntries}} entries
	- it is{{if not .IsMainSeries}}n't{{end}} present in the main series games
	- its first entry is {{(index .PokemonEntries 0).PokemonSpecies.Name}}
	- its last entry is {{(index .PokemonEntries (subOne (len .PokemonEntries))).PokemonSpecies.Name}}
`,
		),
)

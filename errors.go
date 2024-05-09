package pokeapi

import (
	"fmt"
)

var (
	// ErrListExhausted is returned by Page.GetNext and Page.GetPrevious when
	// there are no more results in that direction.
	ErrListExhausted = fmt.Errorf("no more pages to fetch")

	// ErrNotFound is returned if you request a resource path that doesn't exist.
	ErrNotFound = fmt.Errorf("resource not found")
)

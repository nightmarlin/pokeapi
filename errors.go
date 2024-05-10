package pokeapi

import (
	"fmt"
)

var (
	// ErrListExhausted is returned by Page.GetNext and Page.GetPrevious when
	// there are no more results in that direction.
	ErrListExhausted = fmt.Errorf("no more pages to fetch")
)

// HTTPErr represents the error returned by a failed HTTP request.
type HTTPErr struct {
	Status     string
	StatusCode int // The 4xx / 5xx error code
}

func (e HTTPErr) Error() string {
	return fmt.Sprintf("%d %s", e.StatusCode, e.Status)
}

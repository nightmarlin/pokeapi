package pokeapi

import (
	"fmt"
	"net/http"
)

var (
	// ErrListExhausted is returned by Page.GetNext and Page.GetPrevious when
	// there are no more results in that direction.
	ErrListExhausted = fmt.Errorf("no more pages to fetch")

	// ErrNotFound is the error returned when attempting to retrieve a resource
	// that does not exist.
	ErrNotFound = HTTPError{Code: 404}
)

// HTTPError represents an error returned by a failed HTTP request. As a special
// case, 404 Not Found returns ErrNotFound instead.
type HTTPError struct{ Code int }

func (e HTTPError) Error() string { return fmt.Sprintf("%d %s", e.Code, http.StatusText(e.Code)) }

func NewHTTPError(resp *http.Response) error { return HTTPError{Code: resp.StatusCode} }

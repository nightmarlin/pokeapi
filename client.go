package pokeapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

//go:generate go run cmd/gettergen/gettergen.go -- $GOFILE "getters.gen.go"

// The DefaultPokeAPIRoot is the standard URL for PokéAPI. An alternative URL
// can be provided via [ClientOpts.PokeAPIRoot] for use with alternative
// builds of the API.
const DefaultPokeAPIRoot = `https://pokeapi.co/api/v2`

// The Client wraps a http.Client and a Cache to perform requests to PokéAPI.
//
// All methods of the form `Get*` accept the id or name of the resource (unless
// otherwise stated) & return one instance of that resource.
//
// All methods of the form `List*` will return the first Page of results, and
// accept an optional ListOptions parameter to permit you to start iteration
// wherever you like. This parameter may always be nil to start iteration from
// the beginning.
//
// Return types are exact as possible. Pointer types are used to represent
// "optional" fields. Slice fields are is always potentially empty.
type Client struct {
	client      *http.Client
	cache       Cache
	pokeAPIRoot string
}

type ClientOpts struct {
	HTTPClient  *http.Client // Set the HTTP client to use when making lookups.
	Cache       Cache        // Provide a Cache for use in lookups.
	PokeAPIRoot string       // Change the base PokéAPI URL to make lookups to.
}

// A CacheLoader is called on cache misses to retrieve the value of the resource
// from an external source.
type CacheLoader func(context.Context) (any, error)

// A Cache allows the Client to Lookup a URL and retrieve the corresponding
// resource if it has been fetched before. `loadOnMiss` should only be called if
// the cache does not contain a value for the requested `url`. It is
// recommended (but not required) that concurrent lookups for the same `url`
// only make one call to a `loadOnMiss` between them.
type Cache interface {
	Lookup(ctx context.Context, url string, loadOnMiss CacheLoader) (any, error)
}

// NewClient creates and returns a new Client with the provided ClientOpts
// applied. It is safe to use as NewClient(nil), but you are expected to do your
// own caching.
func NewClient(opts *ClientOpts) *Client {
	c := Client{client: http.DefaultClient, pokeAPIRoot: DefaultPokeAPIRoot, cache: noCache{}}

	if opts != nil {
		if opts.HTTPClient != nil {
			c.client = opts.HTTPClient
		}

		if opts.Cache != nil {
			c.cache = opts.Cache
		}

		if opts.PokeAPIRoot != "" {
			c.pokeAPIRoot = trimSlash(opts.PokeAPIRoot)
		}
	}

	return &c
}

// A Resource is the kebab-case name for a PokéAPI endpoint. Resources can
// typically be called with
//
//	GET {root}/resource`
//
// to return a list of links to instances of the resource or
//
//	`GET {root}/resource/{id or name}`
//
// to get a single instance of that Resource. Exceptions to this rule are
// documented.
type Resource string

func trimSlash[S ~string](s S) string { return strings.Trim(string(s), "/") }

func (c *Client) listURL(resource Resource) string {
	return fmt.Sprintf("%s/%s/", c.pokeAPIRoot, trimSlash(resource))
}
func (c *Client) getURL(resource Resource, name string) string {
	return fmt.Sprintf("%s/%s/%s/", c.pokeAPIRoot, trimSlash(resource), trimSlash(name))
}

// An Identifier is embedded into all retrievable resources. It makes it easy to
// convert the numbered ID field to a string for use in api calls.
//
// A resource directly embedding an Identifier will have an unnamed get/list
// client function pair generated for it by gettergen.
type Identifier struct {
	ID int `json:"id"`
}

// Ident returns the api id for this resource as a string for use in api calls.
func (id Identifier) Ident() string { return strconv.Itoa(id.ID) }

// An APIResource represents the indirect link to another resource. Use the Get
// method to retrieve the full resource being referred to.
type APIResource[T any] struct {
	URL string `json:"url"`
}

// Get uses the passed Client to retrieve the full details of the given APIResource.
func (r APIResource[T]) Get(ctx context.Context, client *Client) (*T, error) {
	return do[*T](ctx, client, r.URL, nil)
}

// A NamedIdentifier is embedded into resources that are named.
//
// A resource directly embedding a NamedIdentifier will have a named get/list
// client function pair generated for it by gettergen.
type NamedIdentifier struct {
	Identifier
	Name string `json:"name"`
}

// A NamedAPIResource is similar to an APIResource, but it provides an
// additional human-readable Name. Use the [APIResource.Get] method to retrieve
// the full resource being referred to.
type NamedAPIResource[T any] struct {
	APIResource[T]
	Name string `json:"name"`
}

// ListOptions are available on all List* endpoints, allowing you to set up your
// own pagination start point.
type ListOptions struct {
	Limit  int
	Offset int
}

// urlValues converts the ListOptions to its corresponding url.Values. It is
// safe to call on nil ListOptions.
func (lo *ListOptions) urlValues() url.Values {
	if lo == nil {
		return nil
	}

	v := url.Values{}
	if lo.Offset != 0 {
		v.Set("offset", strconv.Itoa(lo.Offset))
	}
	if lo.Limit != 0 {
		v.Set("limit", strconv.Itoa(lo.Limit))
	}
	return v
}

// A Page represents a list of APIResource s or NamedAPIResource s. It also
// includes information on the total number of resources in the result set, and
// how to view the Next & Previous Page s.
type Page[R APIResource[T] | NamedAPIResource[T], T any] struct {
	Count    int     `json:"count"`    // The total number of resources available from this API.
	Next     *string `json:"next"`     // The URL for the next page in the list.
	Previous *string `json:"previous"` // The URL for the previous page in the list.
	Results  []R     `json:"results"`
}

// GetNext retrieves the Page at Page.Next. If there is no next page,
// ErrListExhausted is returned.
func (p *Page[R, T]) GetNext(ctx context.Context, client *Client) (*Page[R, T], error) {
	if p.Next == nil {
		return nil, ErrListExhausted
	}

	return do[*Page[R, T]](ctx, client, *p.Next, nil)
}

// GetPrevious retrieves the Page at Page.Previous. If there is no previous
// page, ErrListExhausted is returned.
func (p *Page[R, T]) GetPrevious(ctx context.Context, client *Client) (*Page[R, T], error) {
	if p.Previous == nil {
		return nil, ErrListExhausted
	}

	return do[*Page[R, T]](ctx, client, *p.Previous, nil)
}

// The noCache is the default Cache implementation used by a Client. While it is
// valid for use, it does not perform any actual caching.
type noCache struct{}

func (noCache) Lookup(
	ctx context.Context,
	_ string,
	loader CacheLoader,
) (any, error) {
	return loader(ctx)
}

// zero creates and returns the zero value of T.
func zero[T any]() (z T) { return }

// do performs a type-safe http GET operation, using the Client's cache &
// http.Client.
func do[T any](ctx context.Context, c *Client, path string, values url.Values) (T, error) {
	res, err := c.cache.Lookup(
		ctx,
		path,
		func(ctx context.Context) (any, error) {
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
			if err != nil {
				return nil, fmt.Errorf("creating request: %w", err)
			}
			req.Header.Set("Accept", "application/json")

			if len(values) != 0 {
				qry := req.URL.Query()
				for field, val := range values {
					qry[field] = val
				}
				req.URL.RawQuery = qry.Encode()
			}

			resp, err := c.client.Do(req)
			if err != nil {
				return nil, fmt.Errorf("performing request: %w", err)
			}
			defer func() { _ = resp.Body.Close() }()
			if resp.StatusCode != http.StatusOK {
				return nil, NewHTTPError(resp)
			}

			var res T
			if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
				return nil, fmt.Errorf("decoding json response: %w", err)
			}
			return res, nil
		},
	)

	if err != nil {
		return zero[T](), err
	}
	return res.(T), nil
}

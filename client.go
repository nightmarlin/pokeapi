package pokeapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

//go:generate go run cmd/gettergen/gettergen.go -- "getters.gen.go"

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
// accept an optional ListOpts parameter to permit you to start iteration
// wherever you like. This parameter may always be nil to start iteration from
// the beginning.
//
// Return types are exact as possible. Pointer types are used to represent
// "optional" fields. Slice fields are always potentially empty.
type Client struct {
	client      *http.Client
	cache       Cache
	pokeAPIRoot string
}

type ClientOpts struct {
	HTTPClient  *http.Client // Set the HTTP client to use when making lookups. Can be used to add tracing.
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

// A ResourceName is the kebab-case name for a PokéAPI endpoint. Resources can
// always be called with
//
//	GET {root}/resource
//
// to return a Page of APIResource(s) pointing to instances of the resource or
//
//	GET {root}/resource/{id or name}
//
// to get a single instance of that resource. The helper methods List and Get
// correspond to these api calls.
type ResourceName[R GettableAPIResource[T], T any] string

func (rn ResourceName[R, T]) String() string { return string(rn) }
func (rn ResourceName[R, T]) isResource()    {}

// Get allows for the retrieval of a single instance of the desired resource.
func (rn ResourceName[R, T]) Get(ctx context.Context, c *Client, ident string) (*T, error) {
	return do[*T](ctx, c, c.getURL(rn, ident), nil)
}

func (rn ResourceName[R, T]) List(
	ctx context.Context,
	c *Client,
	opts *ListOpts,
) (*Page[R, T], error) {
	return doPage[R, T](ctx, c, c.listURL(rn), opts)
}

func trimSlash[S ~string](s S) string { return strings.Trim(string(s), "/") }

type resourceStringer interface {
	isResource()
	String() string
}

func (c *Client) listURL(resource resourceStringer) string {
	return fmt.Sprintf("%s/%s/", c.pokeAPIRoot, trimSlash(resource.String()))
}
func (c *Client) getURL(resource resourceStringer, ident string) string {
	return fmt.Sprintf("%s/%s/%s/", c.pokeAPIRoot, trimSlash(resource.String()), trimSlash(ident))
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
func (r APIResource[T]) Get(ctx context.Context, c *Client) (*T, error) {
	return do[*T](ctx, c, r.URL, nil)
}

// A NamedIdentifier is embedded into resources that are named.
//
// A resource directly embedding a NamedIdentifier will have a named get/list
// client function pair generated for it by gettergen.
type NamedIdentifier struct {
	//gettergen:ignore
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

// ListOpts are available on all List* endpoints, allowing you to set up your
// own pagination start point. Pagination will continue using the provided
// Limit for every page.
type ListOpts struct {
	Limit, Offset int
}

// urlValues converts the ListOpts to its corresponding url.Values. It is
// safe to call on nil ListOpts.
func (lo *ListOpts) urlValues() url.Values {
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

// GettableAPIResource is implemented by APIResource and NamedAPIResource to
// allow resources to be directly retrieved by their reference, either as part
// of a Page or when referenced by another resource.
type GettableAPIResource[T any] interface {
	Get(ctx context.Context, client *Client) (*T, error)
}

// A Page represents a list of APIResource s or NamedAPIResource s. It also
// includes information on the total number of resources in the result set, and
// how to view the Next & Previous Page s.
//
// If a page is requested that does not exist, ErrListExhausted is returned.
// PokéAPI does not distinguish between "resource not found" and "no items at
// page index", so this is a design decision taken to optimise the common case
// (as all Resources exported by this package are guaranteed to be able to be
// List-ed)
//
//	client.ListBerries(ctx, &pokeapi.ListOpts{Offset: 1000}) => ErrListExhausted
type Page[R GettableAPIResource[T], T any] struct {
	Count    int     `json:"count"`    // The total number of resources available from this API.
	Next     *string `json:"next"`     // The URL for the next page in the list.
	Previous *string `json:"previous"` // The URL for the previous page in the list.
	Results  []R     `json:"results"`
}

// GetNext retrieves the Page at Page.Next. If there is no next page,
// ErrListExhausted is returned.
func (p *Page[R, T]) GetNext(ctx context.Context, c *Client) (*Page[R, T], error) {
	if p.Next == nil {
		return nil, ErrListExhausted
	}
	return doPage[R, T](ctx, c, *p.Next, nil)
}

// GetPrevious retrieves the Page at Page.Previous. If there is no previous
// page, ErrListExhausted is returned.
func (p *Page[R, T]) GetPrevious(ctx context.Context, c *Client) (*Page[R, T], error) {
	if p.Previous == nil {
		return nil, ErrListExhausted
	}

	return doPage[R, T](ctx, c, *p.Previous, nil)
}

// The noCache is the default Cache implementation used by a Client. While it is
// valid for use, it does not perform any actual caching.
type noCache struct{}

func (noCache) Lookup(ctx context.Context, _ string, loader CacheLoader) (any, error) {
	return loader(ctx)
}

// zero creates and returns the zero value of T.
func zero[T any]() (z T) { return }

// do performs a type-safe http GET operation, using the Client's cache &
// http.Client.
func do[T any](ctx context.Context, c *Client, url string, values url.Values) (T, error) {
	res, err := c.cache.Lookup(
		ctx,
		url,
		func(ctx context.Context) (any, error) {
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
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

// doPage calls do to get the requested Page, and then performs the common error
// check - ErrNotFound for a Page get returns ErrListExhausted (set of resources
// is empty).
func doPage[R GettableAPIResource[T], T any](
	ctx context.Context,
	c *Client,
	url string,
	opts *ListOpts,
) (*Page[R, T], error) {
	p, err := do[*Page[R, T]](ctx, c, url, opts.urlValues())
	switch {
	case errors.Is(err, ErrNotFound):
		return nil, ErrListExhausted
	case err != nil:
		return nil, err
	default:
		return p, nil
	}
}

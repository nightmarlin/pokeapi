# PokéAPI

A Go wrapper around the [PokéAPI](https://pokeapi.co/) project!

Featuring a customisable http client, caching, strong types, and helpful
usability features.

## Usage

### Client

To create a client, you can simply do

```go
import "github.com/nightmarlin/pokeapi"

func main() {
	c := pokeapi.NewClient(nil)
}
```

While this is valid, the created client has no caching strategy and will use the
`http.DefaultClient` to perform requests. To provide a caching strategy simply
pass one in - for more details on PokéAPIs caching behaviour, see
[Caching](#caching). A more useful client creation would be:

```go
import "github.com/nightmarlin/pokeapi/cache"

func main() {
	c := pokeapi.NewClient(
		&pokeapi.ClientOpts{
			Cache: cache.NewLRU(nil),
		},
	)
}
```

If the API returns a non-200/404 response, a `HTTPError` will be returned by the
call, containing the returned status code. Other errors, such as ones from the
http client itself, are returned untouched. As a special case, `404 Not Found`
is represented as `ErrNotFound`, but you are still able to cast the error to a
`HTTPError` to retrieve the status code if you need.

### Resources

PokéAPI resources always have a numeric ID, and most have a name. To save you
the concern of working out which to use when fetching, just use
`resource.Ident()` when calling `Get*` methods.

`(Named)APIResource`s represent _references_ to other resources, which can be
retrieved by calling `(Named)ApiResource.Get(ctx, c)`. This returns the exact
resource, correctly typed - no casting required!

`Page`s are sections of a paginated list of `(Named)APIResource`s. The
next/previous page of results can be retrieved with `Page.Get(Next|Previous)`.

### Caching

> [The PokéAPI docs request that users of the API cache responses to reduce load](https://pokeapi.co/docs/v2#fairuse).
> Callers that don't respect this are liable to be permanently banned.

This library uses a caching strategy that means a cache implementation only
needs to do the work for any concurrent request for the same resource once.

It's relatively straightforward to implement a `pokeapi.Cache`, and users are
more than welcome to write their own implementations. If you want to do so, all
need to do is implement `pokeapi.Cache`. To verify correctness, a test suite is
provided at
`pokeapi/cache/cachetest.TestCache(*testing.T, func(int) pokeapi.Cache)`. You
should test any additional behaviours of your cache yourself, such as its
eviction strategy or TTLs.

An LRU cache that supports TTL expiry is provided out of the box, as well as a
`Wrapper` that should be suitable for use with most cache implementations.

> You're also more than welcome to set no cache and use your own implementation
> external to the pokeapi client if that better suits your needs.

A notable implementation detail: the full URL is always provided to the cache,
including the query parameters.

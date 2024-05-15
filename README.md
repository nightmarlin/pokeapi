# PokéAPI

A Go wrapper around the [PokéAPI](https://pokeapi.co/) project!

Featuring a customisable http client, caching, strong types, and helpful
usability features.

## Usage

### Client

To create a client, you can simply do

```go
c := pokeapi.NewClient(nil)
```

While this is valid, the created client has no caching strategy and will use the 
`http.DefaultClient` to perform requests. To provide a caching strategy simply 
pass one in - for more details on PokéAPIs caching behaviour, see 
[Caching](#caching).

If the API returns an error response, a `HTTPError` will be returned by the 
call, containing the returned status code. Other errors are returned untouched.

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

This library uses a caching strategy that ensures it only needs to do the work 
for any concurrent request for the same resource once. 

Implementing the transaction-like behaviour can be a bit of a difficult task,
so a simple LRU cache is provided out-of-the-box, alongside a basic wrapper for
traditional get/put caches.

If you want to write your own cache implementation, you can do so by
implementing `pokeapi.Cache`. To verify its correctness, use 
`pokeapi/cache/cachetest.TestCache(*testing.T, func(int) pokeapi.Cache)`

> You're also more than welcome to set no cache and use your own implementation 
> external to the pokeapi client.

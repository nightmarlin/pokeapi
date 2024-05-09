// Package pokeapi provides a Client and strong set of types for use with the
// PokeAPI project (https://pokeapi.co/docs/v2).
//
// The Client also supports custom caching strategies - or cache.None if you
// intend on implementing your own. Under PokeAPIs Fair Use Policy
// (https://pokeapi.co/docs/v2#fairuse), you should "locally cache resources
// whenever you request them". This package aims to make meeting that criteria
// easy for you!
//
// Resource doc comments in this module are generally taken directly from
// PokeAPIs own docs - with some edits to allow go doc to link between
// identifiers where possible.
//
// A Massive Thank You to Paul Hallett & the PokéAPI contributor team for
// maintaining the api this package wraps :D.
package pokeapi

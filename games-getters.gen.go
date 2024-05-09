// Code generated by github.com/nightmarlin/pokeapi/cmd/gettergen@v0"; DO NOT EDIT.

package pokeapi

import "context"

const GenerationResource Resource = "generation"

func (c *Client) GetGeneration(ctx context.Context, ident string) (*Generation, error) {
	return do[*Generation](ctx, c, c.getURL(GenerationResource, ident), nil)
}
func (c *Client) ListGenerations(ctx context.Context, opts *ListOptions) (*Page[NamedAPIResource[Generation], Generation], error) {
	return do[*Page[NamedAPIResource[Generation], Generation]](ctx, c, c.listURL(GenerationResource), opts.urlValues())
}

const PokedexResource Resource = "pokedex"

func (c *Client) GetPokedex(ctx context.Context, ident string) (*Pokedex, error) {
	return do[*Pokedex](ctx, c, c.getURL(PokedexResource, ident), nil)
}
func (c *Client) ListPokedexs(ctx context.Context, opts *ListOptions) (*Page[NamedAPIResource[Pokedex], Pokedex], error) {
	return do[*Page[NamedAPIResource[Pokedex], Pokedex]](ctx, c, c.listURL(PokedexResource), opts.urlValues())
}

const VersionResource Resource = "version"

func (c *Client) GetVersion(ctx context.Context, ident string) (*Version, error) {
	return do[*Version](ctx, c, c.getURL(VersionResource, ident), nil)
}
func (c *Client) ListVersions(ctx context.Context, opts *ListOptions) (*Page[NamedAPIResource[Version], Version], error) {
	return do[*Page[NamedAPIResource[Version], Version]](ctx, c, c.listURL(VersionResource), opts.urlValues())
}

const VersionGroupResource Resource = "version-group"

func (c *Client) GetVersionGroup(ctx context.Context, ident string) (*VersionGroup, error) {
	return do[*VersionGroup](ctx, c, c.getURL(VersionGroupResource, ident), nil)
}
func (c *Client) ListVersionGroups(ctx context.Context, opts *ListOptions) (*Page[NamedAPIResource[VersionGroup], VersionGroup], error) {
	return do[*Page[NamedAPIResource[VersionGroup], VersionGroup]](ctx, c, c.listURL(VersionGroupResource), opts.urlValues())
}

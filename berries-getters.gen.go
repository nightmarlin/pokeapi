// Code generated by github.com/nightmarlin/pokeapi/cmd/gettergen@v0"; DO NOT EDIT.

package pokeapi

import "context"

const BerryResource Resource = "berry"

func (c *Client) GetBerry(ctx context.Context, ident string) (*Berry, error) {
	return do[*Berry](ctx, c, c.getURL(BerryResource, ident), nil)
}
func (c *Client) ListBerries(ctx context.Context, opts *ListOptions) (*Page[NamedAPIResource[Berry], Berry], error) {
	return do[*Page[NamedAPIResource[Berry], Berry]](ctx, c, c.listURL(BerryResource), opts.urlValues())
}

const BerryFirmnessResource Resource = "berry-firmness"

func (c *Client) GetBerryFirmness(ctx context.Context, ident string) (*BerryFirmness, error) {
	return do[*BerryFirmness](ctx, c, c.getURL(BerryFirmnessResource, ident), nil)
}
func (c *Client) ListBerryFirmnesses(ctx context.Context, opts *ListOptions) (*Page[NamedAPIResource[BerryFirmness], BerryFirmness], error) {
	return do[*Page[NamedAPIResource[BerryFirmness], BerryFirmness]](ctx, c, c.listURL(BerryFirmnessResource), opts.urlValues())
}

const BerryFlavorResource Resource = "berry-flavor"

func (c *Client) GetBerryFlavor(ctx context.Context, ident string) (*BerryFlavor, error) {
	return do[*BerryFlavor](ctx, c, c.getURL(BerryFlavorResource, ident), nil)
}
func (c *Client) ListBerryFlavors(ctx context.Context, opts *ListOptions) (*Page[NamedAPIResource[BerryFlavor], BerryFlavor], error) {
	return do[*Page[NamedAPIResource[BerryFlavor], BerryFlavor]](ctx, c, c.listURL(BerryFlavorResource), opts.urlValues())
}

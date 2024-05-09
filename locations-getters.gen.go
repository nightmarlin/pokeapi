// Code generated by github.com/nightmarlin/pokeapi/cmd/gettergen@v0"; DO NOT EDIT.

package pokeapi

import "context"

const LocationResource Resource = "location"

func (c *Client) GetLocation(ctx context.Context, ident string) (*Location, error) {
	return do[*Location](ctx, c, c.getURL(LocationResource, ident), nil)
}
func (c *Client) ListLocations(ctx context.Context, opts *ListOptions) (*Page[NamedAPIResource[Location], Location], error) {
	return do[*Page[NamedAPIResource[Location], Location]](ctx, c, c.listURL(LocationResource), opts.urlValues())
}

const LocationAreaResource Resource = "location-area"

func (c *Client) GetLocationArea(ctx context.Context, ident string) (*LocationArea, error) {
	return do[*LocationArea](ctx, c, c.getURL(LocationAreaResource, ident), nil)
}
func (c *Client) ListLocationAreas(ctx context.Context, opts *ListOptions) (*Page[NamedAPIResource[LocationArea], LocationArea], error) {
	return do[*Page[NamedAPIResource[LocationArea], LocationArea]](ctx, c, c.listURL(LocationAreaResource), opts.urlValues())
}

const PalParkAreaResource Resource = "pal-park-area"

func (c *Client) GetPalParkArea(ctx context.Context, ident string) (*PalParkArea, error) {
	return do[*PalParkArea](ctx, c, c.getURL(PalParkAreaResource, ident), nil)
}
func (c *Client) ListPalParkAreas(ctx context.Context, opts *ListOptions) (*Page[NamedAPIResource[PalParkArea], PalParkArea], error) {
	return do[*Page[NamedAPIResource[PalParkArea], PalParkArea]](ctx, c, c.listURL(PalParkAreaResource), opts.urlValues())
}

const RegionResource Resource = "region"

func (c *Client) GetRegion(ctx context.Context, ident string) (*Region, error) {
	return do[*Region](ctx, c, c.getURL(RegionResource, ident), nil)
}
func (c *Client) ListRegions(ctx context.Context, opts *ListOptions) (*Page[NamedAPIResource[Region], Region], error) {
	return do[*Page[NamedAPIResource[Region], Region]](ctx, c, c.listURL(RegionResource), opts.urlValues())
}
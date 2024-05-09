// Code generated by github.com/nightmarlin/pokeapi/cmd/gettergen@v0"; DO NOT EDIT.

package pokeapi

import "context"

const ContestTypeResource Resource = "contest-type"

func (c *Client) GetContestType(ctx context.Context, ident string) (*ContestType, error) {
	return do[*ContestType](ctx, c, c.getURL(ContestTypeResource, ident), nil)
}
func (c *Client) ListContestTypes(ctx context.Context, opts *ListOptions) (*Page[NamedAPIResource[ContestType], ContestType], error) {
	return do[*Page[NamedAPIResource[ContestType], ContestType]](ctx, c, c.listURL(ContestTypeResource), opts.urlValues())
}

const ContestEffectResource Resource = "contest-effect"

// GetContestEffect only accepts the ID of the desired ContestEffect.
func (c *Client) GetContestEffect(ctx context.Context, id string) (*ContestEffect, error) {
	return do[*ContestEffect](ctx, c, c.getURL(ContestEffectResource, id), nil)
}
func (c *Client) ListContestEffects(ctx context.Context, opts *ListOptions) (*Page[APIResource[ContestEffect], ContestEffect], error) {
	return do[*Page[APIResource[ContestEffect], ContestEffect]](ctx, c, c.listURL(ContestEffectResource), opts.urlValues())
}

const SuperContestEffectResource Resource = "super-contest-effect"

// GetSuperContestEffect only accepts the ID of the desired SuperContestEffect.
func (c *Client) GetSuperContestEffect(ctx context.Context, id string) (*SuperContestEffect, error) {
	return do[*SuperContestEffect](ctx, c, c.getURL(SuperContestEffectResource, id), nil)
}
func (c *Client) ListSuperContestEffects(ctx context.Context, opts *ListOptions) (*Page[APIResource[SuperContestEffect], SuperContestEffect], error) {
	return do[*Page[APIResource[SuperContestEffect], SuperContestEffect]](ctx, c, c.listURL(SuperContestEffectResource), opts.urlValues())
}

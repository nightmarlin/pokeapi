package cache

// The None cache does what it says on the tin: no caching is done. Use this
// cache if you intend to use your own dedicated caching strategy.
type None struct{}

func (None) Cache(string, any)         {}
func (None) Lookup(string) (any, bool) { return nil, false }

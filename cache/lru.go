package cache

import (
	"context"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"

	"github.com/nightmarlin/pokeapi"
)

const (
	defaultLRUCacheSize        = 500
	defaultLRUCacheTTL         = time.Duration(0)
	defaultLRUCacheExpiryDelay = 24 * 7 * time.Hour
)

// LRU implements a Least-Recently-Used pokeapi.Cache. An LRU cache is a
// queue-like structure, with some extra semantics on Lookup & a map to allow
// for O(1) lookups:
//
//   - On Lookup, if a cache hit occurs, the entry is moved to the top of the list
//     (it is now the youngest entry).
//   - On Put, the entry is moved to/inserted at the top of the list. If an insert
//     causes the length of the list to exceed the cache capacity, the oldest
//     entry is dropped. A Put occurs on the first call to
//     [pokeapi.CacheLookup.Hydrate] on a [pokeapi.CacheLookup] returned by
//     [LRU.Lookup] (as long as [pokeapi.CacheLookup.Close] has not yet been called).
//
// If multiple cache lookups are opened for the same url, the LRU cache will
// ensure that they are not executed in parallel - instead ensuring these
// lookups occur one-after-another. Parallel cache lookups for multiple urls
// will be serviced as normal.
//
// A TTL may be set, in which case the LRU cache will expire entries in
// accordance with it. As key expiry briefly locks the entire cache, it is only
// run once within a specified period. See LRUOpts.ExpiryDelay for details.
type LRU struct {
	mux sync.Mutex

	ttl         time.Duration    // how long cache entries should be stored before eviction
	clock       func() time.Time // get the current time
	expiryDelay time.Duration    // how often to wait between expiry runs
	lastExpiry  time.Time        // when the last expiry was run

	capacity int                       // the maximum size of the cache
	length   int                       // the current number of values in the cache
	youngest *lruCacheEntry            // the most recently accessed/inserted cache item
	oldest   *lruCacheEntry            // the next cache item to evict
	entries  map[string]*lruCacheEntry // lookup map for O(1) lookups

	ongoing singleflight.Group
}

type LRUOpts struct {
	// The maximum capacity of the cache. Default 500.
	Size int

	// Provide a custom time function - useful for testing. Default time.Now().
	Clock func() time.Time

	// How long cached entries should be stored for. Default 0 (forever). Key
	// expiry briefly locks the cache, so ExpiryDelay can be used to limit how
	// often TTL is checked.
	TTL time.Duration

	// How long to wait between expiry runs. Default ~1 week. If set to 0, will
	// check for expired keys every Lookup. Only applies if TTL != 0.
	ExpiryDelay *time.Duration
}

// NewLRU constructs a new LRU cache for use in accordance with the provided
// LRUOpts.
func NewLRU(opts *LRUOpts) *LRU {
	cacheSize := defaultLRUCacheSize

	lru := LRU{
		ttl:         defaultLRUCacheTTL,
		expiryDelay: defaultLRUCacheExpiryDelay,
		clock:       func() time.Time { return time.Now().UTC() },
	}

	if opts != nil {
		if opts.Size > 0 {
			cacheSize = opts.Size
		}
		if opts.TTL > 0 {
			lru.ttl = opts.TTL
		}
		if opts.Clock != nil {
			lru.clock = opts.Clock
		}
		if opts.ExpiryDelay != nil && *opts.ExpiryDelay >= 0 {
			lru.expiryDelay = *opts.ExpiryDelay
		}
	}

	// only allocate the map once we know how big it needs to be
	lru.capacity = cacheSize
	lru.entries = make(map[string]*lruCacheEntry, cacheSize)
	return &lru
}

type lruCacheEntry struct {
	url   string
	value any

	expireAt time.Time

	older   *lruCacheEntry
	younger *lruCacheEntry
}

// extractEntry removes the requested lruCacheEntry from the LRU cache and
// updates any relationships to reflect this. e will also have its relationships
// removed.
func (lru *LRU) extractEntry(e *lruCacheEntry) {
	// update linked list tails
	if lru.youngest == e {
		lru.youngest = e.older
	}
	if lru.oldest == e {
		lru.oldest = e.younger
	}

	// update neighbours
	if e.older != nil {
		e.older.younger = e.younger
		e.older = nil
	}
	if e.younger != nil {
		e.younger.older = e.older
		e.younger = nil
	}

	delete(lru.entries, e.url)
	lru.length -= 1
}

// insertValue creates a new lruCacheEntry and inserts it at the top of the
// cache. if insertion caused the cache length to exceed its capacity,
// extraneous elements are dropped from the cache.
func (lru *LRU) insertValue(url string, value any) {
	e := &lruCacheEntry{
		url:      url,
		value:    value,
		older:    lru.youngest,
		expireAt: lru.clock().Add(lru.ttl),
	}

	if lru.youngest != nil {
		lru.youngest.younger = e
	}
	lru.youngest = e

	if lru.oldest == nil {
		lru.oldest = e
	}

	lru.entries[url] = e
	lru.length += 1

	// delete the oldest entry if capacity exceeded
	if lru.length > lru.capacity {
		o := lru.oldest
		if o == nil {
			return
		}

		if o.younger != nil {
			o.younger.older = nil
		}
		lru.oldest = o.younger

		delete(lru.entries, o.url)

		o.value = nil
		o.older = nil
		o.younger = nil
		lru.length -= 1
	}
}

func (lru *LRU) Lookup(
	ctx context.Context,
	url string,
	loadOnMiss pokeapi.CacheLoader,
) (any, error) {
	res, err, _ := lru.ongoing.Do(
		"url",
		func() (any, error) {
			lru.mux.Lock()

			if e := lru.entries[url]; e != nil {
				// bump entry to top of list
				lru.extractEntry(e)
				lru.insertValue(url, e.value)

				lru.mux.Unlock()
				return e.value, nil
			}

			lru.mux.Unlock()

			res, err := loadOnMiss(ctx)
			if err != nil {
				return nil, err
			}

			lru.mux.Lock()

			if e, ok := lru.entries[url]; ok {
				lru.extractEntry(e)
			}

			lru.insertValue(url, res)

			lru.mux.Unlock()
			return res, nil
		},
	)

	go lru.expire() // run expiry in the background

	return res, err
}

// expire scans through the LRU cache and deletes entries that were written
// before now-ttl, as long as lru.ttl != 0.
func (lru *LRU) expire() {
	if lru.ttl == 0 {
		return
	}

	defer lru.mux.Unlock()
	lru.mux.Lock()

	now := lru.clock()

	// don't run expiry if expired recently
	if lru.lastExpiry.Add(lru.expiryDelay).After(now) {
		return
	}
	lru.lastExpiry = now

	for _, e := range lru.entries {
		if e.expireAt.Before(now) {
			// fun fact: it's safe to delete entries from a map as you iterate through
			// that map! See https://go.dev/ref/spec#For_statements for more details.
			lru.extractEntry(e)
			e.value = nil
		}
	}
}

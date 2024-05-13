package cache

import (
	"context"
	"sync"

	"github.com/nightmarlin/pokeapi"
)

const defaultLRUCacheSize = 50

// LRU implements a Least-Recently-Used pokeapi.Cache. An LRU cache is a
// queue-like structure, with some extra semantics on Lookup & a map to allow
// for O(1) lookups:
//
// - On Lookup, if a cache hit occurs, the entry is moved to the top of the list
// (it is now the youngest entry).
//
// - On Put, the entry is moved to/inserted at the top of the list. If an insert
// causes the length of the list to exceed the cache capacity, the oldest
// entry is dropped. A Put occurs on the first call to
// [pokeapi.CacheLookup.Hydrate] on a [pokeapi.CacheLookup] returned by
// [LRU.Lookup] (as long as [pokeapi.CacheLookup.Close] has not yet been called).
//
// If multiple cache lookups are opened for the same url, the LRU cache will
// ensure that they are not executed in parallel - instead ensuring these
// lookups occur one-after-another. Parallel cache lookups for multiple urls
// will be serviced as normal.
type LRU struct {
	mux sync.Mutex

	capacity int                       // the maximum size of the cache
	length   int                       // the current number of values in the cache
	youngest *lruCacheEntry            // the most recently accessed/inserted cache item
	oldest   *lruCacheEntry            // the next cache item to evict
	entries  map[string]*lruCacheEntry // lookup map for O(1) lookups

	ongoing sync.Map // ongoing implements a key-level lock - no lookups for the same key may occur in parallel
}

// NewLRU constructs a new LRU cache for use. It can hold at most `size`
// entries, and once that limit is exceeded the least-recently-used (read or
// written) cache entry will be evicted. If `size` <= 0, the default cache size
// of 50 will be used.
func NewLRU(size int) *LRU {
	if size <= 0 {
		size = defaultLRUCacheSize
	}

	return &LRU{capacity: size, entries: make(map[string]*lruCacheEntry, size)}
}

type lruCacheEntry struct {
	url   string
	value any

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
	e := &lruCacheEntry{url: url, value: value, older: lru.youngest}

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

func (lru *LRU) Lookup(_ context.Context, url string) pokeapi.CacheLookup {
	// 1. acquire lock to process url
	//   - ensures two concurrent lookups for the same url are ordered, so work is never duplicated
	// 2. acquire lock to modify cache
	// 3. if url already in cache, mark it as most-recently-used
	// 4. return lookup with ability to insert resource & release the url-level lock

	for {
		// spin: attempt to acquire lock on individual url.
		// todo: reduce spins?

		_, urlIsLocked := lru.ongoing.LoadOrStore(url, struct{}{})
		if !urlIsLocked {
			break
		}
	}

	lru.mux.Lock()
	defer lru.mux.Unlock()

	var (
		value    any
		hasValue bool
	)

	if e := lru.entries[url]; e != nil {
		value = e.value
		hasValue = true

		// bump entry to top of list
		lru.extractEntry(e)
		lru.insertValue(url, value)

		// clean up references
		e.value = nil
	}

	return &lruLookup{
		hasValue:  hasValue,
		value:     value,
		putFn:     func(resource any) { lru.Put(url, resource) },
		cleanupFn: func() { lru.ongoing.Delete(url) },
	}
}

// Put inserts an item into the cache & unlocks lookups waiting on that
// resource.
func (lru *LRU) Put(url string, resource any) {
	defer lru.mux.Unlock()
	lru.mux.Lock()

	// 1. if entry already in list, remove it
	// 2. insert entry at top of list
	// 3. if list length now exceeds capacity, eject the oldest item(s)

	if e, ok := lru.entries[url]; ok {
		lru.extractEntry(e)
		e.value = nil
	}

	lru.insertValue(url, resource)
}

type lruLookup struct {
	mux  sync.RWMutex
	once sync.Once

	hasValue bool
	value    any

	// insert (new) resource value on Hydrate.
	putFn func(resource any)

	// release url on Close or Hydrate. must always be called, and when called must be after putFn
	cleanupFn func()
}

func (l *lruLookup) Value(context.Context) (_ any, ok bool) {
	defer l.mux.RUnlock()
	l.mux.RLock()
	return l.value, l.hasValue
}

// cleanup closes the lruLookup and removes any resources it is using.
func (l *lruLookup) cleanup() {
	l.cleanupFn()
	l.hasValue = false
	l.value = nil
	l.putFn = nil
	l.cleanupFn = nil
}

func (l *lruLookup) Hydrate(_ context.Context, resource any) {
	defer l.mux.Unlock()
	l.mux.Lock()

	l.once.Do(
		func() {
			l.putFn(resource)
			l.cleanup()
		},
	)

}

func (l *lruLookup) Close(context.Context) {
	defer l.mux.Unlock()
	l.mux.Lock()

	l.once.Do(l.cleanup)
}

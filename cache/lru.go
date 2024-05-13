package cache

import (
	"sync"

	"github.com/nightmarlin/pokeapi"
)

const defaultLRUCacheSize = 50

// LRU implements a Least-Recently-Used pokeapi.Cache. An LRU cache is a
// waterfall-like structure:
//
// - On Lookup, if a cache hit occurs, the entry is moved to the top of the list.
//
// - On Put, the entry is moved to/inserted at the top of the list. If an insert
// causes the length of the list to exceed the cache capacity, the oldest
// entry is dropped. A Put occurs on the first call to
// [pokeapi.CacheLookup.Hydrate] on a [pokeapi.CacheLookup] returned by
// [LRU.Lookup] (as long as [pokeapi.CacheLookup.Close] has not yet been called).
type LRU struct {
	mux sync.Mutex

	capacity int
	length   int

	youngest *lruCacheEntry
	oldest   *lruCacheEntry

	entries map[string]*lruCacheEntry
}

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
	for lru.length > lru.capacity {
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

func (lru *LRU) Lookup(url string) pokeapi.CacheLookup {
	// 1. if url in open lookups list, wait until it isn't. this ensures two concurrent lookups for the same resource are ordered (requires sync)
	// 2. if url in list, move it to top of the list
	// 4. add url to open lookups list
	// 3. return lookup with reference to lru in putFn

	// todo: if url in open lookups list, wait until it isn't... (requires sync)

	defer lru.mux.Unlock()
	lru.mux.Lock()

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

	// todo: add url to open lookups list

	return &lruLookup{
		hasValue: hasValue,
		value:    value,
		putFn:    func(resource any) { lru.Put(url, resource) },
	}
}

// Put inserts an item into the cache & unlocks lookups waiting on that
// resource.
func (lru *LRU) Put(url string, resource any) {
	defer lru.mux.Unlock()
	lru.mux.Lock()

	// 1. remove entry from ongoing lookups list
	// 2. if entry already in list, remove it
	// 3. insert entry at top of list
	// 4. if list length now exceeds capacity, eject the oldest item(s)

	// todo: remove url from ongoing lookups list

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

	putFn func(resource any)
}

func (l *lruLookup) Value() (_ any, ok bool) {
	defer l.mux.RUnlock()
	l.mux.RLock()
	return l.value, l.hasValue
}

// cleanup closes the lruLookup and removes any resources it is using.
func (l *lruLookup) cleanup() {
	l.hasValue = false
	l.value = nil
	l.putFn = nil
}

func (l *lruLookup) Hydrate(resource any) {
	defer l.mux.Unlock()
	l.mux.Lock()

	l.once.Do(
		func() {
			l.putFn(resource)
			l.cleanup()
		},
	)

}

func (l *lruLookup) Close() {
	defer l.mux.Unlock()
	l.mux.Lock()

	l.once.Do(l.cleanup)
}

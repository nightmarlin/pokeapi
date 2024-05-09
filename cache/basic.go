package cache

import (
	"sync"
)

const defaultBasicCacheSize = 50

type basicCacheEntry struct {
	url string
	val any
}

// The Basic cache provides a simple wraparound cache for the last N requests.
// Once N responses are cached, new responses will overwrite the oldest ones.
type Basic struct {
	mux sync.RWMutex

	store          []basicCacheEntry
	nextWriteIndex int
}

// NewBasic constructs a new Basic cache with the specified size. If size <= 0,
// the default size of 50 records is used.
func NewBasic(size int) *Basic {
	if size <= 0 {
		size = defaultBasicCacheSize
	}
	return &Basic{store: make([]basicCacheEntry, size)}
}

func (b *Basic) Cache(url string, val any) {
	defer b.mux.Unlock()
	b.mux.Lock()

	if len(b.store) == 0 {
		return
	}

	for idx := range b.store {
		if b.store[idx].url == url {
			b.store[idx].val = val
			return
		}
	}

	b.store[b.nextWriteIndex] = basicCacheEntry{url: url, val: val}
	b.nextWriteIndex = (b.nextWriteIndex + 1) % len(b.store)
}

func (b *Basic) Lookup(url string) (val any, ok bool) {
	defer b.mux.RUnlock()
	b.mux.RLock()
	for _, entry := range b.store {
		if entry.url == url {
			return entry.val, true
		}
	}
	return nil, false
}

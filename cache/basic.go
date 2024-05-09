package cache

import (
	"sync"

	"github.com/nightmarlin/pokeapi"
)

const defaultBasicCacheSize = 50

type basicCacheEntry struct {
	url string
	val any
}

type basicLookup struct {
	mux    sync.RWMutex
	val    any
	hasVal bool

	closeOnce sync.Once
	c         chan<- any
}

func (b *basicLookup) Value() (any, bool) {
	defer b.mux.RUnlock()
	b.mux.RLock()

	return b.val, b.hasVal
}

func (b *basicLookup) cleanup() {
	defer b.mux.Unlock()
	b.mux.Lock()

	close(b.c)
	b.val = nil
	b.hasVal = false
}

func (b *basicLookup) Hydrate(resource any) {
	b.closeOnce.Do(
		func() {
			b.c <- resource
			b.cleanup()
		},
	)
}

func (b *basicLookup) Close() { b.closeOnce.Do(b.cleanup) }

// The Basic cache provides a simple wraparound cache for the last N requests.
// Once N responses are cached, new responses will overwrite the oldest ones.
type Basic struct {
	mux sync.Mutex

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

func (b *Basic) Lookup(url string) pokeapi.CacheLookup {
	// acquire cache-level lock
	b.mux.Lock()

	var (
		storedValue      any
		storedValueIndex = -1
	)

	for idx := range b.store {
		if b.store[idx].url == url {
			storedValue, storedValueIndex = b.store[idx].val, idx
			break
		}
	}

	c := make(chan any)

	bl := &basicLookup{val: storedValue, hasVal: storedValueIndex != -1, c: c}

	go func() {
		// unlocks cache once basicLookup closes chan / passes a value on it
		defer b.mux.Unlock()

		select {
		case v, ok := <-c:
			if !ok {
				return
			}

			// refresh cache

			if storedValueIndex != -1 {
				// overwrite current value
				b.store[storedValueIndex].val = v
				return
			}

			// drop the oldest entry & write the new one
			if len(b.store) == 0 {
				return
			}

			b.store[b.nextWriteIndex] = basicCacheEntry{url: url, val: v}
			b.nextWriteIndex = (b.nextWriteIndex + 1) % len(b.store)
		}
	}()

	return bl
}

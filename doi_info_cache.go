package main

import (
	"sync"
	"time"

	lru "github.com/hashicorp/golang-lru"
)

const (
	CacheDefaultSize = 100
	CacheDefaultAge  = 1 * time.Minute
)

// DoiCacheConfig configures a doi cache
type DoiCacheConfig struct {
	MaxAge  time.Duration // Maximum age before evicting a cache entry
	MaxSize int           // Maximum number of entries
}

// DoiCache caches information for a limited number of DOIs, for a specified amount of time.
type DoiCache struct {
	m      sync.Mutex
	config DoiCacheConfig
	cache  *lru.Cache
}

type cacheEntry struct {
	sync.RWMutex
	info *DoiInfo
	ok   bool
	err  error
}

// NewDoiCache initializes a new doi cache
func NewDoiCache(cfg DoiCacheConfig) *DoiCache {
	if cfg.MaxSize <= 0 {
		cfg.MaxSize = CacheDefaultSize
	}

	if cfg.MaxAge <= 0 {
		cfg.MaxAge = CacheDefaultAge
	}

	cache, _ := lru.New(cfg.MaxSize)

	return &DoiCache{
		cache:  cache,
		config: cfg,
	}
}

// GetOrAdd adds an entry to the cache via invoking the given generator
// function, if there isn't one already.   If there is already  a cache
// entry, it just gets the old cache value, and the doi fetch function is never
// invoked.
//
// The doi fetch function provides the doi info to cache, possibly performing
// a fetch that blocks for a while.  Future calls to GetOrAdd for the same doi
// will block until a value is available or the function returns an error.
// In the case of an error, the value will not be added to the cache,
// and all pending Get requests will return the error
func (c *DoiCache) GetOrAdd(doi string, fetchDoi func() (*DoiInfo, error)) (*DoiInfo, error) {

	// Critical section.  Check that we don't have a cached entry, and create/add a locked one if not
	cached, entry, found, err := func() (*DoiInfo, *cacheEntry, bool, error) {
		c.m.Lock()
		defer c.m.Unlock()

		cached, ok, err := c.get(doi)
		if ok {
			return cached, nil, ok, err
		}
		entry := &cacheEntry{}
		entry.Lock()
		c.cache.Add(doi, entry)

		return nil, entry, false, nil
	}()

	if found {
		return cached, err
	}

	// OK, now execute the doi getch function and unlock the cache entry when done.
	defer entry.Unlock()

	if entry.info, err = fetchDoi(); err != nil {
		entry.ok = false
		c.cache.Remove(doi)
		return nil, err
	}

	entry.ok = true
	time.AfterFunc(c.config.MaxAge, func() {
		c.cache.Remove(doi)
	})

	return entry.info, nil
}

func (c *DoiCache) get(doi string) (info *DoiInfo, ok bool, err error) {
	v, ok := c.cache.Get(doi)

	// Nothing in cache
	if !ok {
		return nil, false, nil
	}

	e := v.(*cacheEntry)
	e.RLock()
	defer e.RUnlock()
	return e.info, e.ok, e.err
}

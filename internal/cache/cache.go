package cache

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

// Eviction interface
type bailiff interface {
	// Returns a list of items in order of eviction
	itemList(immunity string) []evictionItem
	// Returns a list of items that will be evicted
	evictionList(list []evictionItem, space int64) []string
	// Evicts the specified items
	eviction(evicted []string)
}

// Cache interface
type Memorizer interface {
	// Add an item to the cache or update a value in the cache by key
	Set(key string, value any)
	// Get the value of an item from the cache by key
	Get(key string) (any, bool)
	// Returns the number of items in the cache
	Len() int64
	// Free the specified amount of memory in the cache
	Eviction(immunity string, space int64)
	// Memorizer implement the eviction interface (bailiff)
	bailiff
}

type Cache struct {
	Limit int64
	size  atomic.Int64
	len   atomic.Int64
	data  map[string]*CacheItem
	mutex *sync.RWMutex
}

type CacheItem struct {
	UsedCount atomic.Int64
	Size      int64
	Value     any
}

func (ci *CacheItem) Use() {
	ci.UsedCount.Add(1)
}
func (ci *CacheItem) Used() int64 {
	return ci.UsedCount.Load()
}

func NewCacheItem(value any, size int64) *CacheItem {
	return &CacheItem{
		UsedCount: atomic.Int64{},
		Size:      size,
		Value:     value,
	}
}

func New(limit int64) *Cache {
	return &Cache{
		Limit: limit,
		size:  atomic.Int64{},
		len:   atomic.Int64{},
		mutex: &sync.RWMutex{},
		data:  map[string]*CacheItem{},
	}
}
func (c *Cache) Len() int64 {
	return c.len.Load()
}

func (c *Cache) Set(key string, value any) {
	var oldValueSize int64
	newValueSize := int64(unsafe.Sizeof(value))
	oldSize := c.size.Load()
	c.mutex.RLock()
	oldItem, ok := c.data[key]
	c.mutex.RUnlock()
	if ok {
		oldValueSize = oldItem.Size
	} else {
		c.len.Add(1)
	}
	newSize := oldSize - oldValueSize + newValueSize
	newItem := NewCacheItem(value, newValueSize)
	space := newSize - oldSize
	if newSize > c.Limit {
		c.Eviction(key, space)
	}
	c.mutex.Lock()
	c.data[key] = newItem
	c.mutex.Unlock()
	c.size.Store(newSize)
}

func (c *Cache) Get(key string) (any, bool) {
	c.mutex.RLock()
	item, ok := c.data[key]
	c.mutex.RUnlock()
	if ok {
		item.Use()
	}
	return item, ok
}

// {key, usedCounter, size}
type evictionItem struct {
	name string
	used int64
	size int64
}

func (c *Cache) Eviction(immunity string, space int64) {
	list := c.itemList(immunity)
	evicted := c.evictionList(list, space)
	c.eviction(evicted)
}

func (c *Cache) itemList(immunity string) []evictionItem {
	stack := []evictionItem{}
	c.mutex.RLock()
	for key, value := range c.data {
		if key != immunity {
			stack = append(stack,
				evictionItem{
					name: key,
					used: value.Used(),
					size: value.Size,
				},
			)
			prev := len(stack) - 2
			current := len(stack) - 1
			if len(stack) > 1 && stack[prev].used < stack[current].used {
				stack[prev], stack[current] = stack[current], stack[prev]
			}
		}
	}
	c.mutex.Unlock()
	return stack
}

func (c *Cache) evictionList(list []evictionItem, space int64) []string {
	var evictionSize int64
	eviction := []string{}
	for i := len(list) - 1; i >= 0; i-- {
		if evictionSize < space {
			eviction = append(eviction, list[i].name)
			evictionSize += list[i].size
		}
	}
	return eviction
}

func (c *Cache) eviction(evicted []string) {
	for _, key := range evicted {
		c.mutex.Lock()
		delete(c.data, key)
		c.mutex.Unlock()
		c.len.Add(-1)
	}
}

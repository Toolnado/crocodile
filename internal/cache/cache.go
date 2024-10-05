package cache

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

type Cache struct {
	Limit int64
	Size  atomic.Int64
	Len   atomic.Int64
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
		Size:  atomic.Int64{},
		Len:   atomic.Int64{},
		mutex: &sync.RWMutex{},
		data:  map[string]*CacheItem{},
	}
}

func (c *Cache) Set(key string, value any) {
	var oldValueSize int64
	newValueSize := int64(unsafe.Sizeof(value))
	oldSize := c.Size.Load()
	c.mutex.RLock()
	oldItem, ok := c.data[key]
	c.mutex.RUnlock()
	if ok {
		oldValueSize = oldItem.Size
	} else {
		c.Len.Add(1)
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
	c.Size.Store(newSize)
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

func (c *Cache) Eviction(usedItemKey string, space int64) {
	list := c.itemList(usedItemKey)
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
	c.mutex.Lock()
	for _, key := range evicted {
		delete(c.data, key)
	}
	c.mutex.Unlock()
}

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
		c.Eviction(oldItem, newItem, space)
	} else {
		c.mutex.Lock()
		c.data[key] = newItem
		c.mutex.Unlock()
	}
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

func (c *Cache) Eviction(oldItem, newItem *CacheItem, space int64) {}

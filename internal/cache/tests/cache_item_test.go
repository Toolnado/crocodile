package tests

import (
	"crocodile/internal/cache"
	"testing"
)

func TestNewCacheItem(t *testing.T) {
	item := cache.NewCacheItem("hello", 5)
	if item.Size != 5 {
		t.Errorf("invalid size value, expected: %d, received: %d\n", 5, item.Size)
	}
	if item.Value.(string) != "hello" {
		t.Errorf("invalid value, expected: %s, received: %s", "hello\n", item.Value)
	}
}

func TestCacheItemUse(t *testing.T) {
	item := cache.NewCacheItem("hello", 5)
	item.Use()
	if item.UsedCount.Load() != 1 {
		t.Errorf("invalid userCount value, expected: %d, received: %d\n", 1, item.UsedCount.Load())
	}
}

func TestCacheItemUsed(t *testing.T) {
	item := cache.NewCacheItem("hello", 5)
	item.Use()
	if item.Used() != 1 {
		t.Errorf("invalid userCount value, expected: %d, received: %d\n", 1, item.Used())
	}
}

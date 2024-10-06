package benchmark

import (
	"crocodile/internal/cache"
	"fmt"
	"sync"
	"testing"
)

type testData struct {
	key   string
	value []byte
}

func genData(count int) []testData {
	data := []testData{}
	for i := 0; i < count; i++ {
		data = append(data, testData{
			key:   fmt.Sprintf("key_%d", i),
			value: []byte(fmt.Sprintf("value_%d", i)),
		})
	}
	return data
}

// benchmark: simple cache set
func BenchmarkCacheSet(t *testing.B) {
	data := genData(1000)
	instance := cache.NewCache(1024 * 1024)
	t.StartTimer()
	for i := 0; i < t.N; i++ {
		for _, item := range data {
			instance.Set(item.key, item.value)
		}
	}
}

// benchmark: simple cache set if need use eviction
func BenchmarkCacheSetWithEviction(t *testing.B) {
	data := genData(1000)
	instance := cache.NewCache(10)
	t.StartTimer()
	for i := 0; i < t.N; i++ {
		for _, item := range data {
			instance.Set(item.key, item.value)
		}
	}
}

// benchmark: simple cache get
func BenchmarkCacheGet(t *testing.B) {
	data := genData(1000)
	instance := cache.NewCache(1024 * 1024)
	for _, item := range data {
		instance.Set(item.key, item.value)
	}

	t.StartTimer()
	for i := 0; i < t.N; i++ {
		for _, item := range data {
			if _, ok := instance.Get(item.key); !ok {
				t.Errorf("value for key=%s not found", item.key)
			}
		}
	}
}

// benchmark: concurrency cache set if need use eviction
func BenchmarkCacheConcurrencySet(t *testing.B) {
	data := genData(1000)
	instance := cache.NewCache(1024 * 1024)
	t.StartTimer()
	wg := sync.WaitGroup{}
	for i := 0; i < t.N; i++ {
		wg.Add(1)
		go func() {
			for _, item := range data {
				instance.Set(item.key, item.value)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

// benchmark: concurrency cache get
func BenchmarkCacheConcurrencyGet(t *testing.B) {
	data := genData(1000)
	instance := cache.NewCache(1024 * 1024)
	for _, item := range data {
		instance.Set(item.key, item.value)
	}
	t.StartTimer()
	wg := sync.WaitGroup{}
	for i := 0; i < t.N; i++ {
		wg.Add(1)
		go func() {
			for _, item := range data {
				if _, ok := instance.Get(item.key); !ok {
					t.Errorf("value for key=%s not found", item.key)
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

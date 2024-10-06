package tests

import (
	"crocodile/internal/cache"
	"testing"
)

const TestSizeBig int64 = 1 << 20

// interface test
var _ cache.Memorizer = cache.New(0)

func TestNewCache(t *testing.T) {
	instance := cache.New(TestSizeBig)
	t.Logf("\ncache: %-v\n", instance)
	if instance.Limit != TestSizeBig {
		t.Error("invalid cache size:", instance.Limit)
	}
}

func TestCacheLen(t *testing.T) {
	instance := cache.New(TestSizeBig)
	instance.Set("key", []byte("value"))
	length := instance.Len()
	var exp int64 = 1
	if length != exp {
		t.Errorf("invalid value, expected: %d, received: %d", exp, length)
	}
}

func TestCacheSet(t *testing.T) {
	tests := []struct {
		data struct {
			key   string
			value []byte
		}
		name     string
		expected any
	}{
		{struct {
			key   string
			value []byte
		}{"key1", []byte("value1")}, "first", "value1"},
		{struct {
			key   string
			value []byte
		}{"key2", []byte("value2")}, "second", "value2"},
		{struct {
			key   string
			value []byte
		}{"key3", []byte("value3")}, "third", "value3"},
	}
	instance := cache.New(12)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			instance.Set(test.data.key, test.data.value)
			res, _ := instance.Get(test.data.key)
			if string(res) != test.expected {
				t.Errorf("invalid value, expected: %d, received: %d", test.expected, res)
			}
		})
	}
}

func TestCacheGet(t *testing.T) {
	tests := []struct {
		data struct {
			key   string
			value []byte
		}
		needFound bool
		name      string
		expected  any
	}{
		{struct {
			key   string
			value []byte
		}{"key1", []byte("value1")}, true, "first", "value1"},
		{struct {
			key   string
			value []byte
		}{"key2", []byte("value2")}, false, "second", "value2"},
		{struct {
			key   string
			value []byte
		}{"key3", []byte("value3")}, true, "third", "value3"},
		{struct {
			key   string
			value []byte
		}{"key2", []byte("value2")}, true, "four", string([]byte{})},
	}
	instance := cache.New(12)
	instance.Set(tests[0].data.key, tests[0].data.value)
	instance.Set(tests[1].data.key, tests[1].data.value)
	instance.Set(tests[2].data.key, tests[2].data.value)
	t.Log("cache:", instance)
	for i, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, _ := instance.Get(test.data.key)
			if test.needFound {
				if string(res) != test.expected {
					t.Errorf("case %d error: invalid value, expected: %s, received: %s", i, test.expected, res)
				}
			}
		})
	}
}

func TestCacheEviction(t *testing.T) {
	instance := cache.New(6)
	instance.Set("key1", []byte("value1"))
	instance.Set("key2", []byte("value2"))
	res, ok := instance.Get("key1")
	if ok || res != nil {
		t.Errorf("invalid value, expected: %v, received: %d", nil, res)
	}
}

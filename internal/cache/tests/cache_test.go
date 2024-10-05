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

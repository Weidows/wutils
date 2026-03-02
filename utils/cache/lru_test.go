package cache

import (
	"testing"
)

func TestLRU_BasicGetSet(t *testing.T) {
	lru := NewLRU(100)

	lru.Set("key1", []byte("value1"))

	val, ok := lru.Get("key1")
	if !ok {
		t.Fatal("expected to find key1")
	}
	if string(val) != "value1" {
		t.Fatalf("expected value1, got %s", string(val))
	}
}

func TestLRU_Miss(t *testing.T) {
	lru := NewLRU(100)

	val, ok := lru.Get("nonexistent")
	if ok {
		t.Fatal("expected miss")
	}
	if val != nil {
		t.Fatalf("expected nil, got %s", string(val))
	}
}

func TestLRU_Eviction(t *testing.T) {
	lru := NewLRU(10)

	lru.Set("key1", []byte("1234567890"))
	lru.Set("key2", []byte("1234567890"))

	_, ok := lru.Get("key1")
	if ok {
		t.Fatal("key1 should have been evicted")
	}

	_, ok = lru.Get("key2")
	if !ok {
		t.Fatal("key2 should still exist")
	}
}

func TestLRU_Update(t *testing.T) {
	lru := NewLRU(100)

	lru.Set("key1", []byte("value1"))
	lru.Set("key1", []byte("value2"))

	val, ok := lru.Get("key1")
	if !ok {
		t.Fatal("expected to find key1")
	}
	if string(val) != "value2" {
		t.Fatalf("expected value2, got %s", string(val))
	}
}

func TestLRU_Stats(t *testing.T) {
	lru := NewLRU(100)

	lru.Set("key1", []byte("value1"))
	lru.Get("key1")
	lru.Get("nonexistent")

	hits, misses, _ := lru.Stats()
	if hits != 1 {
		t.Fatalf("expected 1 hit, got %d", hits)
	}
	if misses != 1 {
		t.Fatalf("expected 1 miss, got %d", misses)
	}
}

func TestLRU_Invalidate(t *testing.T) {
	lru := NewLRU(100)

	lru.Set("key1", []byte("value1"))
	lru.Invalidate("key1")

	_, ok := lru.Get("key1")
	if ok {
		t.Fatal("key1 should have been invalidated")
	}
}

func TestLRU_Clear(t *testing.T) {
	lru := NewLRU(100)

	lru.Set("key1", []byte("value1"))
	lru.Set("key2", []byte("value2"))
	lru.Clear()

	_, ok := lru.Get("key1")
	if ok {
		t.Fatal("key1 should have been cleared")
	}
	_, ok = lru.Get("key2")
	if ok {
		t.Fatal("key2 should have been cleared")
	}
}

func BenchmarkLRU_Set(b *testing.B) {
	lru := NewLRU(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lru.Set("key", []byte("value"))
	}
}

func BenchmarkLRU_Get(b *testing.B) {
	lru := NewLRU(10000)
	lru.Set("key", []byte("value"))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lru.Get("key")
	}
}

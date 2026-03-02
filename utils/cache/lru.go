package cache

import (
	"container/list"
	"sync"
)

type Entry struct {
	Key      string
	Data     []byte
	Size     int64
	Accesses int
}

type LRU struct {
	mu          sync.Mutex
	lru         *list.List
	cache       map[string]*list.Element
	MemoryLimit int64
	CurrentMem  int64
	Hits        int64
	Misses      int64
}

func NewLRU(limit int64) *LRU {
	return &LRU{
		lru:         list.New(),
		cache:       make(map[string]*list.Element),
		MemoryLimit: limit,
	}
}

func (c *LRU) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	elem, exists := c.cache[key]
	if !exists {
		c.Misses++
		return nil, false
	}

	c.Hits++
	c.lru.MoveToFront(elem)
	entry := elem.Value.(*Entry)
	entry.Accesses++
	return entry.Data, true
}

func (c *LRU) Set(key string, data []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	size := int64(len(data))
	if size > c.MemoryLimit {
		return
	}

	if elem, exists := c.cache[key]; exists {
		entry := elem.Value.(*Entry)
		c.CurrentMem -= entry.Size
		entry.Data = data
		entry.Size = size
		c.CurrentMem += size
		c.lru.MoveToFront(elem)
		return
	}

	for c.CurrentMem+size > c.MemoryLimit && c.lru.Len() > 0 {
		backElem := c.lru.Back()
		entry := backElem.Value.(*Entry)
		delete(c.cache, entry.Key)
		c.CurrentMem -= entry.Size
		c.lru.Remove(backElem)
	}

	entry := &Entry{
		Key:  key,
		Data: data,
		Size: size,
	}
	elem := c.lru.PushFront(entry)
	c.cache[key] = elem
	c.CurrentMem += size
}

func (c *LRU) evictOldest() {
	if c.lru.Len() == 0 {
		return
	}
	elem := c.lru.Back()
	entry := elem.Value.(*Entry)
	delete(c.cache, entry.Key)
	c.CurrentMem -= entry.Size
	c.lru.Remove(elem)
}

func (c *LRU) Invalidate(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, exists := c.cache[key]; exists {
		entry := elem.Value.(*Entry)
		c.CurrentMem -= entry.Size
		delete(c.cache, key)
		c.lru.Remove(elem)
	}
}

func (c *LRU) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache = make(map[string]*list.Element)
	c.lru = list.New()
	c.CurrentMem = 0
}

func (c *LRU) Stats() (hits, misses int64, size int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.Hits, c.Misses, c.CurrentMem
}

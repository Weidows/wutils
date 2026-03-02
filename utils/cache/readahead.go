package cache

import (
	"sync"
)

// ReadaheadConfig holds configuration for the readahead cache.
type ReadaheadConfig struct {
	// BlockSize is the size of each block in bytes (default: 4KB)
	BlockSize int64
	// PrefetchSize is the number of blocks to prefetch ahead (default: 4)
	PrefetchSize int
	// MemoryLimit is the maximum memory the underlying LRU cache can use
	MemoryLimit int64
}

// DefaultReadaheadConfig returns a config with default values.
func DefaultReadaheadConfig() *ReadaheadConfig {
	return &ReadaheadConfig{
		BlockSize:    4 * 1024,         // 4KB default block size
		PrefetchSize: 4,                // prefetch 4 blocks ahead by default
		MemoryLimit:  64 * 1024 * 1024, // 64MB default memory limit
	}
}

// readPosition tracks the current read position for a file.
type readPosition struct {
	File   string
	Offset int64
}

// ReadaheadCache is an LRU cache with read-ahead prefetch capabilities.
// It monitors sequential read patterns and prefetches upcoming blocks
// in the background to improve read performance.
type ReadaheadCache struct {
	lru        *LRU
	config     *ReadaheadConfig
	mu         sync.RWMutex
	positions  map[string]*readPosition // tracks current read positions per file
	prefetchCh chan *prefetchRequest
	stopCh     chan struct{}
	wg         sync.WaitGroup
}

// prefetchRequest represents a prefetch request for a file block.
type prefetchRequest struct {
	File   string
	Offset int64
	Key    string
}

// NewReadaheadCache creates a new readahead cache with the given configuration.
// If config is nil, DefaultReadaheadConfig() is used.
func NewReadaheadCache(config *ReadaheadConfig) *ReadaheadCache {
	if config == nil {
		config = DefaultReadaheadConfig()
	}

	rc := &ReadaheadCache{
		lru:        NewLRU(config.MemoryLimit),
		config:     config,
		positions:  make(map[string]*readPosition),
		prefetchCh: make(chan *prefetchRequest, 64), // buffered channel for batching
		stopCh:     make(chan struct{}),
	}

	// Start background prefetch worker
	rc.wg.Add(1)
	go rc.prefetchWorker()

	return rc
}

// Get retrieves data from the cache and triggers prefetch if sequential access is detected.
// The key should be formatted as "filename:offset" for readahead to work properly.
func (rc *ReadaheadCache) Get(key string) ([]byte, bool) {
	// Try to get from cache first
	data, found := rc.lru.Get(key)
	if found {
		rc.updatePosition(key)
		return data, true
	}
	return nil, false
}

// Set stores data in the cache.
func (rc *ReadaheadCache) Set(key string, data []byte) {
	rc.lru.Set(key, data)
}

// GetWithPrefetch retrieves data from cache and triggers prefetch for sequential reads.
// keyFormat: "filename:offset" - e.g., "data.txt:4096"
// Returns the cached data and whether it was found.
func (rc *ReadaheadCache) GetWithPrefetch(key string) ([]byte, bool) {
	data, found := rc.lru.Get(key)
	if found {
		rc.updatePositionAndPrefetch(key)
	}
	return data, found
}

// updatePositionAndPrefetch extracts file and offset from key, updates position,
// detects sequential access, and triggers prefetch if needed.
func (rc *ReadaheadCache) updatePositionAndPrefetch(key string) {
	file, offset := rc.parseKey(key)
	if file == "" {
		return
	}

	rc.mu.Lock()
	prevPos, exists := rc.positions[file]
	isSequential := false

	if exists && offset == prevPos.Offset+rc.config.BlockSize {
		// Sequential access detected - current offset equals previous offset + block size
		isSequential = true
	}

	// Update position
	rc.positions[file] = &readPosition{File: file, Offset: offset}
	rc.mu.Unlock()

	// Trigger prefetch if sequential access detected
	if isSequential {
		rc.triggerPrefetch(file, offset)
	}
}

// updatePosition updates the read position for prefetch tracking without triggering prefetch.
func (rc *ReadaheadCache) updatePosition(key string) {
	file, offset := rc.parseKey(key)
	if file == "" {
		return
	}

	rc.mu.Lock()
	rc.positions[file] = &readPosition{File: file, Offset: offset}
	rc.mu.Unlock()
}

// parseKey extracts file path and offset from a key formatted as "filename:offset".
func (rc *ReadaheadCache) parseKey(key string) (string, int64) {
	// Key format: "filename:offset"
	// Find the last colon to handle paths with colons (e.g., C:\path\file.txt)
	for i := len(key) - 1; i >= 0; i-- {
		if key[i] == ':' {
			if i+1 < len(key) {
				var offset int64
				for j := i + 1; j < len(key); j++ {
					if key[j] < '0' || key[j] > '9' {
						break
					}
					offset = offset*10 + int64(key[j]-'0')
				}
				return key[:i], offset
			}
			break
		}
	}
	return "", 0
}

// triggerPrefetch sends prefetch requests for the next N blocks.
func (rc *ReadaheadCache) triggerPrefetch(file string, currentOffset int64) {
	for i := 1; i <= rc.config.PrefetchSize; i++ {
		prefetchOffset := currentOffset + int64(i)*rc.config.BlockSize
		key := rc.makeKey(file, prefetchOffset)

		// Only prefetch if not already in cache
		if _, exists := rc.lru.Get(key); !exists {
			select {
			case rc.prefetchCh <- &prefetchRequest{
				File:   file,
				Offset: prefetchOffset,
				Key:    key,
			}:
			default:
				// Channel full, skip prefetch request
			}
		}
	}
}

// makeKey creates a cache key from file path and offset.
func (rc *ReadaheadCache) makeKey(file string, offset int64) string {
	return file + ":" + int64ToString(offset)
}

// int64ToString converts int64 to string efficiently.
func int64ToString(n int64) string {
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	neg := n < 0
	if neg {
		n = -n
	}
	for n > 0 {
		i--
		buf[i] = byte(n%10) + '0'
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}

// prefetchWorker is a background goroutine that processes prefetch requests.
// In a real implementation, this would read data from disk. For now, it
// demonstrates the prefetch architecture - actual data loading would need
// to be implemented with a reader callback or similar mechanism.
func (rc *ReadaheadCache) prefetchWorker() {
	defer rc.wg.Done()

	for {
		select {
		case <-rc.stopCh:
			return
		case req := <-rc.prefetchCh:
			// In a real implementation, we would:
			// 1. Read the block from the underlying storage
			// 2. Store it in the LRU cache
			//
			// For now, this demonstrates the prefetch trigger mechanism.
			// The actual data loading would require a storage reader to be
			// injected into the ReadaheadCache.
			_ = req
		}
	}
}

// Prefetch explicitly prefetches a specific block. This can be called
// manually when the read pattern is known ahead of time.
func (rc *ReadaheadCache) Prefetch(file string, offset int64) {
	key := rc.makeKey(file, offset)
	if _, exists := rc.lru.Get(key); !exists {
		select {
		case rc.prefetchCh <- &prefetchRequest{
			File:   file,
			Offset: offset,
			Key:    key,
		}:
		default:
		}
	}
}

// Invalidate removes a specific key from the cache.
func (rc *ReadaheadCache) Invalidate(key string) {
	rc.lru.Invalidate(key)
}

// Clear clears all cached data and resets positions.
func (rc *ReadaheadCache) Clear() {
	rc.lru.Clear()
	rc.mu.Lock()
	rc.positions = make(map[string]*readPosition)
	rc.mu.Unlock()
}

// Stats returns cache statistics (hits, misses, current size).
func (rc *ReadaheadCache) Stats() (hits, misses int64, size int64) {
	return rc.lru.Stats()
}

// GetBlockSize returns the configured block size.
func (rc *ReadaheadCache) GetBlockSize() int64 {
	return rc.config.BlockSize
}

// GetPrefetchSize returns the configured prefetch size.
func (rc *ReadaheadCache) GetPrefetchSize() int {
	return rc.config.PrefetchSize
}

// Close stops the prefetch worker and cleans up resources.
func (rc *ReadaheadCache) Close() {
	close(rc.stopCh)
	rc.wg.Wait()
	close(rc.prefetchCh)
}

// SetPrefetchData is called by the prefetch worker to store prefetched data.
// In a full implementation, this would be used by the background goroutine
// to store data after reading from disk.
func (rc *ReadaheadCache) SetPrefetchData(key string, data []byte) {
	rc.lru.Set(key, data)
}

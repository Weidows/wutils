package config

import (
	"sync"
	"time"

	"github.com/jinzhu/configor"
)

// Watcher reloads configuration periodically, supporting hot-reload.
type Watcher struct {
	mu       sync.RWMutex
	cfg      *Config
	path     string
	interval time.Duration
	stopCh   chan struct{}
	doneCh   chan struct{}
}

// NewWatcher creates a new configuration watcher that reloads
// from the given path at the specified interval.
func NewWatcher(path string, interval time.Duration) *Watcher {
	w := &Watcher{
		path:     path,
		interval: interval,
		stopCh:   make(chan struct{}),
		doneCh:   make(chan struct{}),
	}
	// Initial load
	cfg, err := Load(path)
	if err == nil {
		w.cfg = cfg
	}
	return w
}

// Start begins the hot-reload loop in a background goroutine.
func (w *Watcher) Start() {
	go w.loop()
}

// Stop signals the reload loop to exit and waits for it to finish.
func (w *Watcher) Stop() {
	close(w.stopCh)
	<-w.doneCh
}

// Get returns the current configuration snapshot.
func (w *Watcher) Get() *Config {
	w.mu.RLock()
	defer w.mu.RUnlock()
	if w.cfg == nil {
		cfg := DefaultConfig()
		return &cfg
	}
	cp := *w.cfg
	return &cp
}

// RefreshInterval returns the configured refresh interval.
func (w *Watcher) RefreshInterval() time.Duration {
	return w.interval
}

func (w *Watcher) loop() {
	defer close(w.doneCh)

	// Enforce minimum refresh interval
	if w.interval < time.Second {
		w.interval = time.Second
	}

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			var cfg Config
			if err := configor.Load(&cfg, w.path); err == nil {
				w.mu.Lock()
				w.cfg = &cfg
				w.mu.Unlock()
			}
		case <-w.stopCh:
			return
		}
	}
}

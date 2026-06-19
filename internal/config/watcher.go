package config

import (
	"os"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

// ConfigProvider is the interface for accessing the current config.
// Both the Watcher and simple wrappers can implement it, allowing
// services to receive live config updates without knowing the reload mechanism.
type ConfigProvider interface {
	// Get returns the current configuration.
	Get() *Config
}

// Watcher reloads configuration periodically, supporting hot-reload.
// It implements ConfigProvider so services can use it for live config.
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
// Implements ConfigProvider.
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
	if w.interval < time.Second {
		return time.Second
	}
	return w.interval
}

func (w *Watcher) loop() {
	defer close(w.doneCh)

	if w.interval < time.Second {
		w.interval = time.Second
	}

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.reload()
		case <-w.stopCh:
			return
		}
	}
}

func (w *Watcher) reload() {
	data, err := os.ReadFile(w.path)
	if err != nil {
		return
	}
	base := DefaultConfig()
	var fileCfg Config
	if err := yaml.Unmarshal(data, &fileCfg); err != nil {
		return
	}
	base.Merge(&fileCfg)
	RunMigrations(&base)
	if err := base.Validate(); err != nil {
		return
	}

	// Save back if upgraded
	if base.ConfigVersion > fileCfg.ConfigVersion {
		if err := Save(&base, w.path); err != nil {
			// Non-fatal
		}
	}

	w.mu.Lock()
	w.cfg = &base
	w.mu.Unlock()
}

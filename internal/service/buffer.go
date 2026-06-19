package service

import (
	"fmt"
	"sync"

	"github.com/Weidows/wutils/cmd/wutils/buffer"
	"github.com/Weidows/wutils/internal/app"
)

// BufferService wraps the Dokan-based IO buffer filesystem as a service.
type BufferService struct {
	mu      sync.Mutex
	config  *buffer.BufferConfig
	drive   string
	status  app.ServiceStatus
}

// NewBufferService creates a new buffer service.
func NewBufferService(drive string, config *buffer.BufferConfig) *BufferService {
	return &BufferService{
		drive:  drive,
		config: config,
		status: app.StatusStopped,
	}
}

func (s *BufferService) Name() string        { return "buffer" }
func (s *BufferService) Description() string { return "Buffer filesystem — IO 缓冲虚拟文件系统" }

func (s *BufferService) Status() app.ServiceStatus {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.status
}

func (s *BufferService) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.status == app.StatusRunning {
		return fmt.Errorf("buffer: already mounted")
	}

	if err := buffer.Mount(s.drive, s.config); err != nil {
		s.status = app.StatusError
		return fmt.Errorf("buffer mount failed: %w", err)
	}

	s.status = app.StatusRunning
	return nil
}

func (s *BufferService) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.status != app.StatusRunning {
		return nil
	}

	if err := buffer.Unmount(); err != nil {
		s.status = app.StatusError
		return fmt.Errorf("buffer unmount failed: %w", err)
	}

	s.status = app.StatusStopped
	return nil
}

// MountedDrive returns the drive letter where the buffer is mounted.
func (s *BufferService) MountedDrive() string {
	return s.drive
}

// Config returns the current buffer configuration.
func (s *BufferService) Config() *buffer.BufferConfig {
	return s.config
}

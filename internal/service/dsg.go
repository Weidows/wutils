package service

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Weidows/wutils/internal/app"
	"github.com/Weidows/wutils/internal/config"
	"github.com/Weidows/wutils/internal/i18n"
	"github.com/Weidows/wutils/utils/collection"
	"github.com/sirupsen/logrus"
)

// DSGService prevents hard drives from sleeping by periodically writing timestamps.
type DSGService struct {
	mu       sync.Mutex
	cfg      config.ConfigProvider
	logger   *logrus.Logger
	cancel   context.CancelFunc
	status   app.ServiceStatus
}

// NewDSGService creates a new DSG service.
// cfg provides live configuration — use a *config.Watcher for hot-reload support,
// or a simple wrapper for static config.
func NewDSGService(cfg config.ConfigProvider, logger *logrus.Logger) *DSGService {
	if logger == nil {
		logger = logrus.New()
	}
	return &DSGService{
		cfg:    cfg,
		logger: logger,
		status: app.StatusStopped,
	}
}

func (s *DSGService) Name() string        { return "dsg" }
func (s *DSGService) Description() string { return i18n.G("dsg.description") }

func (s *DSGService) Status() app.ServiceStatus {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.status
}

func (s *DSGService) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.status == app.StatusRunning {
		return nil
	}

	// Check disks at start time
	cfg := s.cfg.Get()
	if len(cfg.Cmd.Dsg.Disk) == 0 {
		return fmt.Errorf("dsg: no disks configured")
	}

	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel
	s.status = app.StatusRunning

	go s.run(ctx)
	return nil
}

func (s *DSGService) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.status != app.StatusRunning {
		return nil
	}

	s.cancel()
	s.status = app.StatusStopped
	return nil
}

func (s *DSGService) run(ctx context.Context) {
	defer func() {
		s.mu.Lock()
		s.status = app.StatusStopped
		s.mu.Unlock()
	}()

	for {
		// Read config on each iteration to pick up hot-reload changes
		cfg := s.cfg.Get()
		dsgCfg := cfg.Cmd.Dsg

		s.logger.Infoln(dsgCfg)
		delay := time.Duration(dsgCfg.Delay) * time.Second

		// Write timestamps to all configured disks
		collection.ForEach(dsgCfg.Disk, func(i int, disk string) {
			writeDSGTimestamp(disk, s.logger)
		})

		// Wait for the delay or cancellation
		select {
		case <-ctx.Done():
			return
		case <-time.After(delay):
		}
	}
}

func writeDSGTimestamp(disk string, logger *logrus.Logger) {
	f := strings.Join([]string{disk, ".dsg"}, "/")
	file, err := os.OpenFile(f, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.Println("disk format error, please input like 'E:'", err)
		return
	}
	defer file.Close()
	_, _ = file.WriteString("dsg running at " + time.Now().String() + "\n")
}

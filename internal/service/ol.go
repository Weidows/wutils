package service

import (
	"context"
	"sync"
	"time"

	"github.com/Weidows/wutils/internal/app"
	"github.com/Weidows/wutils/internal/config"
	"github.com/Weidows/wutils/utils/collection"
	"github.com/Weidows/wutils/utils/grammar"
	os2 "github.com/Weidows/wutils/utils/os"
	"github.com/sirupsen/logrus"
)

// OLService monitors running windows and applies transparency based on matching rules.
type OLService struct {
	mu     sync.Mutex
	cfg    *config.OLConfig
	logger *logrus.Logger
	cancel context.CancelFunc
	status app.ServiceStatus
}

// NewOLService creates a new Opacity Listener service.
func NewOLService(cfg *config.OLConfig, logger *logrus.Logger) *OLService {
	if logger == nil {
		logger = logrus.New()
	}
	return &OLService{
		cfg:    cfg,
		logger: logger,
		status: app.StatusStopped,
	}
}

func (s *OLService) Name() string        { return "ol" }
func (s *OLService) Description() string { return "Opacity Listener — 窗口透明度控制" }

func (s *OLService) Status() app.ServiceStatus {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.status
}

func (s *OLService) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.status == app.StatusRunning {
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel
	s.status = app.StatusRunning

	go s.run(ctx)
	return nil
}

func (s *OLService) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.status != app.StatusRunning {
		return nil
	}

	s.cancel()
	s.status = app.StatusStopped
	return nil
}

// ListWindows returns all visible windows with their titles and handles.
func (s *OLService) ListWindows() []*os2.EnumWindowsResult {
	return os2.GetEnumWindowsInfo(&os2.EnumWindowsFilter{
		IgnoreNoTitled:  true,
		IgnoreInvisible: true,
	})
}

func (s *OLService) run(ctx context.Context) {
	defer func() {
		s.mu.Lock()
		s.status = app.StatusStopped
		s.mu.Unlock()
	}()

	s.logger.Infoln(s.cfg)
	delay := time.Duration(s.cfg.Delay) * time.Second

	for {
		select {
		case <-ctx.Done():
			return
		default:
			s.applyPatterns()
			time.Sleep(delay)
		}
	}
}

func (s *OLService) applyPatterns() {
	windows := os2.GetEnumWindowsInfo(&os2.EnumWindowsFilter{
		IgnoreNoTitled:  true,
		IgnoreInvisible: true,
	})

	collection.ForEach(windows, func(i int, window *os2.EnumWindowsResult) {
		collection.ForEach(s.cfg.Patterns, func(ii int, pattern config.OLPattern) {
			if grammar.Match(pattern.Title, window.Title) && pattern.Opacity != window.Opacity {
				os2.SetWindowOpacity(window.Handle, pattern.Opacity)
			}
		})
	})
}

package service

import (
	"github.com/Weidows/wutils/cmd/wutils/gmm"
	"github.com/Weidows/wutils/internal/app"
)

// GMMService manages Go module proxy mirrors (GOPROXY, GOSUMDB).
type GMMService struct{}

// NewGMMService creates a new GMMService.
func NewGMMService() *GMMService {
	return &GMMService{}
}

func (s *GMMService) Name() string        { return "gmm" }
func (s *GMMService) Description() string { return "Go 代理管理 — 测试/切换模块镜像" }
func (s *GMMService) Status() app.ServiceStatus { return app.StatusStopped }
func (s *GMMService) Start() error               { return nil }
func (s *GMMService) Stop() error                { return nil }

// TestSpeed benchmarks all proxy mirrors.
func (s *GMMService) TestSpeed() {
	gmm.TestSpeed()
}

// SetProxy configures GOPROXY to the named mirror.
func (s *GMMService) SetProxy(name string) error {
	return gmm.SetProxy(name)
}

// SetSumdb configures GOSUMDB to the named mirror.
func (s *GMMService) SetSumdb(name string) error {
	return gmm.SetSumdb(name)
}

// List prints all available mirrors.
func (s *GMMService) List() {
	gmm.List()
}

// Current prints the current go environment proxy configuration.
func (s *GMMService) Current() {
	gmm.Current()
}

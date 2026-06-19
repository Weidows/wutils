package service

import (
	"testing"

	"github.com/Weidows/wutils/internal/app"
	"github.com/Weidows/wutils/internal/config"
)

// testCfgProvider wraps a config for testing.
type testCfgProvider struct {
	cfg *config.Config
}

func (p *testCfgProvider) Get() *config.Config { return p.cfg }

func TestDSGService(t *testing.T) {
	cfg := config.DefaultConfig()
	p := &testCfgProvider{cfg: &cfg}
	s := NewDSGService(p, nil)
	if s.Name() != "dsg" {
		t.Errorf("expected dsg, got %s", s.Name())
	}
	if s.Status() != app.StatusStopped {
		t.Errorf("expected stopped, got %s", s.Status())
	}
}

func TestOLService(t *testing.T) {
	cfg := config.DefaultConfig()
	p := &testCfgProvider{cfg: &cfg}
	s := NewOLService(p, nil)
	if s.Name() != "ol" {
		t.Errorf("expected ol, got %s", s.Name())
	}
	if s.Status() != app.StatusStopped {
		t.Errorf("expected stopped, got %s", s.Status())
	}
}

func TestDiffService(t *testing.T) {
	s := NewDiffService()
	if s.Name() != "diff" {
		t.Errorf("expected diff, got %s", s.Name())
	}
	if s.Status() != app.StatusStopped {
		t.Errorf("expected stopped, got %s", s.Status())
	}
	if err := s.Start(); err != nil {
		t.Errorf("Start should be no-op: %v", err)
	}
	if err := s.Stop(); err != nil {
		t.Errorf("Stop should be no-op: %v", err)
	}
}

func TestZipCrackService(t *testing.T) {
	s := NewZipCrackService()
	if s.Name() != "zipcrack" {
		t.Errorf("expected zipcrack, got %s", s.Name())
	}
}

func TestMediaService(t *testing.T) {
	s := NewMediaService()
	if s.Name() != "media" {
		t.Errorf("expected media, got %s", s.Name())
	}
}

func TestExtractService(t *testing.T) {
	s := NewExtractService()
	if s.Name() != "extract" {
		t.Errorf("expected extract, got %s", s.Name())
	}
}

func TestGMMService(t *testing.T) {
	s := NewGMMService()
	if s.Name() != "gmm" {
		t.Errorf("expected gmm, got %s", s.Name())
	}
}

func TestServiceRegistryIntegration(t *testing.T) {
	r := app.NewServiceRegistry()

	cfg := config.DefaultConfig()
	p := &testCfgProvider{cfg: &cfg}
	dsg := NewDSGService(p, nil)
	ol := NewOLService(p, nil)
	diff := NewDiffService()

	if err := r.Register(dsg); err != nil {
		t.Fatalf("register dsg: %v", err)
	}
	if err := r.Register(ol); err != nil {
		t.Fatalf("register ol: %v", err)
	}
	if err := r.Register(diff); err != nil {
		t.Fatalf("register diff: %v", err)
	}

	svcs := r.List()
	if len(svcs) != 3 {
		t.Fatalf("expected 3 services, got %d", len(svcs))
	}

	if errs := r.StartAll(); len(errs) > 0 {
		t.Fatalf("start all errors: %v", errs)
	}

	dsgStatus := dsg.Status()
	if dsgStatus != app.StatusRunning && dsgStatus != app.StatusError {
		t.Logf("DSG status after start: %s (expected running or error in test env)", dsgStatus)
	}

	if errs := r.StopAll(); len(errs) > 0 {
		t.Fatalf("stop all errors: %v", errs)
	}
}

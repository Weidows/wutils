package service

import (
	"testing"

	"github.com/Weidows/wutils/internal/app"
	"github.com/Weidows/wutils/internal/config"
)

func TestDiffService(t *testing.T) {
	s := NewDiffService()
	if s.Name() != "diff" {
		t.Errorf("expected diff, got %s", s.Name())
	}
	if s.Status() != app.StatusStopped {
		t.Errorf("expected stopped, got %s", s.Status())
	}
	// Start/Stop are no-ops
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

func TestDSGService(t *testing.T) {
	cfg := config.DefaultConfig()
	s := NewDSGService(&cfg.Cmd.Dsg, nil)
	if s.Name() != "dsg" {
		t.Errorf("expected dsg, got %s", s.Name())
	}
	if s.Status() != app.StatusStopped {
		t.Errorf("expected stopped, got %s", s.Status())
	}
}

func TestOLService(t *testing.T) {
	cfg := config.DefaultConfig()
	s := NewOLService(&cfg.Cmd.Ol, nil)
	if s.Name() != "ol" {
		t.Errorf("expected ol, got %s", s.Name())
	}
	if s.Status() != app.StatusStopped {
		t.Errorf("expected stopped, got %s", s.Status())
	}
}

func TestServiceRegistryIntegration(t *testing.T) {
	r := app.NewServiceRegistry()

	cfg := config.DefaultConfig()
	dsg := NewDSGService(&cfg.Cmd.Dsg, nil)
	ol := NewOLService(&cfg.Cmd.Ol, nil)
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

	// Start all
	if errs := r.StartAll(); len(errs) > 0 {
		t.Fatalf("start all errors: %v", errs)
	}

	// DSG should now be running (or try to start)
	dsgStatus := dsg.Status()
	if dsgStatus != app.StatusRunning && dsgStatus != app.StatusError {
		t.Logf("DSG status after start: %s (expected running or error in test env)", dsgStatus)
	}

	// Stop all
	if errs := r.StopAll(); len(errs) > 0 {
		t.Fatalf("stop all errors: %v", errs)
	}
}

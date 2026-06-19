package app

import (
	"errors"
	"sync"
	"testing"
)

// testService is a simple Service implementation for testing.
type testService struct {
	mu          sync.Mutex
	name        string
	description string
	status      ServiceStatus
	startErr    error
	stopErr     error
}

func (s *testService) Name() string               { return s.name }
func (s *testService) Description() string         { return s.description }
func (s *testService) Status() ServiceStatus       { s.mu.Lock(); defer s.mu.Unlock(); return s.status }
func (s *testService) Start() error                { s.mu.Lock(); defer s.mu.Unlock(); s.status = StatusRunning; return s.startErr }
func (s *testService) Stop() error                 { s.mu.Lock(); defer s.mu.Unlock(); s.status = StatusStopped; return s.stopErr }

func newTestService(name string) *testService {
	return &testService{
		name:        name,
		description: "test service " + name,
		status:      StatusStopped,
	}
}

func TestServiceStatusString(t *testing.T) {
	tests := []struct {
		status ServiceStatus
		want   string
	}{
		{StatusStopped, "stopped"},
		{StatusRunning, "running"},
		{StatusError, "error"},
		{ServiceStatus(99), "unknown"},
	}

	for _, tt := range tests {
		if got := tt.status.String(); got != tt.want {
			t.Errorf("ServiceStatus(%d).String() = %q, want %q", tt.status, got, tt.want)
		}
	}
}

func TestRegistryRegister(t *testing.T) {
	r := NewServiceRegistry()
	s1 := newTestService("alpha")
	s2 := newTestService("beta")

	if err := r.Register(s1); err != nil {
		t.Fatalf("Register alpha failed: %v", err)
	}
	if err := r.Register(s2); err != nil {
		t.Fatalf("Register beta failed: %v", err)
	}
	// Duplicate
	if err := r.Register(s1); err == nil {
		t.Fatal("expected error for duplicate registration")
	}
}

func TestRegistryGet(t *testing.T) {
	r := NewServiceRegistry()
	s := newTestService("test-svc")
	r.Register(s)

	got, ok := r.Get("test-svc")
	if !ok {
		t.Fatal("service not found")
	}
	if got.Name() != "test-svc" {
		t.Errorf("expected test-svc, got %s", got.Name())
	}

	_, ok = r.Get("nonexistent")
	if ok {
		t.Fatal("expected false for nonexistent service")
	}
}

func TestRegistryList(t *testing.T) {
	r := NewServiceRegistry()
	r.Register(newTestService("a"))
	r.Register(newTestService("b"))
	r.Register(newTestService("c"))

	svcs := r.List()
	if len(svcs) != 3 {
		t.Fatalf("expected 3 services, got %d", len(svcs))
	}
	// Check insertion order
	if svcs[0].Name() != "a" || svcs[1].Name() != "b" || svcs[2].Name() != "c" {
		t.Errorf("unexpected order: %v", svcs)
	}
}

func TestRegistryNames(t *testing.T) {
	r := NewServiceRegistry()
	r.Register(newTestService("z"))
	r.Register(newTestService("a"))
	r.Register(newTestService("m"))

	names := r.Names()
	if len(names) != 3 {
		t.Fatalf("expected 3 names, got %d", len(names))
	}
	// Should be alphabetically sorted
	if names[0] != "a" || names[1] != "m" || names[2] != "z" {
		t.Errorf("expected sorted names, got %v", names)
	}
}

func TestStartStopAll(t *testing.T) {
	r := NewServiceRegistry()
	r.Register(newTestService("svc1"))
	r.Register(newTestService("svc2"))

	errs := r.StartAll()
	if len(errs) > 0 {
		t.Fatalf("StartAll errors: %v", errs)
	}

	svc, _ := r.Get("svc1")
	if svc.Status() != StatusRunning {
		t.Error("expected svc1 to be running")
	}

	errs = r.StopAll()
	if len(errs) > 0 {
		t.Fatalf("StopAll errors: %v", errs)
	}

	if svc.Status() != StatusStopped {
		t.Error("expected svc1 to be stopped")
	}
}

func TestStartStopErrors(t *testing.T) {
	r := NewServiceRegistry()
	badSvc := newTestService("bad")
	badSvc.startErr = errors.New("start failed")
	badSvc.stopErr = errors.New("stop failed")
	r.Register(badSvc)
	r.Register(newTestService("good"))

	// Start should report the error but not panic
	errs := r.StartAll()
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d: %v", len(errs), errs)
	}

	// Stop should also handle errors
	errs = r.StopAll()
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d: %v", len(errs), errs)
	}
}

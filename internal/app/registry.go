package app

import (
	"fmt"
	"sort"
	"sync"
)

// ServiceRegistry manages the lifecycle of all wutils services.
type ServiceRegistry struct {
	mu       sync.RWMutex
	services map[string]Service
	order    []string // insertion order for deterministic iteration
}

// NewServiceRegistry creates an empty service registry.
func NewServiceRegistry() *ServiceRegistry {
	return &ServiceRegistry{
		services: make(map[string]Service),
	}
}

// Register adds a service to the registry.
// Returns an error if a service with the same name already exists.
func (r *ServiceRegistry) Register(svc Service) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := svc.Name()
	if _, exists := r.services[name]; exists {
		return fmt.Errorf("service %q already registered", name)
	}
	r.services[name] = svc
	r.order = append(r.order, name)
	return nil
}

// Get retrieves a service by name.
func (r *ServiceRegistry) Get(name string) (Service, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	svc, ok := r.services[name]
	return svc, ok
}

// List returns all registered services in insertion order.
func (r *ServiceRegistry) List() []Service {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]Service, 0, len(r.order))
	for _, name := range r.order {
		result = append(result, r.services[name])
	}
	return result
}

// Names returns all registered service names in alphabetical order.
func (r *ServiceRegistry) Names() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.services))
	for name := range r.services {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// StartAll starts every registered service concurrently.
func (r *ServiceRegistry) StartAll() []error {
	r.mu.RLock()
	svcs := make([]Service, 0, len(r.services))
	for _, name := range r.order {
		svcs = append(svcs, r.services[name])
	}
	r.mu.RUnlock()

	var wg sync.WaitGroup
	errCh := make(chan error, len(svcs))

	for _, svc := range svcs {
		wg.Add(1)
		go func(s Service) {
			defer wg.Done()
			if err := s.Start(); err != nil {
				errCh <- fmt.Errorf("%s: %w", s.Name(), err)
			}
		}(svc)
	}

	wg.Wait()
	close(errCh)

	var errs []error
	for err := range errCh {
		errs = append(errs, err)
	}
	return errs
}

// StopAll stops every running service concurrently.
func (r *ServiceRegistry) StopAll() []error {
	r.mu.RLock()
	svcs := make([]Service, 0, len(r.services))
	for _, name := range r.order {
		svcs = append(svcs, r.services[name])
	}
	r.mu.RUnlock()

	var wg sync.WaitGroup
	errCh := make(chan error, len(svcs))

	for _, svc := range svcs {
		wg.Add(1)
		go func(s Service) {
			defer wg.Done()
			if err := s.Stop(); err != nil {
				errCh <- fmt.Errorf("%s: %w", s.Name(), err)
			}
		}(svc)
	}

	wg.Wait()
	close(errCh)

	var errs []error
	for err := range errCh {
		errs = append(errs, err)
	}
	return errs
}

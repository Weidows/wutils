// Package app defines the core abstractions for wutils services.
//
// The Service interface provides a uniform lifecycle contract that both
// CLI commands and GUI controls consume without knowing each other's details.
package app

// ServiceStatus represents the current state of a service.
type ServiceStatus int

const (
	StatusStopped ServiceStatus = iota
	StatusRunning
	StatusError
)

func (s ServiceStatus) String() string {
	switch s {
	case StatusStopped:
		return "stopped"
	case StatusRunning:
		return "running"
	case StatusError:
		return "error"
	default:
		return "unknown"
	}
}

// Service is the uniform lifecycle interface for all wutils services.
//
// Long-running services (DSG, OL, Buffer) use Start/Stop for lifecycle.
// One-shot services (Diff, ZipCrack, etc.) can be invoked via their
// specific methods while implementing Service as a thin wrapper.
type Service interface {
	// Name returns the human-readable service name.
	Name() string

	// Description returns a short description of what the service does.
	Description() string

	// Start begins the service. For one-shot services this is a no-op.
	Start() error

	// Stop gracefully stops the service.
	Stop() error

	// Status returns the current service status.
	Status() ServiceStatus
}

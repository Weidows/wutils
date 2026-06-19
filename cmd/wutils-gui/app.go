package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Weidows/wutils/internal/app"
	"github.com/Weidows/wutils/internal/config"
	"github.com/Weidows/wutils/internal/service"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx       context.Context
	registry  *app.ServiceRegistry
	dsgSvc    *service.DSGService
	olSvc     *service.OLService
	diffSvc   *service.DiffService
	zipSvc    *service.ZipCrackService
	mediaSvc  *service.MediaService
	extractSvc *service.ExtractService
	gmmSvc    *service.GMMService
}

// NewApp creates a new App application struct
func NewApp() *App {
	registry := app.NewServiceRegistry()

	cfg := config.DefaultConfig()

	dsgSvc := service.NewDSGService(&cfg.Cmd.Dsg, nil)
	olSvc := service.NewOLService(&cfg.Cmd.Ol, nil)
	diffSvc := service.NewDiffService()
	zipSvc := service.NewZipCrackService()
	mediaSvc := service.NewMediaService()
	extractSvc := service.NewExtractService()
	gmmSvc := service.NewGMMService()

	// Register all services
	_ = registry.Register(dsgSvc)
	_ = registry.Register(olSvc)
	_ = registry.Register(diffSvc)
	_ = registry.Register(zipSvc)
	_ = registry.Register(mediaSvc)
	_ = registry.Register(extractSvc)
	_ = registry.Register(gmmSvc)

	return &App{
		registry:   registry,
		dsgSvc:     dsgSvc,
		olSvc:      olSvc,
		diffSvc:    diffSvc,
		zipSvc:     zipSvc,
		mediaSvc:   mediaSvc,
		extractSvc: extractSvc,
		gmmSvc:     gmmSvc,
	}
}

// startup is called when the app starts.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// shutdown cleans up all services.
func (a *App) shutdown(ctx context.Context) {
	a.registry.StopAll()
}

// beforeClose intercepts the window close button — hides to tray instead of quitting.
// Returns true to prevent close, false to allow it.
func (a *App) beforeClose(ctx context.Context) bool {
	runtime.WindowHide(ctx)
	return true // prevent close, just hide
}

// ============ Service Registry (for Dashboard) ============

// ServiceInfo represents a service's state for the frontend.
type ServiceInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

// ListServices returns all registered services with their status.
func (a *App) ListServices() []ServiceInfo {
	svcs := a.registry.List()
	result := make([]ServiceInfo, len(svcs))
	for i, svc := range svcs {
		result[i] = ServiceInfo{
			Name:        svc.Name(),
			Description: svc.Description(),
			Status:      svc.Status().String(),
		}
	}
	return result
}

// StartService starts a service by name.
func (a *App) StartService(name string) error {
	svc, ok := a.registry.Get(name)
	if !ok {
		return fmt.Errorf("service %q not found", name)
	}
	return svc.Start()
}

// StopService stops a service by name.
func (a *App) StopService(name string) error {
	svc, ok := a.registry.Get(name)
	if !ok {
		return fmt.Errorf("service %q not found", name)
	}
	return svc.Stop()
}

// StartAllServices starts all registered services.
func (a *App) StartAllServices() []string {
	errs := a.registry.StartAll()
	if len(errs) == 0 {
		return []string{"All services started"}
	}
	result := make([]string, len(errs))
	for i, err := range errs {
		result[i] = err.Error()
	}
	return result
}

// StopAllServices stops all running services.
func (a *App) StopAllServices() []string {
	errs := a.registry.StopAll()
	if len(errs) == 0 {
		return []string{"All services stopped"}
	}
	result := make([]string, len(errs))
	for i, err := range errs {
		result[i] = err.Error()
	}
	return result
}

// ============ DSG Service ============

// StartDSG starts the Disk Sleep Guard.
func (a *App) StartDSG() error {
	return a.dsgSvc.Start()
}

// StopDSG stops the Disk Sleep Guard.
func (a *App) StopDSG() error {
	return a.dsgSvc.Stop()
}

// DSGStatus returns the current status of DSG.
func (a *App) DSGStatus() string {
	return a.dsgSvc.Status().String()
}

// ============ OL Service ============

// StartOL starts the Opacity Listener.
func (a *App) StartOL() error {
	return a.olSvc.Start()
}

// StopOL stops the Opacity Listener.
func (a *App) StopOL() error {
	return a.olSvc.Stop()
}

// OLStatus returns the current status of OL.
func (a *App) OLStatus() string {
	return a.olSvc.Status().String()
}

// ListWindows returns all visible windows.
func (a *App) ListWindows() []interface{} {
	windows := a.olSvc.ListWindows()
	result := make([]interface{}, len(windows))
	for i, w := range windows {
		result[i] = map[string]interface{}{
			"handle": fmt.Sprintf("%d", w.Handle),
			"title":  w.Title,
			"opacity": w.Opacity,
		}
	}
	return result
}

// ============ Diff Service ============

// DiffResult holds the result of a diff operation.
type DiffResult struct {
	MissInA []string `json:"missInA"`
	MissInB []string `json:"missInB"`
}

// Diff compares two files and returns the differences.
func (a *App) Diff(inputA, inputB string) DiffResult {
	missInA, missInB := a.diffSvc.Diff(inputA, inputB)
	return DiffResult{MissInA: missInA, MissInB: missInB}
}

// ============ ZipCrack Service ============

// CrackPassword attempts to crack an archive password.
func (a *App) CrackPassword(archivePath string) string {
	return a.zipSvc.CrackPassword(archivePath)
}

// ============ Media Service ============

// ClusterMedia clusters media files by time and GPS.
func (a *App) ClusterMedia(inputDir string) string {
	a.mediaSvc.ClusterAndCopy(inputDir)
	return "done"
}

// ============ Extract Service ============

// Extract flattens subdirectories.
func (a *App) Extract(mode, rootPath string) error {
	return a.extractSvc.Run(mode, rootPath)
}

// ============ GMM Service ============

// GMMTestResult represents a mirror test result.
type GMMTestResult struct {
	Name   string `json:"name"`
	URL    string `json:"url"`
	LatencyMS int64 `json:"latencyMs"`
}

// TestMirrors tests all proxy mirrors.
func (a *App) TestMirrors() []GMMTestResult {
	// Note: gmm.TestSpeed only prints to stdout. For GUI we need structured data.
	// This is a simplified version - the full version would parse structured results.
	return nil
}

// SetGoproxy sets the GOPROXY environment variable.
func (a *App) SetGoproxy(name string) error {
	return a.gmmSvc.SetProxy(name)
}

// SetGosumdb sets the GOSUMDB environment variable.
func (a *App) SetGosumdb(name string) error {
	return a.gmmSvc.SetSumdb(name)
}

// ============ Utility ============

// Greet returns a greeting for the given name (kept from template for testing).
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, welcome to wutils!", name)
}

// Now returns the current server time.
func (a *App) Now() string {
	return time.Now().Format(time.RFC3339)
}

package service

import (
	"github.com/Weidows/wutils/cmd/wutils/extract"
	"github.com/Weidows/wutils/internal/app"
	"github.com/Weidows/wutils/internal/i18n"
)

// ExtractService flattens subdirectories by moving their contents to the parent directory.
type ExtractService struct{}

// NewExtractService creates a new ExtractService.
func NewExtractService() *ExtractService {
	return &ExtractService{}
}

func (s *ExtractService) Name() string        { return "extract" }
func (s *ExtractService) Description() string { return i18n.G("extract.description") }
func (s *ExtractService) Status() app.ServiceStatus { return app.StatusStopped }
func (s *ExtractService) Start() error               { return nil }
func (s *ExtractService) Stop() error                { return nil }

// Run executes the extraction with the given mode and root path.
func (s *ExtractService) Run(mode, rootPath string) error {
	return extract.Run(mode, rootPath)
}

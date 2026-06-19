package service

import (
	"github.com/Weidows/wutils/cmd/wutils/diff"
	"github.com/Weidows/wutils/internal/app"
	"github.com/Weidows/wutils/internal/i18n"
)

// DiffService compares two files and returns the line differences.
type DiffService struct{}

// NewDiffService creates a new DiffService.
func NewDiffService() *DiffService {
	return &DiffService{}
}

func (s *DiffService) Name() string        { return "diff" }
func (s *DiffService) Description() string { return i18n.G("diff.description") }
func (s *DiffService) Status() app.ServiceStatus { return app.StatusStopped }
func (s *DiffService) Start() error               { return nil }
func (s *DiffService) Stop() error                { return nil }

// Diff returns the sets of lines missing in each file (symmetric difference).
func (s *DiffService) Diff(inputA, inputB string) (missInA, missInB []string) {
	return diff.CheckLinesDiff(inputA, inputB)
}

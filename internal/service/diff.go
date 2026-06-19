package service

import (
	"github.com/Weidows/wutils/cmd/wutils/diff"
	"github.com/Weidows/wutils/internal/app"
)

// DiffService compares two files and returns the line differences.
type DiffService struct{}

// NewDiffService creates a new DiffService.
func NewDiffService() *DiffService {
	return &DiffService{}
}

func (s *DiffService) Name() string        { return "diff" }
func (s *DiffService) Description() string { return "文件对比 — 计算两个文件的行差集" }
func (s *DiffService) Status() app.ServiceStatus { return app.StatusStopped }
func (s *DiffService) Start() error               { return nil }
func (s *DiffService) Stop() error                { return nil }

// Diff returns the sets of lines missing in each file (symmetric difference).
func (s *DiffService) Diff(inputA, inputB string) (missInA, missInB []string) {
	return diff.CheckLinesDiff(inputA, inputB)
}

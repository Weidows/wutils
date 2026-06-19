package service

import (
	"bufio"
	"os"
	"strings"

	"github.com/Weidows/wutils/internal/app"
	"github.com/Weidows/wutils/internal/i18n"
	"github.com/Weidows/wutils/utils/collection"
)

// DiffService compares two files and returns the line differences.
type DiffService struct{}

// NewDiffService creates a new DiffService.
func NewDiffService() *DiffService {
	return &DiffService{}
}

func (s *DiffService) Name() string               { return "diff" }
func (s *DiffService) Description() string        { return i18n.G("diff.description") }
func (s *DiffService) Status() app.ServiceStatus  { return app.StatusStopped }
func (s *DiffService) Start() error               { return nil }
func (s *DiffService) Stop() error                { return nil }

// Diff returns the sets of lines missing in each file (symmetric difference).
func (s *DiffService) Diff(inputA, inputB string) (missInA, missInB []string) {
	fileA, _ := readLines(inputA)
	fileB, _ := readLines(inputB)
	return collection.SliceDiff(fileA, fileB)
}

func readLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "] [")
		if len(parts) > 1 {
			lines = append(lines, strings.TrimSpace(parts[1]))
		}
	}
	return lines, scanner.Err()
}

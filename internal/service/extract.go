package service

import (
	"log"
	"os"
	"path"
	"strings"

	"github.com/Weidows/wutils/internal/app"
	"github.com/Weidows/wutils/internal/i18n"
	"github.com/Weidows/wutils/utils/files"
	"github.com/Weidows/wutils/utils/hash"
)

// ExtractService flattens subdirectories by moving their contents to the parent directory.
type ExtractService struct{}

// NewExtractService creates a new ExtractService.
func NewExtractService() *ExtractService {
	return &ExtractService{}
}

func (s *ExtractService) Name() string               { return "extract" }
func (s *ExtractService) Description() string        { return i18n.G("extract.description") }
func (s *ExtractService) Status() app.ServiceStatus  { return app.StatusStopped }
func (s *ExtractService) Start() error               { return nil }
func (s *ExtractService) Stop() error                { return nil }

const (
	modeAutoCheck = "0"
	modeOverwrite = "1"
	modeSkip      = "2"
)

// Run executes the extraction with the given mode and root path.
func (s *ExtractService) Run(mode, rootPath string) error {
	log.Println("Root: ", rootPath)
	if err := os.Chdir(rootPath); err != nil {
		return err
	}

	items, err := os.ReadDir(rootPath)
	if err != nil {
		return err
	}

	var dirs []os.DirEntry
	for _, item := range items {
		if item.IsDir() {
			dirs = append(dirs, item)
		}
	}

	for _, dir := range dirs {
		subItems, err := os.ReadDir(dir.Name())
		if err != nil {
			return err
		}

		needToDelete := dir.Name()
		for _, sub := range subItems {
			oldPath := path.Join(dir.Name(), sub.Name())
			newPath := path.Join(rootPath, sub.Name())
			if fi, err := os.Stat(newPath); err == nil {
				if fi.IsDir() {
					if dir.Name() == sub.Name() {
						newPath = path.Join(rootPath, "tmp-"+sub.Name())
					} else {
						files.MergeDirs(oldPath, newPath)
						continue
					}
				} else {
					switch mode {
					case modeAutoCheck:
						if hash.CompareFile(oldPath, newPath) {
							continue
						} else {
							newPath = path.Join(rootPath, dir.Name()+"-"+sub.Name())
						}
					case modeOverwrite:
						if err = os.RemoveAll(newPath); err != nil {
							return err
						}
					case modeSkip:
						continue
					}
				}
			}
			if err = os.Rename(oldPath, newPath); err != nil {
				return err
			}
		}
		if err = os.RemoveAll(needToDelete); err != nil {
			return err
		}
	}

	// Clean up tmp- prefixes
	items, err = os.ReadDir(rootPath)
	if err != nil {
		return err
	}
	for _, item := range items {
		if item.IsDir() && strings.Contains(item.Name(), "tmp-") {
			_ = os.Rename(
				path.Join(rootPath, item.Name()),
				path.Join(rootPath, strings.TrimPrefix(item.Name(), "tmp-")),
			)
		}
	}
	return nil
}

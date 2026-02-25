package extract

import (
	"os"
	"path"
	"path/filepath"
	"testing"
)

func TestExtract(t *testing.T) {
	wd, _ := os.Getwd()
	rootPath := filepath.Join(wd, "testdata", "1")

	err := Run(AutoCheck, rootPath)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(path.Join(rootPath, "2.3.txt")); os.IsNotExist(err) {
		t.Error("expected 2.3.txt to exist")
	}

	if _, err := os.Stat(path.Join(rootPath, "3.1.txt")); os.IsNotExist(err) {
		t.Error("expected 3.1.txt to exist (deduplicated)")
	}

	if _, err := os.Stat(path.Join(rootPath, "3.2.txt")); os.IsNotExist(err) {
		t.Error("expected 3.2.txt to exist")
	}
}

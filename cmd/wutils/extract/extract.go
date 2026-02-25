package extract

import (
	"log"
	"os"
	"path"
	"strings"

	"github.com/Weidows/wutils/utils/files"
	"github.com/Weidows/wutils/utils/hash"
)

const (
	AutoCheck = "0"
	Overwrite = "1"
	Skip      = "2"
)

func Run(mode, rootPath string) error {
	log.Println("Root: ", rootPath)

	if err := os.Chdir(rootPath); err != nil {
		return err
	}

	floor2xItems, err := os.ReadDir(rootPath)
	if err != nil {
		return err
	}

	var floor2xDirs []os.DirEntry
	for _, floor2xItem := range floor2xItems {
		if floor2xItem.IsDir() {
			floor2xDirs = append(floor2xDirs, floor2xItem)
		}
	}

	for _, floor2xDir := range floor2xDirs {
		floor3xItems, err := os.ReadDir(floor2xDir.Name())
		if err != nil {
			return err
		}

		needToDelete := floor2xDir.Name()
		for _, floor3xItem := range floor3xItems {
			oldPath := path.Join(floor2xDir.Name(), floor3xItem.Name())
			newPath := path.Join(rootPath, floor3xItem.Name())
			if fileInfo, err := os.Stat(newPath); err == nil {
				if fileInfo.IsDir() {
					if floor2xDir.Name() == floor3xItem.Name() {
						newPath = path.Join(rootPath, "tmp-"+floor3xItem.Name())
					} else {
						files.MergeDirs(oldPath, newPath)
						continue
					}
				} else {
					switch mode {
					case AutoCheck:
						if hash.CompareFile(oldPath, newPath) {
							continue
						} else {
							newPath = path.Join(rootPath, floor2xDir.Name()+"-"+floor3xItem.Name())
						}
					case Overwrite:
						err = os.RemoveAll(newPath)
						if err != nil {
							return err
						}
					case Skip:
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

	floor2xItems, err = os.ReadDir(rootPath)
	if err != nil {
		return err
	}

	for _, floor2xItem := range floor2xItems {
		if floor2xItem.IsDir() && strings.Contains(floor2xItem.Name(), "tmp-") {
			if err = os.Rename(path.Join(rootPath, floor2xItem.Name()), path.Join(rootPath, strings.Trim(floor2xItem.Name(), "tmp-"))); err != nil {
				return err
			}
		}
	}
	return nil
}

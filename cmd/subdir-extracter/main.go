package main

import (
	"github.com/Weidows/wutils/utils/files"
	"github.com/Weidows/wutils/utils/hash"
	"log"
	"os"
	"path"
	"strings"
)

const (
	autoCheck = "0"
	overwrite = "1"
	skip      = "2"
)

func main() {
	if len(os.Args) < 2 {
		return
	}
	mode := os.Args[1]
	floor1x := os.Args[2]
	log.Println("Root: ", floor1x)

	if err := os.Chdir(floor1x); err != nil {
		log.Fatal(err)
	}

	// floor2xDirs 就是你需要的所有子目录
	floor2xItems, err := os.ReadDir(floor1x)
	if err != nil {
		log.Fatal(err)
	}

	var floor2xDirs []os.DirEntry
	for _, floor2xItem := range floor2xItems {
		if floor2xItem.IsDir() {
			floor2xDirs = append(floor2xDirs, floor2xItem)
		}
	}

	// 2.x
	for _, floor2xDir := range floor2xDirs {
		floor3xItems, err := os.ReadDir(floor2xDir.Name())
		if err != nil {
			log.Fatal(err)
		}

		// 3.x
		// move items
		needToDelete := floor2xDir.Name()
		for _, floor3xItem := range floor3xItems {
			oldPath := path.Join(floor2xDir.Name(), floor3xItem.Name())
			newPath := path.Join(floor1x, floor3xItem.Name())
			// check if newPath exists
			if fileInfo, err := os.Stat(newPath); err == nil {
				// 文件夹已存在
				if fileInfo.IsDir() {
					// 父子目录同名
					if floor2xDir.Name() == floor3xItem.Name() {
						newPath = path.Join(floor1x, "tmp-"+floor3xItem.Name())
					} else {
						files.MergeDirs(oldPath, newPath)
						continue
					}
				} else {
					switch mode {
					case autoCheck:
						// 是重复文件
						if hash.CompareFile(oldPath, newPath) {
							continue
						} else {
							newPath = path.Join(floor1x, floor2xDir.Name()+"-"+floor3xItem.Name())
						}
					case overwrite:
						err = os.RemoveAll(newPath)
						if err != nil {
							return
						}
					case skip:
						continue
					}
				}
			}

			if err = os.Rename(oldPath, newPath); err != nil {
				log.Fatal(err)
			}

		}

		// delete floor2xDir and skipped floor2xItems in it
		if err = os.RemoveAll(needToDelete); err != nil {
			log.Fatal(err)
		}
	}

	floor2xItems, err = os.ReadDir(floor1x)
	if err != nil {
		log.Fatal(err)
	}

	for _, floor2xItem := range floor2xItems {
		if floor2xItem.IsDir() && strings.Contains(floor2xItem.Name(), "tmp-") {
			if err = os.Rename(path.Join(floor1x, floor2xItem.Name()), path.Join(floor1x, strings.Trim(floor2xItem.Name(), "tmp-"))); err != nil {
				log.Fatal(err)
			}
		}
	}
}

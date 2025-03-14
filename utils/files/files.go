package files

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/Weidows/wutils/utils/collection"
	"github.com/Weidows/wutils/utils/hash"
	"github.com/Weidows/wutils/utils/log"
)

var (
	logger = log.GetLogger()
)

const (
	HARD_MOVE = 1
	SOFT_MOVE = 2
)

type subFileInfo struct {
	Name string
	Path string
}

func GetSubFiles(path string) (res []*subFileInfo) {
	info, err := os.Stat(path)
	if err != nil || !info.IsDir() {
		return nil
	}

	items, err := os.ReadDir(path)
	if err != nil {
		return nil
	}
	for _, item := range items {
		if !item.IsDir() {
			res = append(res, &subFileInfo{Name: item.Name(), Path: path})
		}
	}
	return res
}

func GetAllSubFiles(path string) []*subFileInfo {
	return getAllSubFiles(path, "")
}

func GetAllSubFilesWithFilter(path string, filter func(*subFileInfo) bool) (res []*subFileInfo) {
	collection.ForEach(getAllSubFiles(path, ""), func(index int, value *subFileInfo) {
		if filter(value) {
			res = append(res, value)
		}
	})
	return res
}

// https://www.jianshu.com/p/92794526fb94?u_atoken=53f2b298-ccca-46c7-9337-b4cf591b2f8f&u_asession=01_qLoBDaRk1M5SeGCRkdHBmPCNvgyscrsoOYpIgpjlbCmOnG34feDPT2CjStSe46IX0KNBwm7Lovlpxjd_P_q4JsKWYrT3W_NKPr8w6oU7K9f7RG5F6-A8wD78mhCffewhUF3o-sVtq6Wun3JL3SJe2BkFo3NEHBv0PZUm6pbxQU&u_asig=053rs23s8kYv45len_QlH5siBvqlpH_kD8gNu2A-4fbA7CigNwanUzsRmVw813XC7AC66GRf0tk5x_5n5HSPIDZ7kM1AX_K0GmXaaINzotIcJwlg2eIsFD8yAT9merTtzp1AbC3-yfHDX2cCGVRkie6w4nVhgD9IIuQKfmOJFukUj9JS7q8ZD7Xtz2Ly-b0kmuyAKRFSVJkkdwVUnyHAIJzcGuoo7G0WovhEe-1j2-gibPIDU0TpaVMUjlPAfVPQvSom7nzSzR1LP16f45fIKp-e3h9VXwMyh6PgyDIVSG1W8T7cfbN2A6Rljwpks-diHLBQwMsil22Y1Tlmo5XouMtJmZflDcW3Ox19TjFEuTMtYJro_EbU9SYnrqkfuUuDGamWspDxyAEEo4kbsryBKb9Q&u_aref=200R4W8ckdeoQMCjbmEw8cxwrek%3D
func getAllSubFiles(path, file string) (res []*subFileInfo) {
	fullPath := filepath.Join(path, file)
	dstF, err := os.Open(fullPath)
	if err != nil {
		logger.Error(err)
	}
	defer dstF.Close()
	fileInfo, err := dstF.Stat()
	if err != nil {
		logger.Error(err)
	}

	if fileInfo.IsDir() {
		fileList, err := dstF.Readdir(0) //获取文件夹下各个文件或文件夹的fileInfo
		if err != nil {
			logger.Error(err)
		}
		collection.ForEach(fileList, func(index int, subFiles os.FileInfo) {
			collection.ForEach(getAllSubFiles(fullPath, subFiles.Name()), func(index int, resFileInfo *subFileInfo) {
				res = append(res, resFileInfo)
			})
		})
	} else {
		res = append(res, &subFileInfo{Name: file, Path: fullPath})
	}
	return res
}

// TODO
// func WithOpen() {

// }

func MergeDirs(mergePath, distPath string) {
	mergePath = filepath.Clean(mergePath)
	distPath = filepath.Clean(distPath)
	if mergePath == distPath || !IsDir(mergePath) || !IsDir(distPath) {
		return
	}

	Move(mergePath, distPath, HARD_MOVE)
	err := os.RemoveAll(mergePath)
	if err != nil {
		return
	}
}

func IsDir(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		logger.Error(err)
	}
	return fileInfo.IsDir()
}

// IsExist 检查文件是否存在
func IsExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true // 文件存在
	}
	if os.IsNotExist(err) {
		return false // 文件不存在
	}
	return false // 其他错误
}

func Move(oldPath, newPath string, mode int) {
	switch mode {
	case HARD_MOVE:
		HardMove(oldPath, newPath)
		//case SOFT_MOVE:
		//	SoftMove(oldPath, newPath)
	}
}

func HardMove(oldPath, newPath string) {
	oldPath = filepath.Clean(oldPath)
	newPath = filepath.Clean(newPath)
	if !IsExist(oldPath) {
		return
	}

	if IsExist(newPath) {
		if IsDir(newPath) {
			items, err := os.ReadDir(oldPath)
			if err != nil {
				return
			}
			for _, item := range items {
				HardMove(filepath.Join(oldPath, item.Name()), filepath.Join(newPath, item.Name()))
			}
		} else {
			//	file
			if hash.CompareFile(oldPath, newPath) {
				_ = os.Remove(oldPath)
			} else {
				filename := filepath.Base(newPath)
				newPath = strings.Replace(newPath, filename, "merged-"+filename, -1)
			}
		}
	}
	err := os.Rename(oldPath, newPath)
	if err != nil {
		return
	}
}

/*
CopyFile Copy file with timestamp

复制文件, 同时恢复时间戳
*/
func CopyFile(src, dst string) error {
	// 打开源文件
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// 获取源文件的信息，包括时间戳
	sourceInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// 创建目标文件
	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// 复制文件内容
	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	// 恢复时间戳
	err = os.Chtimes(dst, sourceInfo.ModTime(), sourceInfo.ModTime())
	if err != nil {
		return err
	}

	return nil
}

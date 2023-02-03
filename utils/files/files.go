package files

import (
	"github.com/Weidows/Golang/utils/collection"
	"github.com/Weidows/Golang/utils/log"
	"os"
	"path/filepath"
)

var (
	logger = log.GetLogger()
)

type subFileInfo struct {
	Name string
	Path string
}

func GetSubFiles(path string) []*subFileInfo {
	return getSubFiles(path, "")
}

func GetSubFilesWithFilter(path string, filter func(*subFileInfo) bool) (res []*subFileInfo) {
	collection.ForEach(getSubFiles(path, ""), func(index int, value *subFileInfo) {
		if filter(value) {
			res = append(res, value)
		}
	})
	return res
}

// https://www.jianshu.com/p/92794526fb94?u_atoken=53f2b298-ccca-46c7-9337-b4cf591b2f8f&u_asession=01_qLoBDaRk1M5SeGCRkdHBmPCNvgyscrsoOYpIgpjlbCmOnG34feDPT2CjStSe46IX0KNBwm7Lovlpxjd_P_q4JsKWYrT3W_NKPr8w6oU7K9f7RG5F6-A8wD78mhCffewhUF3o-sVtq6Wun3JL3SJe2BkFo3NEHBv0PZUm6pbxQU&u_asig=053rs23s8kYv45len_QlH5siBvqlpH_kD8gNu2A-4fbA7CigNwanUzsRmVw813XC7AC66GRf0tk5x_5n5HSPIDZ7kM1AX_K0GmXaaINzotIcJwlg2eIsFD8yAT9merTtzp1AbC3-yfHDX2cCGVRkie6w4nVhgD9IIuQKfmOJFukUj9JS7q8ZD7Xtz2Ly-b0kmuyAKRFSVJkkdwVUnyHAIJzcGuoo7G0WovhEe-1j2-gibPIDU0TpaVMUjlPAfVPQvSom7nzSzR1LP16f45fIKp-e3h9VXwMyh6PgyDIVSG1W8T7cfbN2A6Rljwpks-diHLBQwMsil22Y1Tlmo5XouMtJmZflDcW3Ox19TjFEuTMtYJro_EbU9SYnrqkfuUuDGamWspDxyAEEo4kbsryBKb9Q&u_aref=200R4W8ckdeoQMCjbmEw8cxwrek%3D
func getSubFiles(path, file string) (res []*subFileInfo) {
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
			collection.ForEach(getSubFiles(fullPath, subFiles.Name()), func(index int, resFileInfo *subFileInfo) {
				res = append(res, resFileInfo)
			})
		})
	} else {
		res = append(res, &subFileInfo{Name: file, Path: fullPath})
	}
	return res
}

func WithOpen() {

}

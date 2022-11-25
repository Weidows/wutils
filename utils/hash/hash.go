package hash

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"github.com/Weidows/Golang/utils/log"
	"io"
	"os"
)

const (
	Sha256 = "sha256"
	Md5    = "md5"
)

func SumString(s, fn string) (res string) {
	switch fn {
	case Sha256:
		hash := sha256.Sum256([]byte(s))
		res = fmt.Sprintf("%x", hash)
	case Md5:
		hash := md5.Sum([]byte(s))
		res = fmt.Sprintf("%x", hash)
	default:
		hash := sha256.Sum256([]byte(s))
		res = fmt.Sprintf("%x", hash)
	}
	return
}

// SumFile
//
// by: https://gist.github.com/miguelmota/1f398eb5fb2666a037e6ed6fc9b119b0
func SumFile(filePath, fn string) (res string) {
	switch fn {
	case Sha256:
		file, err := os.Open(filePath)
		if err != nil {
			log.GetLogger().Error(err)
		}
		hash := sha256.New()
		if _, err = io.Copy(hash, file); err != nil {
			log.GetLogger().Error(err)
		}
		res = fmt.Sprintf("%s", hash.Sum(nil))
	}
	return
}

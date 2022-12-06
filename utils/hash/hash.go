package hash

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"github.com/Weidows/Golang/utils/log"
	"github.com/sirupsen/logrus"
	"hash"
	"io"
	"os"
)

const (
	Sha256 = iota
	Md5
)

var (
	logger *logrus.Logger
)

func init() {
	logger = log.GetLogger()
}

func SumString[S interface{ string | []byte }](str S, choice int) string {
	var h any
	switch choice {
	case Sha256:
		h = sha256.Sum256([]byte(str))
	case Md5:
		h = md5.Sum([]byte(str))
	default:
		h = sha256.Sum256([]byte(str))
	}

	return fmt.Sprintf("%x", h)
}

// SumFile
//
// by: https://gist.github.com/miguelmota/1f398eb5fb2666a037e6ed6fc9b119b0
func SumFile(filePath string, choice int) string {
	file, err := os.Open(filePath)
	if err != nil {
		logger.Error(err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	var h hash.Hash
	switch choice {
	case Sha256:
		h = sha256.New()
	case Md5:
		h = md5.New()
	}

	if _, err = io.Copy(h, file); err != nil {
		logger.Error(err)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

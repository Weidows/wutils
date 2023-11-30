package os

import (
	"github.com/Weidows/wutils/utils/log"
	"os"
)

var (
	logger = log.GetLogger()
)

func GetCurrentPath() string {
	getwd, err := os.Getwd()
	if err != nil {
		logger.Error(err.Error())
	}
	return getwd
}

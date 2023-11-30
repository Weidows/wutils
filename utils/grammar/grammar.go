package grammar

import (
	"github.com/Weidows/wutils/utils/log"
	"regexp"
)

var (
	logger = log.GetLogger()
)

func ConditionalEqual[T any](condition bool, value1, value2 T) T {
	if condition {
		return value1
	} else {
		return value2
	}
}

func Match(regex, text string) bool {
	isMatch, err := regexp.Match(regex, []byte(text))
	if err != nil {
		logger.Error(err.Error())
	}
	return isMatch
}

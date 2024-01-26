package grammar

import (
	"regexp"

	"github.com/Weidows/wutils/utils/log"
)

var (
	logger = log.GetLogger()
)

func ConditionalEqual[T any](condition bool, value1, value2 T) T {
	if condition {
		return value1
	}
	return value2
}

func Match(regex, text string) bool {
	isMatch, err := regexp.Match(regex, []byte(text))
	if err != nil {
		logger.Error(err.Error())
	}
	return isMatch
}

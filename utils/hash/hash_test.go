package hash

import (
	"testing"
)

func TestSumString(t *testing.T) {
	logger.Println(SumString("abc", Sha256))
	logger.Println(SumString("abc", Md5))
}

func TestSumFile(t *testing.T) {
	logger.Println(SumFile("hash.go", Sha256))
	logger.Println(SumFile("hash.go", Md5))
}

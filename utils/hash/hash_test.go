package hash

import (
	"log"
	"testing"
)

func TestSumString(t *testing.T) {
	log.Println(SumString("abc", Sha256))
}

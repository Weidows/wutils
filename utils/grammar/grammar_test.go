package grammar

import (
	"fmt"
	"testing"
)

func TestConditionalEqual(t *testing.T) {
	fmt.Println(ConditionalEqual(true, 1, 2))
	fmt.Println(ConditionalEqual(false, 1, 2))
}

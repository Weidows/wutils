package diff

import (
	"fmt"
	"testing"
)

func TestCheckDiff(t *testing.T) {
	missInA, missInB := CheckLinesDiff("./test/inputA.txt", "./test/inputB.txt")

	// 输出结果
	fmt.Println("Missing in A")
	for _, file := range missInA {
		fmt.Println(file)
	}

	fmt.Println("\nMissing in B")
	for _, file := range missInB {
		fmt.Println(file)
	}
}

package diff

import (
	"bufio"
	"os"
	"strings"

	"github.com/Weidows/wutils/utils/collection"
)

func CheckLinesDiff(inputA, inputB string) ([]string, []string) {
	inputFileA, _ := readInputFile(inputA)
	inputFileB, _ := readInputFile(inputB)

	// 比较两个文件名列表
	return collection.SliceDiff(inputFileA, inputFileB)
}

func readInputFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var fileNames []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// 提取文件名
		parts := strings.Split(line, "] [")
		if len(parts) > 1 {
			fileName := strings.TrimSpace(parts[1])
			fileNames = append(fileNames, fileName)
		}
	}
	return fileNames, scanner.Err()
}

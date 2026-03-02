package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/Weidows/wutils/utils/bufioutil"
)

func main() {
	fmt.Println("=== Buffer Performance Test ===")
	fmt.Println()

	const fileCount = 1000
	const fileSize = 512

	fmt.Printf("Generating %d random files of %d bytes each...\n", fileCount, fileSize)
	fmt.Println()

	// Test 1: Direct write
	tmpDir1, _ := os.MkdirTemp("", "direct_test")
	defer os.RemoveAll(tmpDir1)

	start := time.Now()
	for i := 0; i < fileCount; i++ {
		path := filepath.Join(tmpDir1, fmt.Sprintf("file_%04d.txt", i))
		data := make([]byte, fileSize)
		rand.Read(data)
		os.WriteFile(path, data, 0644)
	}
	directTime := time.Since(start)
	directSeconds := directTime.Seconds()
	if directSeconds > 0 {
		fmt.Printf("Direct write:  %v (%d files/sec)\n", directTime, int(float64(fileCount)/directSeconds))
	} else {
		fmt.Printf("Direct write:  %v (too fast to measure)\n", directTime)
	}

	// Test 2: Buffered write
	tmpDir2, _ := os.MkdirTemp("", "buffered_test")
	defer os.RemoveAll(tmpDir2)

	strategy := bufioutil.GetPresetStrategy(bufioutil.UseCaseMigration)
	pool := bufioutil.NewPool(100*1024*1024, strategy)

	start = time.Now()
	for i := 0; i < fileCount; i++ {
		path := filepath.Join(tmpDir2, fmt.Sprintf("file_%04d.txt", i))
		data := make([]byte, fileSize)
		rand.Read(data)
		pool.Write(path, data, 0)
	}
	pool.FlushAll()
	bufferedTime := time.Since(start)
	fmt.Printf("Buffered write: %v (%d files/sec)\n", bufferedTime, int(float64(fileCount)/bufferedTime.Seconds()))

	fmt.Println()
	fmt.Printf("Speedup: %.2fx\n", float64(directTime)/float64(bufferedTime))
	fmt.Printf("Memory used: %d bytes\n", pool.MemoryUsage())
}

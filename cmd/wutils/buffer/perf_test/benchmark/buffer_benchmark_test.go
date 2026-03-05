package main

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Weidows/wutils/utils/bufioutil"
	"github.com/Weidows/wutils/utils/cache"
)

const (
	blockSize   = 4 * 1024
	iterations  = 1000
	memoryLimit = 100 * 1024 * 1024
)

type benchmarkResult struct {
	name       string
	ops        int
	bytes      int64
	duration   time.Duration
	iops       float64
	throughput float64
}

func generateRandomData(size int) []byte {
	data := make([]byte, size)
	rand.Read(data)
	return data
}

func createTestFiles(b testing.TB, dir string, count int, size int) []string {
	paths := make([]string, count)
	for i := 0; i < count; i++ {
		path := filepath.Join(dir, fmt.Sprintf("file_%04d.dat", i))
		data := generateRandomData(size)
		if err := os.WriteFile(path, data, 0644); err != nil {
			b.Fatalf("Failed to create test file: %v", err)
		}
		paths[i] = path
	}
	return paths
}

func BenchmarkSmallFileRandomWriteUnbuffered(b *testing.B) {
	tmpDir, err := os.MkdirTemp("", "benchmark_write_unbuffered")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	b.ResetTimer()
	b.SetBytes(int64(blockSize * iterations))

	for i := 0; i < b.N; i++ {
		for j := 0; j < iterations; j++ {
			path := filepath.Join(tmpDir, fmt.Sprintf("file_%04d.dat", j))
			data := generateRandomData(blockSize)
			if err := os.WriteFile(path, data, 0644); err != nil {
				b.Fatalf("Failed to write file: %v", err)
			}
		}
		for j := 0; j < iterations; j++ {
			path := filepath.Join(tmpDir, fmt.Sprintf("file_%04d.dat", j))
			os.Remove(path)
		}
	}
}

func BenchmarkSmallFileRandomWriteBuffered(b *testing.B) {
	tmpDir, err := os.MkdirTemp("", "benchmark_write_buffered")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	strategy := bufioutil.GetPresetStrategy(bufioutil.UseCaseMigration)
	pool := bufioutil.NewPool(memoryLimit, strategy)
	defer pool.Close()

	b.ResetTimer()
	b.SetBytes(int64(blockSize * iterations))

	for i := 0; i < b.N; i++ {
		for j := 0; j < iterations; j++ {
			path := filepath.Join(tmpDir, fmt.Sprintf("file_%04d.dat", j))
			data := generateRandomData(blockSize)
			pool.Write(path, data, 0)
		}
		pool.FlushAll()
		for j := 0; j < iterations; j++ {
			path := filepath.Join(tmpDir, fmt.Sprintf("file_%04d.dat", j))
			os.Remove(path)
		}
	}
}

func BenchmarkSmallFileRandomReadUnbuffered(b *testing.B) {
	tmpDir, err := os.MkdirTemp("", "benchmark_read_unbuffered")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	createTestFiles(b, tmpDir, iterations, blockSize)

	b.ResetTimer()
	b.SetBytes(int64(blockSize * iterations))

	for i := 0; i < b.N; i++ {
		for j := 0; j < iterations; j++ {
			path := filepath.Join(tmpDir, fmt.Sprintf("file_%04d.dat", j))
			data, err := os.ReadFile(path)
			if err != nil {
				b.Fatalf("Failed to read file: %v", err)
			}
			_ = data
		}
	}
}

func BenchmarkSmallFileRandomReadWithLRU(b *testing.B) {
	tmpDir, err := os.MkdirTemp("", "benchmark_read_lru")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	paths := createTestFiles(b, tmpDir, iterations, blockSize)
	lru := cache.NewLRU(memoryLimit)

	b.ResetTimer()
	b.SetBytes(int64(blockSize * iterations))

	for i := 0; i < b.N; i++ {
		for j := 0; j < iterations; j++ {
			path := paths[j]
			key := filepath.Base(path)
			data, found := lru.Get(key)
			if !found {
				data, err := os.ReadFile(path)
				if err != nil {
					b.Fatalf("Failed to read file: %v", err)
				}
				lru.Set(key, data)
			}
			_ = data
		}
	}
}

func BenchmarkSmallFileRandomReadWithReadahead(b *testing.B) {
	tmpDir, err := os.MkdirTemp("", "benchmark_read_readahead")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	paths := createTestFiles(b, tmpDir, iterations, blockSize)

	readaheadConfig := &cache.ReadaheadConfig{
		BlockSize:    blockSize,
		PrefetchSize: 4,
		MemoryLimit:  memoryLimit,
	}
	readahead := cache.NewReadaheadCache(readaheadConfig)
	defer readahead.Close()

	b.ResetTimer()
	b.SetBytes(int64(blockSize * iterations))

	for i := 0; i < b.N; i++ {
		for j := 0; j < iterations; j++ {
			path := paths[j]
			key := fmt.Sprintf("%s:0", filepath.Base(path))
			data, found := readahead.GetWithPrefetch(key)
			if !found {
				data, err := os.ReadFile(path)
				if err != nil {
					b.Fatalf("Failed to read file: %v", err)
				}
				readahead.Set(key, data)
			}
			_ = data
		}
	}
}

func BenchmarkSequentialWriteUnbuffered(b *testing.B) {
	tmpDir, err := os.MkdirTemp("", "benchmark_seq_write_unbuffered")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	fileSize := 4 * 1024 * 1024
	numFiles := 10

	b.ResetTimer()
	b.SetBytes(int64(fileSize * numFiles))

	for i := 0; i < b.N; i++ {
		for j := 0; j < numFiles; j++ {
			path := filepath.Join(tmpDir, fmt.Sprintf("seq_%04d.dat", j))
			data := generateRandomData(fileSize)
			if err := os.WriteFile(path, data, 0644); err != nil {
				b.Fatalf("Failed to write file: %v", err)
			}
		}
		for j := 0; j < numFiles; j++ {
			path := filepath.Join(tmpDir, fmt.Sprintf("seq_%04d.dat", j))
			os.Remove(path)
		}
	}
}

func BenchmarkSequentialWriteBuffered(b *testing.B) {
	tmpDir, err := os.MkdirTemp("", "benchmark_seq_write_buffered")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	fileSize := 4 * 1024 * 1024
	numFiles := 10

	strategy := bufioutil.GetPresetStrategy(bufioutil.UseCaseMigration)
	pool := bufioutil.NewPool(memoryLimit, strategy)
	defer pool.Close()

	b.ResetTimer()
	b.SetBytes(int64(fileSize * numFiles))

	for i := 0; i < b.N; i++ {
		for j := 0; j < numFiles; j++ {
			path := filepath.Join(tmpDir, fmt.Sprintf("seq_%04d.dat", j))
			data := generateRandomData(fileSize)
			pool.Write(path, data, 0)
		}
		pool.FlushAll()
		for j := 0; j < numFiles; j++ {
			path := filepath.Join(tmpDir, fmt.Sprintf("seq_%04d.dat", j))
			os.Remove(path)
		}
	}
}

func BenchmarkSequentialReadUnbuffered(b *testing.B) {
	tmpDir, err := os.MkdirTemp("", "benchmark_seq_read_unbuffered")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	fileSize := 4 * 1024 * 1024
	numFiles := 10

	for j := 0; j < numFiles; j++ {
		path := filepath.Join(tmpDir, fmt.Sprintf("seq_%04d.dat", j))
		data := generateRandomData(fileSize)
		if err := os.WriteFile(path, data, 0644); err != nil {
			b.Fatalf("Failed to create test file: %v", err)
		}
	}

	b.ResetTimer()
	b.SetBytes(int64(fileSize * numFiles))

	for i := 0; i < b.N; i++ {
		for j := 0; j < numFiles; j++ {
			path := filepath.Join(tmpDir, fmt.Sprintf("seq_%04d.dat", j))
			data, err := os.ReadFile(path)
			if err != nil {
				b.Fatalf("Failed to read file: %v", err)
			}
			_ = data
		}
	}
}

func BenchmarkSequentialReadWithReadahead(b *testing.B) {
	tmpDir, err := os.MkdirTemp("", "benchmark_seq_read_readahead")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	fileSize := 4 * 1024 * 1024
	numFiles := 10

	paths := make([]string, numFiles)
	for j := 0; j < numFiles; j++ {
		path := filepath.Join(tmpDir, fmt.Sprintf("seq_%04d.dat", j))
		data := generateRandomData(fileSize)
		if err := os.WriteFile(path, data, 0644); err != nil {
			b.Fatalf("Failed to create test file: %v", err)
		}
		paths[j] = path
	}

	readaheadConfig := &cache.ReadaheadConfig{
		BlockSize:    blockSize,
		PrefetchSize: 8,
		MemoryLimit:  memoryLimit,
	}
	readahead := cache.NewReadaheadCache(readaheadConfig)
	defer readahead.Close()

	b.ResetTimer()
	b.SetBytes(int64(fileSize * numFiles))

	for i := 0; i < b.N; i++ {
		for j := 0; j < numFiles; j++ {
			path := paths[j]
			file, err := os.Open(path)
			if err != nil {
				b.Fatalf("Failed to open file: %v", err)
			}

			var offset int64 = 0
			for offset < int64(fileSize) {
				key := fmt.Sprintf("%s:%d", filepath.Base(path), offset)
				data, found := readahead.GetWithPrefetch(key)
				if !found {
					buf := make([]byte, blockSize)
					n, err := file.ReadAt(buf, offset)
					if err != nil && err != io.EOF {
						file.Close()
						b.Fatalf("Failed to read file: %v", err)
					}
					if n > 0 {
						data = buf[:n]
						readahead.Set(key, data)
					}
				}
				offset += int64(blockSize)
			}
			file.Close()
		}
	}
}

func BenchmarkStreamingWriteUnbuffered(b *testing.B) {
	tmpDir, err := os.MkdirTemp("", "benchmark_stream_write_unbuffered")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	totalSize := 64 * 1024 * 1024
	chunkSize := 32 * 1024

	b.ResetTimer()
	b.SetBytes(int64(totalSize))

	for i := 0; i < b.N; i++ {
		path := filepath.Join(tmpDir, "stream.dat")
		file, err := os.Create(path)
		if err != nil {
			b.Fatalf("Failed to create file: %v", err)
		}

		for written := 0; written < totalSize; written += chunkSize {
			data := generateRandomData(chunkSize)
			if _, err := file.Write(data); err != nil {
				file.Close()
				b.Fatalf("Failed to write: %v", err)
			}
		}
		file.Close()
		os.Remove(path)
	}
}

func BenchmarkStreamingWriteBuffered(b *testing.B) {
	tmpDir, err := os.MkdirTemp("", "benchmark_stream_write_buffered")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	totalSize := 64 * 1024 * 1024
	chunkSize := 32 * 1024

	strategy := bufioutil.GetPresetStrategy(bufioutil.UseCaseMigration)
	pool := bufioutil.NewPool(memoryLimit, strategy)
	defer pool.Close()

	b.ResetTimer()
	b.SetBytes(int64(totalSize))

	for i := 0; i < b.N; i++ {
		path := filepath.Join(tmpDir, "stream.dat")

		for written := 0; written < totalSize; written += chunkSize {
			data := generateRandomData(chunkSize)
			pool.Write(path, data, int64(written))
		}
		pool.FlushAll()
		os.Remove(path)
	}
}

func BenchmarkStreamingReadUnbuffered(b *testing.B) {
	tmpDir, err := os.MkdirTemp("", "benchmark_stream_read_unbuffered")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	totalSize := 64 * 1024 * 1024
	chunkSize := 32 * 1024

	testPath := filepath.Join(tmpDir, "stream.dat")
	file, err := os.Create(testPath)
	if err != nil {
		b.Fatalf("Failed to create test file: %v", err)
	}
	for written := 0; written < totalSize; written += chunkSize {
		data := generateRandomData(chunkSize)
		file.Write(data)
	}
	file.Close()

	b.ResetTimer()
	b.SetBytes(int64(totalSize))

	for i := 0; i < b.N; i++ {
		file, err := os.Open(testPath)
		if err != nil {
			b.Fatalf("Failed to open file: %v", err)
		}

		buf := make([]byte, chunkSize)
		for {
			_, err := file.Read(buf)
			if err != nil {
				if err == io.EOF {
					break
				}
				file.Close()
				b.Fatalf("Failed to read: %v", err)
			}
		}
		file.Close()
	}
}

func BenchmarkStreamingReadWithBuffer(b *testing.B) {
	tmpDir, err := os.MkdirTemp("", "benchmark_stream_read_buffer")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	totalSize := 64 * 1024 * 1024
	chunkSize := 32 * 1024

	testPath := filepath.Join(tmpDir, "stream.dat")
	file, err := os.Create(testPath)
	if err != nil {
		b.Fatalf("Failed to create test file: %v", err)
	}
	for written := 0; written < totalSize; written += chunkSize {
		data := generateRandomData(chunkSize)
		file.Write(data)
	}
	file.Close()

	b.ResetTimer()
	b.SetBytes(int64(totalSize))

	for i := 0; i < b.N; i++ {
		file, err := os.Open(testPath)
		if err != nil {
			b.Fatalf("Failed to open file: %v", err)
		}

		bufReader := make([]byte, 0, chunkSize)
		for {
			tmpBuf := make([]byte, chunkSize)
			n, err := file.Read(tmpBuf)
			if err != nil {
				if err == io.EOF {
					break
				}
				file.Close()
				b.Fatalf("Failed to read: %v", err)
			}
			bufReader = append(bufReader, tmpBuf[:n]...)
		}
		file.Close()
		_ = bufReader
	}
}

func printBenchmarkResults(results []benchmarkResult) {
	fmt.Println("\n=== Benchmark Results Summary ===")
	fmt.Printf("%-40s %12s %15s %12s\n", "Test Name", "Ops", "IOPS", "Throughput")
	fmt.Println(strings.Repeat("-", 80))

	for _, r := range results {
		fmt.Printf("%-40s %12d %15.2f %12.2f MB/s\n",
			r.name, r.ops, r.iops, r.throughput)
	}
}

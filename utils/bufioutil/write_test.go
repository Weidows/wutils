package bufioutil

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestPool_BasicWrite(t *testing.T) {
	pool := NewPool(1000, NewTimeBasedStrategy(time.Hour))

	pool.Write("test.txt", []byte("hello"), 0)
	time.Sleep(50 * time.Millisecond)

	val, _ := pool.Read("test.txt", 0, 5)
	if string(val) != "hello" {
		t.Fatalf("expected hello, got %s", string(val))
	}
}

func TestPool_Flush(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")

	pool := NewPool(1000, NewTimeBasedStrategy(time.Hour))
	pool.Write(testFile, []byte("hello world"), 0)
	time.Sleep(50 * time.Millisecond)

	err := pool.Flush(testFile)
	if err != nil {
		t.Fatalf("flush failed: %v", err)
	}

	data, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("read file failed: %v", err)
	}
	if string(data) != "hello world" {
		t.Fatalf("expected 'hello world', got '%s'", string(data))
	}
}

func TestPool_FlushAll(t *testing.T) {
	tmpDir := t.TempDir()
	file1 := filepath.Join(tmpDir, "test1.txt")
	file2 := filepath.Join(tmpDir, "test2.txt")

	pool := NewPool(1000, NewTimeBasedStrategy(time.Hour))
	pool.Write(file1, []byte("content1"), 0)
	pool.Write(file2, []byte("content2"), 0)
	time.Sleep(50 * time.Millisecond)

	err := pool.FlushAll()
	if err != nil {
		t.Fatalf("flush all failed: %v", err)
	}

	data1, _ := os.ReadFile(file1)
	data2, _ := os.ReadFile(file2)

	if string(data1) != "content1" {
		t.Fatalf("expected content1, got %s", string(data1))
	}
	if string(data2) != "content2" {
		t.Fatalf("expected content2, got %s", string(data2))
	}
}

func TestPool_OffsetWrite(t *testing.T) {
	pool := NewPool(1000, NewTimeBasedStrategy(time.Hour))

	pool.Write("test.txt", []byte("hello"), 0)
	pool.Write("test.txt", []byte(" world"), 5)
	time.Sleep(50 * time.Millisecond)

	val, _ := pool.Read("test.txt", 0, 11)
	if string(val) != "hello world" {
		t.Fatalf("expected 'hello world', got '%s'", string(val))
	}
}

func TestPool_MemoryUsage(t *testing.T) {
	pool := NewPool(1000, NewTimeBasedStrategy(time.Hour))

	initial := pool.MemoryUsage()

	pool.Write("test.txt", []byte("hello world"), 0)
	time.Sleep(50 * time.Millisecond)

	after := pool.MemoryUsage()
	if after <= initial {
		t.Fatal("memory usage should increase after write")
	}
}

func TestPool_FlushEmpty(t *testing.T) {
	pool := NewPool(1000, NewTimeBasedStrategy(time.Hour))

	err := pool.Flush("nonexistent.txt")
	if err != nil {
		t.Fatalf("flush nonexistent should not fail: %v", err)
	}
}

func TestPool_Strategies(t *testing.T) {
	timeStrat := NewTimeBasedStrategy(time.Millisecond)
	if timeStrat.Name() != "time_based" {
		t.Fatalf("expected time_based, got %s", timeStrat.Name())
	}

	sizeStrat := NewSizeBasedStrategy(100)
	if sizeStrat.Name() != "size_based" {
		t.Fatalf("expected size_based, got %s", sizeStrat.Name())
	}

	sortedStrat := NewSortedBatchStrategy(100, time.Second)
	if sortedStrat.Name() != "sorted_batch" {
		t.Fatalf("expected sorted_batch, got %s", sortedStrat.Name())
	}
}

func TestPool_Presets(t *testing.T) {
	tests := []struct {
		useCase   UseCase
		expectNil bool
	}{
		{UseCaseMonitoring, false},
		{UseCaseDefrag, false},
		{UseCaseDownload, false},
		{UseCaseMigration, false},
		{UseCaseBalanced, false},
		{"unknown", false},
	}

	for _, tt := range tests {
		strat := GetPresetStrategy(tt.useCase)
		if strat == nil && !tt.expectNil {
			t.Errorf("expected strategy for %s, got nil", tt.useCase)
		}
		if strat != nil && strat.Name() == "" {
			t.Errorf("strategy name should not be empty for %s", tt.useCase)
		}
	}
}

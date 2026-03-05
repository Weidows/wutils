package bufioutil

import (
	"errors"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

// 错误定义
var (
	ErrDiskFull      = errors.New("磁盘空间不足")
	ErrStorageFailed = errors.New("存储失败")
)

type FlushStrategy interface {
	ShouldFlush(pool *Pool) bool
	Name() string
}

type TimeBasedStrategy struct {
	interval  time.Duration
	lastFlush time.Time
}

func NewTimeBasedStrategy(interval time.Duration) *TimeBasedStrategy {
	return &TimeBasedStrategy{
		interval:  interval,
		lastFlush: time.Now(),
	}
}

func (s *TimeBasedStrategy) ShouldFlush(pool *Pool) bool {
	if time.Since(s.lastFlush) >= s.interval {
		s.lastFlush = time.Now()
		return true
	}
	return false
}

func (s *TimeBasedStrategy) Name() string { return "time_based" }

type SizeBasedStrategy struct {
	threshold int64
}

func NewSizeBasedStrategy(threshold int64) *SizeBasedStrategy {
	return &SizeBasedStrategy{threshold: threshold}
}

func (s *SizeBasedStrategy) ShouldFlush(pool *Pool) bool {
	return pool.MemoryUsage() >= s.threshold
}

func (s *SizeBasedStrategy) Name() string { return "size_based" }

type SortedBatchStrategy struct {
	minSize    int64
	maxLatency time.Duration
	writeCount int
	lastFlush  time.Time
}

func NewSortedBatchStrategy(minSize int64, maxLatency time.Duration) *SortedBatchStrategy {
	return &SortedBatchStrategy{
		minSize:    minSize,
		maxLatency: maxLatency,
		lastFlush:  time.Now(),
	}
}

func (s *SortedBatchStrategy) ShouldFlush(pool *Pool) bool {
	memUsage := pool.MemoryUsage()
	forceFlush := memUsage >= pool.MemoryLimit
	timeSinceFlush := time.Since(s.lastFlush)

	softFlush := (memUsage >= s.minSize && timeSinceFlush >= s.maxLatency) ||
		s.writeCount >= 100

	if forceFlush || softFlush {
		s.writeCount = 0
		s.lastFlush = time.Now()
		return true
	}
	s.writeCount++
	return false
}

func (s *SortedBatchStrategy) Name() string { return "sorted_batch" }

type UseCase string

const (
	UseCaseMonitoring UseCase = "monitoring"
	UseCaseDefrag     UseCase = "defrag"
	UseCaseDownload   UseCase = "download"
	UseCaseMigration  UseCase = "migration"
	UseCaseBalanced   UseCase = "balanced"
)

func GetPresetStrategy(useCase UseCase) FlushStrategy {
	switch useCase {
	case UseCaseMonitoring:
		return NewTimeBasedStrategy(30 * time.Second)
	case UseCaseDefrag:
		return NewSortedBatchStrategy(64*1024*1024, 10*time.Second)
	case UseCaseDownload:
		return NewTimeBasedStrategy(5 * time.Second)
	case UseCaseMigration:
		return NewSortedBatchStrategy(256*1024*1024, 2*time.Second)
	case UseCaseBalanced:
		return NewSortedBatchStrategy(64*1024*1024, 5*time.Second)
	default:
		return NewSortedBatchStrategy(64*1024*1024, 5*time.Second)
	}
}

type WriteRequest struct {
	FilePath string
	Data     []byte
	Offset   int64
	Time     time.Time
}

type Buffer struct {
	Data     []byte
	Offset   int64
	FilePath string
	Modified time.Time
}

type Pool struct {
	buffers      map[string]*Buffer
	mu           sync.RWMutex
	MemoryLimit  int64
	CurrentMem   int64
	strategy     FlushStrategy
	writeQueue   chan WriteRequest
	sortedWrites []WriteRequest
	workerDone   chan struct{}
	closed       bool
}

func NewPool(limit int64, strategy FlushStrategy) *Pool {
	if strategy == nil {
		strategy = NewSortedBatchStrategy(64*1024*1024, 5*time.Second)
	}

	pool := &Pool{
		buffers:     make(map[string]*Buffer),
		MemoryLimit: limit,
		strategy:    strategy,
		writeQueue:  make(chan WriteRequest, 10000),
		workerDone:  make(chan struct{}),
	}
	go pool.worker()
	return pool
}

func (p *Pool) worker() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer close(p.workerDone)

	for {
		select {
		case req := <-p.writeQueue:
			p.queueWrite(req)
		case <-ticker.C:
			if p.strategy.ShouldFlush(p) && !p.closed {
				p.flushSorted()
			}
		}
	}
}

func (p *Pool) queueWrite(req WriteRequest) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		return
	}

	buf, exists := p.buffers[req.FilePath]
	if !exists {
		buf = &Buffer{
			FilePath: req.FilePath,
			Modified: time.Now(),
		}
		p.buffers[req.FilePath] = buf
	}

	if req.Offset == buf.Offset+int64(len(buf.Data)) {
		buf.Data = append(buf.Data, req.Data...)
	} else if req.Offset < buf.Offset {
		buf.Data = req.Data
		buf.Offset = req.Offset
	} else {
		frontPad := make([]byte, req.Offset-buf.Offset)
		buf.Data = append(frontPad, buf.Data...)
		buf.Data = append(buf.Data, req.Data...)
	}
	buf.Modified = time.Now()
	p.CurrentMem += int64(len(req.Data))

	p.sortedWrites = append(p.sortedWrites, req)
}

func (p *Pool) Write(filePath string, data []byte, offset int64) error {
	select {
	case p.writeQueue <- WriteRequest{FilePath: filePath, Data: data, Offset: offset, Time: time.Now()}:
		return nil
	default:
		return errors.New("写入队列已满，无法写入")
	}
}

func (p *Pool) flushSorted() {
	p.mu.Lock()

	if len(p.sortedWrites) == 0 {
		p.mu.Unlock()
		return
	}

	sort.Slice(p.sortedWrites, func(i, j int) bool {
		return p.sortedWrites[i].FilePath < p.sortedWrites[j].FilePath
	})

	currentPath := ""
	var currentData []byte
	var currentOffset int64

	flushBatch := func() {
		if currentPath != "" && len(currentData) > 0 {
			p.flushToDisk(currentPath, currentData, currentOffset)
			delete(p.buffers, currentPath)
		}
	}

	for _, req := range p.sortedWrites {
		if req.FilePath != currentPath {
			flushBatch()
			currentPath = req.FilePath
			currentData = req.Data
			currentOffset = req.Offset
		} else {
			currentData = append(currentData, req.Data...)
		}
	}
	flushBatch()

	p.sortedWrites = p.sortedWrites[:0]
	p.CurrentMem = 0

	p.mu.Unlock()
}

func (p *Pool) flushToDisk(filePath string, data []byte, offset int64) error {
	if len(data) == 0 {
		return nil
	}

	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		if os.IsPermission(err) {
			return errors.New("存储失败，权限不足，请检查文件路径权限")
		}
		if os.IsNotExist(err) {
			return errors.New("存储失败，父目录不存在，请确保目标路径有效")
		}
		return errors.New("存储失败，无法创建文件: " + err.Error())
	}
	defer f.Close()

	_, err = f.WriteAt(data, offset)
	if err != nil {
		if isDiskFullError(err) {
			return errors.New("磁盘空间不足，请清理磁盘后重试")
		}
		return errors.New("存储失败，写入数据时发生错误: " + err.Error())
	}
	return nil
}

func isDiskFullError(err error) bool {
	if err == nil {
		return false
	}
	errMsg := strings.ToLower(err.Error())
	return strings.Contains(errMsg, "no space left") ||
		strings.Contains(errMsg, "enospc") ||
		strings.Contains(errMsg, "disk full") ||
		strings.Contains(errMsg, "quota exceeded")
}

func (p *Pool) Read(filePath string, offset int64, size int) ([]byte, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	buf, exists := p.buffers[filePath]
	if !exists {
		return nil, nil
	}

	start := offset - buf.Offset
	if start < 0 || start >= int64(len(buf.Data)) {
		return nil, nil
	}

	end := start + int64(size)
	if end > int64(len(buf.Data)) {
		end = int64(len(buf.Data))
	}

	result := make([]byte, end-start)
	copy(result, buf.Data[start:end])
	return result, nil
}

func (p *Pool) Flush(filePath string) error {
	p.mu.Lock()
	buf, exists := p.buffers[filePath]
	if !exists {
		p.mu.Unlock()
		return nil
	}
	delete(p.buffers, filePath)
	p.CurrentMem -= int64(len(buf.Data))
	data := buf.Data
	offset := buf.Offset
	p.mu.Unlock()

	return p.flushToDisk(filePath, data, offset)
}

func (p *Pool) FlushAll() error {
	p.mu.Lock()
	paths := make([]string, 0, len(p.buffers))
	for path := range p.buffers {
		paths = append(paths, path)
	}
	p.mu.Unlock()

	var lastErr error
	for _, path := range paths {
		if err := p.Flush(path); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

func (p *Pool) MemoryUsage() int64 {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.CurrentMem
}

func (p *Pool) PendingCount() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.sortedWrites)
}

func (p *Pool) SetStrategy(strategy FlushStrategy) {
	p.mu.Lock()
	p.strategy = strategy
	p.mu.Unlock()
}

func (p *Pool) Close() error {
	p.mu.Lock()
	p.closed = true
	p.mu.Unlock()
	return p.FlushAll()
}

package buffer

import (
	"io"
	"io/fs"
	"os"
	"path"
	"time"

	"github.com/Weidows/wutils/utils/bufioutil"
	"github.com/Weidows/wutils/utils/cache"
	"github.com/binzume/dkango"
)

type BufferConfig struct {
	SourcePath        string
	MemoryLimit       int64
	FlushInterval     int64
	Strategy          string
	EnableReadCache   bool
	EnableWriteBuffer bool
}

type BufferFS struct {
	config      *BufferConfig
	writeBuffer *bufioutil.Pool
	readCache   *cache.LRU
}

func NewBufferFS(config *BufferConfig) *BufferFS {
	fs := &BufferFS{config: config}
	if config.EnableWriteBuffer {
		strategy := bufioutil.GetPresetStrategy(bufioutil.UseCase(config.Strategy))
		if config.FlushInterval > 0 {
			strategy = bufioutil.NewTimeBasedStrategy(time.Duration(config.FlushInterval) * time.Second)
		}
		fs.writeBuffer = bufioutil.NewPool(config.MemoryLimit, strategy)
	}
	if config.EnableReadCache {
		fs.readCache = cache.NewLRU(config.MemoryLimit)
	}
	if config.SourcePath == "" {
		config.SourcePath = "."
	}
	return fs
}

func (b *BufferFS) Open(name string) (fs.File, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrInvalid}
	}
	return os.Open(path.Join(b.config.SourcePath, name))
}

func (b *BufferFS) Stat(name string) (fs.FileInfo, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{Op: "stat", Path: name, Err: fs.ErrInvalid}
	}
	return os.Stat(path.Join(b.config.SourcePath, name))
}

func (b *BufferFS) OpenWriter(name string, flag int) (io.WriteCloser, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrInvalid}
	}
	return os.OpenFile(path.Join(b.config.SourcePath, name), flag, fs.ModePerm)
}

func (b *BufferFS) Remove(name string) error {
	if !fs.ValidPath(name) {
		return &fs.PathError{Op: "remove", Path: name, Err: fs.ErrInvalid}
	}
	return os.Remove(path.Join(b.config.SourcePath, name))
}

func (b *BufferFS) Mkdir(name string, mode fs.FileMode) error {
	if !fs.ValidPath(name) {
		return &fs.PathError{Op: "mkdir", Path: name, Err: fs.ErrInvalid}
	}
	return os.Mkdir(path.Join(b.config.SourcePath, name), mode)
}

func (b *BufferFS) Rename(name, newName string) error {
	if !fs.ValidPath(name) || !fs.ValidPath(newName) {
		return &fs.PathError{Op: "rename", Path: name, Err: fs.ErrInvalid}
	}
	return os.Rename(path.Join(b.config.SourcePath, name), path.Join(b.config.SourcePath, newName))
}

func (b *BufferFS) Truncate(name string, size int64) error {
	if !fs.ValidPath(name) {
		return &fs.PathError{Op: "truncate", Path: name, Err: fs.ErrInvalid}
	}
	return os.Truncate(path.Join(b.config.SourcePath, name), size)
}

var mountInstance interface{ Close() error }

func Mount(drive string, config *BufferConfig) error {
	fs := NewBufferFS(config)
	mount, err := dkango.MountFS(drive, fs, nil)
	if err != nil {
		return err
	}
	mountInstance = mount
	return nil
}

func Unmount() error {
	if mountInstance != nil {
		return mountInstance.Close()
	}
	return nil
}

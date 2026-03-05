package buffer

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"time"

	"github.com/Weidows/wutils/utils/bufioutil"
	"github.com/Weidows/wutils/utils/cache"
	"github.com/binzume/dkango"
)

// 错误定义
var (
	ErrSourcePathNotExist = errors.New("源路径不存在或无法访问")
	ErrSourcePathNotDir   = errors.New("源路径不是有效的目录")
	ErrMountPointInvalid  = errors.New("挂载点无效")
	ErrMountFailed        = errors.New("挂载失败")
	ErrAlreadyMounted     = errors.New("已经存在挂载实例，请先卸载")
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

func NewBufferFS(config *BufferConfig) (*BufferFS, error) {
	if config.SourcePath == "" {
		config.SourcePath = "."
	}
	if err := validateSourcePath(config.SourcePath); err != nil {
		return nil, err
	}

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
	return fs, nil
}

func validateSourcePath(sourcePath string) error {
	info, err := os.Stat(sourcePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%w: %s", ErrSourcePathNotExist, sourcePath)
		}
		return fmt.Errorf("%w: %s", ErrSourcePathNotExist, err.Error())
	}
	if !info.IsDir() {
		return fmt.Errorf("%w: %s", ErrSourcePathNotDir, sourcePath)
	}
	return nil
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
	f, err := os.OpenFile(path.Join(b.config.SourcePath, name), flag, 0644)
	if err != nil {
		if os.IsPermission(err) {
			return nil, fmt.Errorf("文件写入失败，权限不足: %s, 错误: %w", name, err)
		}
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("文件写入失败，父目录不存在: %s, 错误: %w", name, err)
		}
		return nil, fmt.Errorf("文件写入失败: %s, 错误: %w", name, err)
	}
	return f, nil
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
	if drive == "" {
		return fmt.Errorf("%w: 挂载点不能为空", ErrMountPointInvalid)
	}
	if mountInstance != nil {
		return fmt.Errorf("%w", ErrAlreadyMounted)
	}

	fs, err := NewBufferFS(config)
	if err != nil {
		return fmt.Errorf("初始化BufferFS失败: %w", err)
	}

	mount, err := dkango.MountFS(drive, fs, nil)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrMountFailed, err)
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

package service

import (
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/Weidows/wutils/internal/app"
	"github.com/Weidows/wutils/internal/i18n"
	"github.com/Weidows/wutils/utils/files"
	"github.com/rwcarlsen/goexif/exif"
)

// Supported media file extensions
var (
	imageExts = []string{".jpg", ".jpeg", ".png"}
	videoExts = []string{".mp4", ".mov"}
)

const (
	limitTime     = 12.0 // hours
	limitDistance = 1.0  // km
)

// MediaFile represents a photo or video with metadata.
type MediaFile struct {
	Path      string
	Lat       float64
	Lng       float64
	Timestamp time.Time
}

// Cluster represents a group of related media files.
type Cluster []string

// MediaService groups photos and videos by time and GPS proximity.
type MediaService struct{}

// NewMediaService creates a new MediaService.
func NewMediaService() *MediaService {
	return &MediaService{}
}

func (s *MediaService) Name() string               { return "media" }
func (s *MediaService) Description() string        { return i18n.G("media.description") }
func (s *MediaService) Status() app.ServiceStatus  { return app.StatusStopped }
func (s *MediaService) Start() error               { return nil }
func (s *MediaService) Stop() error                { return nil }

// ClusterAndCopy clusters media files in inputDir and copies them to output/.
func (s *MediaService) ClusterAndCopy(inputDir string) {
	outputDir := filepath.Join(inputDir, "output")
	clusters, err := clusterMediaFiles(inputDir)
	if err != nil {
		log.Fatalf("聚类失败: %v", err)
	}
	for i, cluster := range clusters {
		fmt.Printf("聚类 %d: %v\n", i+1, cluster)
	}
	if err := copyFilesToClusters(clusters, outputDir); err != nil {
		log.Fatalf("复制文件失败: %v", err)
	}
	fmt.Println("文件已成功复制到对应目录")
}

func clusterMediaFiles(dirPath string) ([]Cluster, error) {
	mediaFiles, err := getMediaFiles(dirPath)
	if err != nil {
		return nil, fmt.Errorf("获取媒体文件失败: %v", err)
	}
	if len(mediaFiles) == 0 {
		return nil, fmt.Errorf("未找到有效的媒体文件")
	}

	sort.Slice(mediaFiles, func(i, j int) bool {
		return mediaFiles[i].Timestamp.Before(mediaFiles[j].Timestamp)
	})

	var clusters []Cluster
	currentCluster := Cluster{mediaFiles[0].Path}
	currentLat := mediaFiles[0].Lat
	currentLng := mediaFiles[0].Lng
	currentTime := mediaFiles[0].Timestamp

	for i := 1; i < len(mediaFiles); i++ {
		f := mediaFiles[i]
		dist := calcDistance(currentLat, currentLng, f.Lat, f.Lng)
		tdiff := f.Timestamp.Sub(currentTime).Hours()
		if dist > limitDistance || tdiff > limitTime {
			clusters = append(clusters, currentCluster)
			currentCluster = Cluster{f.Path}
			currentLat, currentLng, currentTime = f.Lat, f.Lng, f.Timestamp
		} else {
			currentCluster = append(currentCluster, f.Path)
			currentLat, currentLng, currentTime = f.Lat, f.Lng, f.Timestamp
		}
	}
	if len(currentCluster) > 0 {
		clusters = append(clusters, currentCluster)
	}
	return clusters, nil
}

func calcDistance(lat1, lng1, lat2, lng2 float64) float64 {
	latDiff := (lat1 - lat2) * 111.0
	lngDiff := (lng1 - lng2) * 111.0 * math.Cos(lat1*math.Pi/180.0)
	return math.Sqrt(latDiff*latDiff + lngDiff*lngDiff)
}

func getMediaFiles(dirPath string) ([]MediaFile, error) {
	var result []MediaFile
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		if contains(imageExts, ext) || contains(videoExts, ext) {
			mf, err := extractMetadata(path)
			if err != nil {
				fmt.Printf("警告: 无法从 %s 提取元数据: %v\n", path, err)
				return nil
			}
			result = append(result, mf)
		}
		return nil
	})
	return result, err
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func extractMetadata(filePath string) (MediaFile, error) {
	ext := strings.ToLower(filepath.Ext(filePath))
	var lat, lng float64
	var ts time.Time
	var err error

	switch {
	case contains(imageExts, ext):
		if ext == ".jpg" || ext == ".jpeg" {
			lat, lng, ts, err = extractJPEG(filePath)
		} else {
			// PNG: no EXIF geo data, use file modtime
			fi, e := os.Stat(filePath)
			if e != nil {
				return MediaFile{}, e
			}
			ts = fi.ModTime()
		}
	case contains(videoExts, ext):
		fi, e := os.Stat(filePath)
		if e != nil {
			return MediaFile{}, e
		}
		ts = fi.ModTime()
	default:
		return MediaFile{}, fmt.Errorf("不支持的文件类型: %s", ext)
	}
	if err != nil {
		return MediaFile{}, err
	}
	return MediaFile{Path: filePath, Lat: lat, Lng: lng, Timestamp: ts}, nil
}

func extractJPEG(filePath string) (float64, float64, time.Time, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return 0, 0, time.Time{}, err
	}
	defer f.Close()

	x, err := exif.Decode(f)
	if err != nil {
		fi, e := os.Stat(filePath)
		if e != nil {
			return 0, 0, time.Time{}, e
		}
		return 0, 0, fi.ModTime(), nil
	}

	lat, lng, err := x.LatLong()
	if err != nil {
		lat, lng = 0, 0
	}
	timeStr, err := x.DateTime()
	if err != nil {
		fi, e := os.Stat(filePath)
		if e != nil {
			return 0, 0, time.Time{}, e
		}
		return lat, lng, fi.ModTime(), nil
	}
	return lat, lng, timeStr, nil
}

func copyFilesToClusters(clusters []Cluster, outputDir string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}
	unclassifiedDir := filepath.Join(outputDir, "Unclassified")
	if err := os.MkdirAll(unclassifiedDir, 0755); err != nil {
		return err
	}

	for _, cluster := range clusters {
		firstPath := cluster[0]
		mf, err := extractMetadata(firstPath)
		if err != nil || (mf.Timestamp.IsZero() && mf.Lat == 0 && mf.Lng == 0) {
			for _, fp := range cluster {
				dest := filepath.Join(unclassifiedDir, filepath.Base(fp))
				if err := files.CopyFile(fp, dest); err != nil {
					fmt.Printf("警告: 无法复制文件 %s: %v\n", fp, err)
				}
			}
			continue
		}
		clusterDir := filepath.Join(outputDir, fmt.Sprintf("%s-(%.0f,%.0f)",
			mf.Timestamp.Format("06.1.2-15"), mf.Lat, mf.Lng))
		if err := os.MkdirAll(clusterDir, 0755); err != nil {
			return err
		}
		for _, fp := range cluster {
			dest := filepath.Join(clusterDir, filepath.Base(fp))
			if err := files.CopyFile(fp, dest); err != nil {
				fmt.Printf("警告: 无法复制文件 %s: %v\n", fp, err)
			}
		}
	}
	return nil
}

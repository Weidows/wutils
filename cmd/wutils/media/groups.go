package media

import (
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

// MediaFile 表示一个媒体文件（照片或视频）
type MediaFile struct {
	Path      string
	Lat       float64
	Lng       float64
	Timestamp time.Time
}

// Cluster 表示一个聚类
type Cluster []string

func ClusterAndCopy(inputDir string) {
	outputDir := filepath.Join(inputDir, "output")

	// 执行聚类
	clusters, err := clusterMediaFiles(inputDir)
	if err != nil {
		log.Fatalf("聚类失败: %v", err)
	}

	// 打印聚类结果
	for i, cluster := range clusters {
		fmt.Printf("聚类 %d: %v\n", i+1, cluster)
	}

	// copy files
	err = copyFilesToClusters(clusters, outputDir)
	if err != nil {
		log.Fatalf("复制文件失败: %v", err)
	}
	fmt.Println("文件已成功复制到对应目录")
}

// clusterMediaFiles 根据地理位置和时间戳对媒体文件进行聚类
func clusterMediaFiles(dirPath string) ([]Cluster, error) {
	// 获取所有媒体文件
	mediaFiles, err := getMediaFiles(dirPath)
	if err != nil {
		return nil, fmt.Errorf("获取媒体文件失败: %v", err)
	}

	if len(mediaFiles) == 0 {
		return nil, fmt.Errorf("未找到有效的媒体文件")
	}

	// 按时间戳排序
	sort.Slice(mediaFiles, func(i, j int) bool {
		return mediaFiles[i].Timestamp.Before(mediaFiles[j].Timestamp)
	})

	// 聚类
	clusters := []Cluster{}
	currentCluster := Cluster{mediaFiles[0].Path}
	currentLat := mediaFiles[0].Lat
	currentLng := mediaFiles[0].Lng
	currentTime := mediaFiles[0].Timestamp

	for i := 1; i < len(mediaFiles); i++ {
		file := mediaFiles[i]

		// 计算地理距离（简化版，使用欧几里得距离）
		distance := calculateDistance(currentLat, currentLng, file.Lat, file.Lng)

		// 计算时间差（小时）
		timeDiff := file.Timestamp.Sub(currentTime).Hours()

		// 如果距离大于200米或时间差大于12小时，创建新聚类
		if distance > 0.2 || timeDiff > 12 {
			clusters = append(clusters, currentCluster)
			currentCluster = Cluster{file.Path}
			currentLat = file.Lat
			currentLng = file.Lng
			currentTime = file.Timestamp
		} else {
			currentCluster = append(currentCluster, file.Path)
			// 更新当前聚类的参考点（使用最新的文件）
			currentLat = file.Lat
			currentLng = file.Lng
			currentTime = file.Timestamp
		}
	}

	// 添加最后一个聚类
	if len(currentCluster) > 0 {
		clusters = append(clusters, currentCluster)
	}

	return clusters, nil
}

// calculateDistance 计算两点之间的距离（千米）
func calculateDistance(lat1, lng1, lat2, lng2 float64) float64 {
	// 简化版的距离计算，使用欧几里得距离
	// 在小范围内这种简化是可接受的
	// 1度纬度约等于111千米
	latDiff := (lat1 - lat2) * 111.0
	// 1度经度在赤道约等于111千米，但随着纬度增加而减少
	lngDiff := (lng1 - lng2) * 111.0 * math.Cos(lat1*math.Pi/180.0)

	return math.Sqrt(latDiff*latDiff + lngDiff*lngDiff)
}

// getMediaFiles 获取目录中所有媒体文件的元数据
func getMediaFiles(dirPath string) ([]MediaFile, error) {
	var mediaFiles []MediaFile

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录
		if info.IsDir() {
			return nil
		}

		// 检查文件扩展名
		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".mp4" || ext == ".mov" {
			// 提取元数据
			mediaFile, err := extractMetadata(path)
			if err != nil {
				fmt.Printf("警告: 无法从 %s 提取元数据: %v\n", path, err)
				return nil
			}

			mediaFiles = append(mediaFiles, mediaFile)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return mediaFiles, nil
}

// extractMetadata 从媒体文件中提取元数据
func extractMetadata(filePath string) (MediaFile, error) {
	ext := strings.ToLower(filepath.Ext(filePath))

	var lat, lng float64
	var timestamp time.Time
	var err error

	if ext == ".jpg" || ext == ".jpeg" {
		lat, lng, timestamp, err = extractJpegMetadata(filePath)
	} else if ext == ".png" {
		lat, lng, timestamp, err = extractPngMetadata(filePath)
	} else if ext == ".mp4" || ext == ".mov" {
		// 对于视频文件，我们可能需要使用其他库
		// 这里简化处理，使用文件修改时间
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			return MediaFile{}, err
		}
		timestamp = fileInfo.ModTime()
		// 视频文件可能没有地理信息，使用默认值
		lat, lng = 0, 0
		err = nil
	} else {
		return MediaFile{}, fmt.Errorf("不支持的文件类型: %s", ext)
	}

	if err != nil {
		return MediaFile{}, err
	}

	return MediaFile{
		Path:      filePath,
		Lat:       lat,
		Lng:       lng,
		Timestamp: timestamp,
	}, nil
}

// extractJpegMetadata 从JPEG文件中提取元数据
func extractJpegMetadata(filePath string) (float64, float64, time.Time, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return 0, 0, time.Time{}, err
	}
	defer f.Close()

	// 使用rwcarlsen/goexif库解析EXIF数据
	x, err := exif.Decode(f)
	if err != nil {
		// 如果没有EXIF数据，使用文件修改时间
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			return 0, 0, time.Time{}, err
		}
		return 0, 0, fileInfo.ModTime(), nil
	}

	// 提取地理位置
	lat, lng, err := x.LatLong()
	if err != nil {
		// 如果没有地理信息，使用默认值
		lat, lng = 0, 0
	}

	// 提取时间戳
	timeStr, err := x.DateTime()
	if err != nil {
		// 如果没有时间信息，使用文件修改时间
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			return 0, 0, time.Time{}, err
		}
		return lat, lng, fileInfo.ModTime(), nil
	}

	return lat, lng, timeStr, nil
}

// extractPngMetadata 从PNG文件中提取元数据
func extractPngMetadata(filePath string) (float64, float64, time.Time, error) {
	// PNG文件可能没有地理信息，使用默认值
	lat, lng := 0.0, 0.0

	// 使用文件修改时间作为时间戳
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0, 0, time.Time{}, err
	}
	timestamp := fileInfo.ModTime()

	return lat, lng, timestamp, nil
}

// copyFilesToClusters 将文件复制到对应的聚类目录
func copyFilesToClusters(clusters []Cluster, outputDir string) error {
	// 确保输出目录存在
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		return err
	}

	// 为每个聚类创建目录并复制文件
	for i, cluster := range clusters {
		clusterDir := filepath.Join(outputDir, fmt.Sprintf("cluster_%d", i+1))
		err := os.MkdirAll(clusterDir, 0755)
		if err != nil {
			return err
		}

		for _, filePath := range cluster {
			// 获取文件名
			fileName := filepath.Base(filePath)
			// 目标路径
			destPath := filepath.Join(clusterDir, fileName)

			// 复制文件
			err := copyFile(filePath, destPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// copyFile 复制文件
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return nil
}

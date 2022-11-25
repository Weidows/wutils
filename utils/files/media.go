package files

import (
	"bytes"
	"encoding/binary"
	"io"
	"time"
)

// boxHeader 信息头
type boxHeader struct {
	Size       uint32
	FourccType [4]byte
	Size64     uint64
}

// GetMP4Duration 获取视频时长，以秒计
//
// Modified from https://github.com/akkuman/mp4info
func GetMP4Duration(reader io.ReaderAt) (lengthOfTime time.Duration, err error) {
	var (
		info   = make([]byte, 0x10)
		header boxHeader
		offset int64 = 0
	)

	// 获取moov结构偏移
	for {
		_, err = reader.ReadAt(info, offset)
		if err != nil {
			return
		}
		// getHeaderBoxInfo 获取头信息
		if err = binary.Read(bytes.NewBuffer(info), binary.BigEndian, &header); err != nil {
			header = boxHeader{}
		}

		fourccType := string(header.FourccType[:])
		if fourccType == "moov" {
			break
		} else if fourccType == "mdat" {
			// 有一部分mp4 mdat尺寸过大需要特殊处理
			if header.Size == 1 {
				offset += int64(header.Size64)
				continue
			}
		}
		offset += int64(header.Size)
	}
	// 获取moov结构开头一部分
	moovStartBytes := make([]byte, 0x100)
	_, err = reader.ReadAt(moovStartBytes, offset)
	if err != nil {
		return
	}
	// 定义timeScale与Duration偏移
	timeScaleOffset := 0x1C
	durationOffest := 0x20
	timeScale := binary.BigEndian.Uint32(moovStartBytes[timeScaleOffset : timeScaleOffset+4])
	Duration := binary.BigEndian.Uint32(moovStartBytes[durationOffest : durationOffest+4])
	return time.Second * time.Duration(Duration/timeScale), nil
}

// https://lexica.art/docs
// 爬取信息

//func GetFileFromURL(url, path string) error {
//	resp, err := http.Get(url)
//	if err != nil {
//		return err
//	}
//	defer func(Body io.ReadCloser) {
//		_ = Body.Close()
//	}(resp.Body)
//
//	if err = os.MkdirAll(path, 0750); err != nil {
//		return err
//	}
//
//	file, err := os.Open(path)
//	if err != nil {
//		return err
//	}
//	io.CopyBuffer()
//}

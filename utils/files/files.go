package files

//func getVideos(path string, file string) error {
//	fullPath := filepath.Join(path, file)
//	dstF, err := os.Open(fullPath)
//	if err != nil {
//		return err
//	}
//	defer dstF.Close()
//	fileInfo, err := dstF.Stat()
//	if err != nil {
//		return err
//	}
//
//	if fileInfo.IsDir() {
//		fileList, err := dstF.Readdir(0) //获取文件夹下各个文件或文件夹的fileInfo
//		if err != nil {
//			return err
//		}
//		for _, fileInfo = range fileList {
//			if err = getVideos(fullPath, fileInfo.Name()); err != nil {
//				return err
//			}
//		}
//		return nil
//	} else {
//		if strings.Contains(file, "video.mp4") {
//			videos = append(videos, map[string]string{
//				"path": path,
//				"file": file,
//			})
//		}
//		return nil
//	}
//}

func WithOpen() {

}

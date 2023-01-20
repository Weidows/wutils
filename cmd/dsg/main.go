package main

import (
	"fmt"
	"github.com/Weidows/Golang/utils/log"
	"github.com/cheggaaa/pb/v3"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	logger   = log.GetLogger()
	filePath string
	count    = 30
)

func init() {
	args := os.Args
	if len(args) == 1 {
		fmt.Println(
			"please start with params like: 'dsg.exe E: 30'\n",
			"\t1. disk\n",
			"\t2. delay seconds",
		)
		return
	}

	f := strings.Join([]string{args[1], ".dsg"}, "/")
	// RemoveAll 不会因为文件不存在 return error
	if err := os.RemoveAll(f); err != nil {
		logger.Printf("file: '%s' state error, maybe use by other processes", f)
		return
	}
	filePath = f

	if len(args) > 2 {
		c, err := strconv.Atoi(args[2])
		if err != nil {
			logger.Println("delay seconds number is not valid.")
			return
		}
		count = c
	}
}

func main() {
	bar := pb.Simple.Start(count)
	for {
		for i := 0; i < count; i++ {
			bar.Increment()
			time.Sleep(time.Second)
		}
		go WriteString()
		bar.SetCurrent(0)
	}
}

func WriteString() {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		logger.Println("disk format error, please input like 'E:'", err)
		return
	}
	_, _ = file.WriteString("dsg running at " + time.Now().String() + "\n")
	_ = file.Close()
}

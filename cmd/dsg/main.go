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
	logger = log.GetLogger()
	count  = 30
)

func main() {
	file := parseArgs(os.Args)

	bar := pb.Simple.Start(count)
	for file != nil {
		for i := 0; i < count; i++ {
			bar.Increment()
			time.Sleep(time.Second)
		}
		_, _ = file.WriteString("dsg running at " + time.Now().String() + "\n")
		bar.SetCurrent(0)
	}
}

func parseArgs(args []string) *os.File {
	if len(args) == 1 {
		fmt.Println(
			"please start with params like: 'dsg.exe E: 30'\n",
			"\t1. disk\n",
			"\t2. delay seconds",
		)
		return nil
	}

	f := strings.Join([]string{args[1], ".dsg"}, "/")
	if err := os.Remove(f); err != nil {
		logger.Printf("file: '%s' state error, maybe not exist or use by other processes", f)
		return nil
	}
	file, err := os.OpenFile(f, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		logger.Println("disk format error, please input like 'E:'", err)
		return nil
	}

	if len(args) > 2 {
		c, err := strconv.Atoi(args[2])
		if err != nil {
			logger.Println("delay seconds number is not valid.")
			return nil
		}
		count = c
	}
	return file
}

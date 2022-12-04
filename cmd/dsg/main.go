package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cheggaaa/pb/v3"
)

func main() {
	file, count := parseArgs(os.Args)

	bar := pb.Simple.Start(count)
	for true {
		for i := 0; i < count; i++ {
			bar.Increment()
			time.Sleep(time.Second)
		}
		_, _ = file.WriteString("dsg running at " + time.Now().String() + "\n")
		bar.SetCurrent(0)
	}
}

func parseArgs(args []string) (*os.File, int) {
	if len(args) == 1 {
		fmt.Println(`
			please start with params like: 
			dsg.exe E: 30
			1. disk
			2. delay seconds
		`)
		return nil, 0
	}

	f := strings.Join([]string{args[1], ".dsg"}, "/")
	file, err := os.OpenFile(f, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println(`
			disk format error, please input like 'E:'
		`, err)
		return nil, 0
	}

	count := 30
	if len(args) > 2 {
		c, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Println(`
				delay seconds number is not valid.
			`)
			return nil, 0
		}
		count = c
	}
	return file, count
}

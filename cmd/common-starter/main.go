package main

import (
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		return
	}
	script := os.Args[1]
	println("script: ", script)
	var err error

	// 根据文件后缀来决定如何执行脚本
	switch filepath.Ext(script) {
	case ".bat", ".exe":
		// 在新的窗口中执行脚本
		err = exec.Command("cmd", "/C", "start", script).Start()
	case ".sh":
		err = exec.Command("cmd", "/C", "start", "bash", "-c", script).Start()
	case ".ps1":
		err = exec.Command("powershell", script).Start()
	case ".cmd":
		err = exec.Command("cmd", "/c", script).Start()
	default:
		println("Unknown file extension")
	}

	if err != nil {
		println(err)
	}
}

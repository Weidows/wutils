/*
# jpu

# Jetbrains Portable Upgrader

	go install github.com/Weidows/wutils/cmd/jpu

	echo "Default to use $SCOOP/persist/jetbrains-toolbox/apps"
	set JPU_PATH=D:/Scoop/persist/jetbrains-toolbox/apps && jpu.exe

***

通过改配置实现 Portable 效果

	```
	- PyCharm-P
	  - ch-0
	  - 223.8214.51
	  - bin
	  - idea.properties
	  - 223.8214.51.plugins
	  - 223.8617.48
	  - 223.8617.48.plugins
	  - config
	  - system

	- Goland
	- datagrip
	```

在 `IDE/bin/idea.properties` 顶部添加

	idea.config.path=${idea.home.path}/../../config
	idea.system.path=${idea.home.path}/../../system

由于路径中含带版本号, 用脚本不易操作, 所以用 go 写
*/
package main

import (
	"fmt"
	"github.com/magiconair/properties"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	appPath string
)

func init() {
	appPath = os.Getenv("JPU_PATH")
	if appPath != "" {
		return
	}

	// 默认取 scoop 位置
	scoopPath := os.Getenv("SCOOP")
	appPath = filepath.Join(scoopPath, "persist/jetbrains-toolbox/apps")
}

func main() {
	fmt.Println("Using app path: ", appPath)
	ides, err := os.ReadDir(appPath)
	if err != nil {
		return
	}
	for _, ide := range ides {
		versions, _ := os.ReadDir(filepath.Join(appPath, ide.Name(), "ch-0"))
		maxVersion := ""
		for _, version := range versions {
			if !version.IsDir() || filepath.Ext(version.Name()) == ".plugins" {
				continue
			}
			if version.Name() > maxVersion {
				maxVersion = version.Name()
			}
		}

		switch ide.Name() {
		case "Fleet":
			continue
		}
		modifyProperties(ide.Name(), maxVersion)
	}
}

func modifyProperties(ideName, maxVersion string) {
	property := filepath.Join(appPath, ideName, "ch-0", maxVersion, "bin/idea.properties")
	p := properties.NewProperties()
	p.Set("idea.config.path", "${idea.home.path}/../../config")
	p.Set("idea.system.path", "${idea.home.path}/../../system")
	f, _ := os.OpenFile(property, os.O_RDWR, 0644)
	p.Write(f, properties.UTF8)
	f.Close()
}

// Deprecated
func linkCurrentVersion(ideName, maxVersion string) {
	os.Chdir(filepath.Join(appPath, ideName, "ch-0"))
	var exe string
	switch ideName {
	case "Aqua":
		exe = "aqua64.exe"
	case "datagrip":
		exe = "datagrip64.exe"
	case "Goland":
		exe = "goland64.exe"
	case "IDEA-U":
		exe = "idea64.exe"
	case "PyCharm-P":
		exe = "pycharm64.exe"
	}
	os.Remove(exe)
	exePath := strings.ReplaceAll(filepath.Join(maxVersion, "bin", exe), "/", "\\")
	_ = exec.Command("cmd", "/c", "mklink", "/h", exe, exePath).Start()
}

package main

import (
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
	args := os.Args
	if len(args) > 0 {
		appPath = args[0]
		return
	}

	// 默认取 scoop 位置
	scoopPath := os.Getenv("SCOOP")
	appPath = filepath.Join(scoopPath, "persist/jetbrains-toolbox/apps")
}

func main() {
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

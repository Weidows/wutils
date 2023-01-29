package main

import (
	"github.com/magiconair/properties"
	"os"
	"path/filepath"
)

var (
	appPath string
)

// ide 中不会执行, cmd 中会
func init() {
	args := os.Args
	if len(args) > 0 {
		appPath = args[0]
	}
}

func initialize() {
	if appPath != "" {
		return
	}
	// 默认取 scoop 位置
	scoopPath := os.Getenv("SCOOP")
	appPath = filepath.Join(scoopPath, "persist/jetbrains-toolbox/apps")
}

func main() {
	initialize()
	ides, err := os.ReadDir(appPath)
	if err != nil {
		return
	}
	for _, ide := range ides {
		switch ide.Name() {
		case "Fleet":
			continue
		}

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
		property := filepath.Join(appPath, ide.Name(), "ch-0", maxVersion, "bin/idea.properties")
		p := properties.NewProperties()
		p.Set("idea.config.path", "${idea.home.path}/../../config")
		p.Set("idea.system.path", "${idea.home.path}/../../system")
		f, _ := os.OpenFile(property, os.O_RDWR, 0644)
		p.Write(f, properties.UTF8)
		f.Close()
	}
}

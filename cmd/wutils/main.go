package main

import (
	"os"

	"github.com/Weidows/wutils/internal/cli"
	"github.com/Weidows/wutils/utils/log"
)

func main() {
	app := cli.NewApp()
	if err := app.Run(os.Args); err != nil {
		log.GetLogger().Fatal(err.Error())
	}
}

package command

import (
	"fmt"
	"github.com/urfave/cli"
	"log"
	"os"
	"strings"
)

func NewApp() {
	app := cli.NewApp()
	app.Name = "pineapple-backups"
	app.Version = "V.1.9.7"
	app.Usage = "https://github.com/AlexiaVeronica/pineapple-backups"
	app.Flags = commandArgs

	app.Action = func(c *cli.Context) {
		if !strings.Contains(Command.AppType, "cat") && !strings.Contains(Command.AppType, "sfacg") {
			log.Fatalf("%s app type error", Command.AppType)
		}
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Println(app.HideHelp)
		log.Fatal(err)
	}
}

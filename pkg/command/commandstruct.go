package command

import (
	"fmt"
	"github.com/urfave/cli"
)

var Command = &commandLines{
	AppType:   "sfacg",
	MaxThread: 32,
}

type commandLines struct {
	BookID    string
	Account   string
	Password  string
	AppType   string
	SearchKey string
	MaxThread int
	Token     bool
	Login     bool
	ShowInfo  bool
	Update    bool
	Epub      bool
}

const (
	FlagAppType   = "appType"
	FlagDownload  = "download"
	FlagToken     = "token"
	FlagMaxThread = "maxThread"
	FlagUser      = "user"
	FlagPassword  = "password"
	FlagUpdate    = "update"
	FlagSearch    = "search"
	FlagLogin     = "login"
	FlagEpub      = "epub"
)

var commandArgs = []cli.Flag{
	cli.StringFlag{
		Name:        fmt.Sprintf("a, %s", FlagAppType),
		Value:       "sfacg",
		Usage:       "cheng app type",
		Destination: &Command.AppType,
	},
	cli.StringFlag{
		Name:        fmt.Sprintf("d, %s", FlagDownload),
		Value:       "",
		Usage:       "book id",
		Destination: &Command.BookID,
	},
	cli.BoolFlag{
		Name:        fmt.Sprintf("t, %s", FlagToken),
		Usage:       "input hbooker token",
		Destination: &Command.Token,
	},
	cli.IntFlag{
		Name:        fmt.Sprintf("m, %s", FlagMaxThread),
		Value:       16,
		Usage:       "change max thread number",
		Destination: &Command.MaxThread,
	},
	cli.StringFlag{
		Name:        fmt.Sprintf("u, %s", FlagUser),
		Value:       "",
		Usage:       "input account name",
		Destination: &Command.Account,
	},
	cli.StringFlag{
		Name:        fmt.Sprintf("p, %s", FlagPassword),
		Value:       "",
		Usage:       "input password",
		Destination: &Command.Password,
	},
	cli.BoolFlag{
		Name:        FlagUpdate,
		Usage:       "update book",
		Destination: &Command.Update,
	},
	cli.StringFlag{
		Name:        fmt.Sprintf("s, %s", FlagSearch),
		Value:       "",
		Usage:       "search book by keyword",
		Destination: &Command.SearchKey,
	},
	cli.BoolFlag{
		Name:        fmt.Sprintf("l, %s", FlagLogin),
		Usage:       "login local account",
		Destination: &Command.Login,
	},
	cli.BoolFlag{
		Name:        fmt.Sprintf("e, %s", FlagEpub),
		Usage:       "start epub",
		Destination: &Command.Epub,
	},
}

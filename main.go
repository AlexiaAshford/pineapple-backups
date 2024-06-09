package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/AlexiaVeronica/input"
	"github.com/AlexiaVeronica/pineapple-backups/config"
	"github.com/AlexiaVeronica/pineapple-backups/pkg/app"
	"github.com/AlexiaVeronica/pineapple-backups/pkg/file"
	"github.com/AlexiaVeronica/pineapple-backups/pkg/tools"
	"github.com/urfave/cli"
)

var (
	apps   *app.APP
	newCli *cli.App
	cmd    = &commandLines{
		AppType:   "sfacg",
		MaxThread: 32,
	}
)

type commandLines struct {
	BookID, Account, Password, AppType, SearchKey   string
	MaxThread                                       int
	Token, Login, ShowInfo, Update, Epub, BookShelf bool
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
	FlagBookShelf = "bookshelf"
)

func init() {
	if !config.Exist("./config.json") || file.SizeFile("./config.json") == 0 {
		fmt.Println("config.json does not exist, creating a new one!")
	} else {
		config.LoadJson()
	}
	config.UpdateConfig()

	newCli = cli.NewApp()
	newCli.Name = "pineapple-backups"
	newCli.Version = "V.2.2.1"
	newCli.Usage = "https://github.com/AlexiaVeronica/pineapple-backups"
	newCli.Flags = defineFlags()
	newCli.Action = validateAppType

	if err := newCli.Run(os.Args); err != nil {
		log.Fatal(err)
	}

	setupTokens()
}

func defineFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{Name: fmt.Sprintf("a, %s", FlagAppType), Value: "sfacg", Usage: "change app type", Destination: &cmd.AppType},
		cli.StringFlag{Name: fmt.Sprintf("d, %s", FlagDownload), Usage: "book id", Destination: &cmd.BookID},
		cli.BoolFlag{Name: fmt.Sprintf("t, %s", FlagToken), Usage: "input hbooker token", Destination: &cmd.Token},
		cli.IntFlag{Name: fmt.Sprintf("m, %s", FlagMaxThread), Value: 16, Usage: "change max thread number", Destination: &cmd.MaxThread},
		cli.StringFlag{Name: fmt.Sprintf("u, %s", FlagUser), Usage: "input account name", Destination: &cmd.Account},
		cli.StringFlag{Name: fmt.Sprintf("p, %s", FlagPassword), Usage: "input password", Destination: &cmd.Password},
		cli.BoolFlag{Name: FlagUpdate, Usage: "update book", Destination: &cmd.Update},
		cli.StringFlag{Name: fmt.Sprintf("s, %s", FlagSearch), Usage: "search book by keyword", Destination: &cmd.SearchKey},
		cli.BoolFlag{Name: fmt.Sprintf("l, %s", FlagLogin), Usage: "login local account", Destination: &cmd.Login},
		cli.BoolFlag{Name: fmt.Sprintf("e, %s", FlagEpub), Usage: "start epub", Destination: &cmd.Epub},
		cli.BoolFlag{Name: fmt.Sprintf("b, %s", FlagBookShelf), Usage: "show bookshelf", Destination: &cmd.BookShelf},
	}
}

func validateAppType(c *cli.Context) {
	if !strings.Contains(cmd.AppType, "cat") && !strings.Contains(cmd.AppType, "sfacg") {
		log.Fatalf("%s app type error", cmd.AppType)
	}
}

func setupTokens() {
	apps.Ciweimao.SetToken(config.Apps.Hbooker.Account, config.Apps.Hbooker.LoginToken)
	apps.Boluobao.Cookie = config.Apps.Sfacg.Cookie
}

func shellSwitch(inputs []string) {
	switch inputs[0] {
	case "up", "update":
		// Update function placeholder
	case "a", "app":
		if tools.TestList([]string{app.BoluobaoLibAPP, app.CiweimaoLibAPP}, inputs[1]) {
			apps.CurrentApp = inputs[1]
		} else {
			fmt.Println("app type error, please input again.")
		}
	case "d", "download":
		if len(inputs) == 2 {
			fmt.Println("download book by book id:", inputs)
			apps.DownloadBookByBookId(inputs[1])
		} else {
			fmt.Println("input book id or url, like: download <bookid/url>")
		}
	case "bs", "bookshelf":
		apps.Bookshelf()
	case "s", "search":
		apps.SearchDetailed(inputs[1])
	case "l", "login":
		handleLogin(inputs)
	case "t", "token":
		apps.Ciweimao.SetToken(inputs[1], inputs[2])
	default:
		fmt.Println("command not found, please input help to see the command list:", inputs[0])
	}
}

func handleLogin(inputs []string) {
	if len(inputs) < 3 {
		fmt.Println("you must input account and password, like: -login account password")
		return
	}
	switch apps.CurrentApp {
	case app.CiweimaoLibAPP:
		apps.Ciweimao.SetToken(inputs[1], inputs[2])
	case app.BoluobaoLibAPP:
		loginStatus, err := apps.Boluobao.API().Login(inputs[1], inputs[2])
		if err == nil {
			apps.Boluobao.Cookie = loginStatus.Cookie
		}
	}
}

func shell(messageOpen bool) {
	if messageOpen {
		for _, message := range config.HelpMessage {
			fmt.Println("[info]", message)
		}
	}
	for {
		if inputRes := input.StringInput(">"); len(inputRes) > 0 {
			shellSwitch(strings.Split(inputRes, " "))
		}
	}
}

func main() {
	if len(os.Args) > 1 {
		handleCommandLine()
	} else {
		shell(true)
	}
}

func handleCommandLine() {
	switch {
	case cmd.Login:
		loginStatus, err := apps.Boluobao.API().Login(cmd.Account, cmd.Password)
		if err == nil {
			apps.Boluobao.Cookie = loginStatus.Cookie
		}
	case cmd.BookID != "":
		apps.DownloadBookByBookId(cmd.BookID)
	case cmd.SearchKey != "":
		apps.SearchDetailed(cmd.SearchKey)
	case cmd.Update:
		// Update function placeholder
	case cmd.Token:
		apps.Ciweimao.SetToken(input.StringInput("Please input account:"), input.StringInput("Please input token:"))
	case cmd.BookShelf:
		apps.Bookshelf()
	default:
		fmt.Println("command not found, please input help to see the command list:", os.Args[1])
	}
}

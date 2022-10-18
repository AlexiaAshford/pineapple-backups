package config

import (
	"gopkg.in/urfave/cli.v1"
)

var CommandLines = struct {
	BookId    string
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
}{}

var Args = []cli.Flag{
	cli.StringFlag{
		Name:        "a, app",
		Value:       "cat",
		Usage:       "cheng app type",
		Destination: &CommandLines.AppType,
	},
	cli.StringFlag{
		Name:        "d, download",
		Value:       "",
		Usage:       "book id",
		Destination: &CommandLines.BookId,
	},
	cli.BoolFlag{
		Name:        "t, token",
		Usage:       "input hbooker token",
		Destination: &CommandLines.Token,
	},
	cli.IntFlag{
		Name:        "m, max",
		Value:       16,
		Usage:       "change max thread number",
		Destination: &CommandLines.MaxThread,
	},
	cli.StringFlag{
		Name:        "u, user",
		Value:       "",
		Usage:       "input account name",
		Destination: &CommandLines.Account,
	},
	cli.StringFlag{
		Name:        "p, password",
		Value:       "",
		Usage:       "input password",
		Destination: &CommandLines.Password,
	},
	cli.BoolFlag{
		Name:        "update",
		Usage:       "update book",
		Destination: &CommandLines.Update,
	},
	cli.StringFlag{
		Name:        "s, search",
		Value:       "",
		Usage:       "show config",
		Destination: &CommandLines.SearchKey,
	},
	cli.BoolFlag{
		Name:        "l, login",
		Usage:       "login local account",
		Destination: &CommandLines.Login,
	},
	cli.BoolFlag{
		Name:        "e, epub",
		Usage:       "start epub",
		Destination: &CommandLines.Epub,
	},
}

//func InitCommand() {
//	app := cli.NewApp()
//	app.Name = "pineapple-backups"
//	app.Version = "V.1.5.2"
//	app.Usage = "https://github.com/VeronicaAlexia/pineapple-backups"
//	app.Flags = Args
//	app.Action = func(c *cli.Context) {
//		fmt.Println("you can input -h and --help to see the command list.")
//		if CommandLines.AppType != "cat" && CommandLines.AppType != "sfacg" {
//			fmt.Println(CommandLines.AppType, "app type error, default app type is cat.")
//			CommandLines.AppType = "cat" // default app type is cat
//		}
//	}
//	if err := app.Run(os.Args); err != nil {
//		log.Fatal(err)
//	}
//
//}

// delete cobra command line
//func CommandInit() []string {
//	rule_cmd := ruleCmd()
//	AddFlags := rule_cmd.Flags()
//	AddFlags.StringVarP(&Book_id, "download", "d", "", "")
//	AddFlags.StringVar(&Account, "account", "", "input account")
//	AddFlags.StringVar(&Password, "password", "", "input password")
//	AddFlags.StringVarP(&Token, "token", "t", "", "input password")
//	AddFlags.StringVarP(&App_type, "app", "a", "sfacg", "input app type")
//	AddFlags.StringVarP(&Search_key, "search", "s", "", "input search keyword")
//	AddFlags.IntVarP(&max_thread, "max", "m", 32, "change max thread number")
//	AddFlags.BoolVar(&show_info, "show", false, "show config")
//	AddFlags.BoolVar(&up_date, "update", false, "update config")
//	if err := rule_cmd.Execute(); err != nil {
//		fmt.Println("ruleCmd error:", err)
//	} else {
//		if tool.TestList([]string{"sfacg", "cat"}, App_type) {
//			Vars.AppType = App_type
//		} else {
//			fmt.Println("app type error, default sfacg")
//		}
//		Vars.ThreadNum = max_thread
//
//		if show_info {
//			tool.FormatJson(ReadConfig(""))
//		}
//		if Book_id != "" {
//			command_line = []string{"download", Book_id}
//		} else if Search_key != "" {
//			command_line = []string{"search", Search_key}
//		} else if up_date {
//			command_line = []string{"update"}
//		} else if Account != "" && Password != "" {
//			command_line = []string{"login", Account, Password}
//		} else {
//			command_line = []string{"console", ""}
//		}
//
//	}
//	return command_line
//}

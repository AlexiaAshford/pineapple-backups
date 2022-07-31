package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"sf/cfg"
	"sf/multi"
	"sf/src"
	"strings"
)

func ShellLoginAccount(account, password string) {
	if account == "" {
		fmt.Println("you must input account, like: sf username")
	} else if password == "" {
		fmt.Println("you must input password, like: sf password")
	} else {
		src.LoginAccount(account, password, 0)
	}
}

func shellBookDownload(downloadId any) {
	switch downloadId.(type) {
	case string:
		start := src.BookInits{BookID: downloadId.(string), Index: 0, Locks: nil, ShowBook: true}
		catalogues := start.DownloadBookInit() // get book catalogues
		if catalogues.TestBookResult {
			catalogues.InitCatalogue()
		}

	case []string:
		Locks := multi.NewGoLimit(7)
		for BookIndex, BookId := range downloadId.([]string) {
			Locks.Add()
			start := src.BookInits{BookID: BookId, Index: BookIndex, Locks: Locks, ShowBook: true}
			catalogues := start.DownloadBookInit() // get book catalogues
			if catalogues.TestBookResult {
				catalogues.InitCatalogue()
			}
		}
		Locks.WaitZero() // wait for all goroutines to finish
	}
	os.Exit(0) // exit the program if no error
}

func ParseCommand() map[string]string {
	AppList := []string{"sfacg", "cat"}
	CommandMap := make(map[string]string)
	download := flag.String("download", "", "input book id or url, like:download <bookid/url>")
	account := flag.String("account", "", "input account")
	password := flag.String("password", "", "input password")
	appType := flag.String("app", "sfacg", "input app type, like: app sfacg")
	search := flag.String("search", "", "input search keyword, like: search keyword")
	Thread := flag.Int("max", 0, "input thread number, like: thread 1")
	showConfig := flag.Bool("show", false, "show config, like: show config")
	flag.Parse()
	CommandMap["book_id"] = cfg.ExtractBookID(*download)
	CommandMap["account"] = *account
	CommandMap["password"] = *password
	CommandMap["app_type"] = *appType
	CommandMap["key_word"] = *search
	if *showConfig {
		// print config json file information to console
		cfg.FormatJson(cfg.ReadConfig(""))
	}
	// check app type and cheng edit config
	if cfg.TestList(AppList, *appType) {
		cfg.Vars.AppType = *appType
	} else {
		fmt.Printf("app type %v is invalid, please input again:", *appType)
		cfg.Vars.AppType = "sfacg"
	}
	if *Thread != 0 {
		if *Thread >= 64 {
			fmt.Println("thread number is too large, please input again:")
		} else {
			cfg.Vars.MaxThreadNumber = *Thread
			fmt.Println("change thread number to:", *Thread, "thread")
		}
	}
	return CommandMap
}

func TestCatAccount() bool {
	if cfg.Vars.AppType == "cat" {
		if cfg.Vars.Cat.Params.Account != "" && cfg.Vars.Cat.Params.LoginToken != "" {
			return false
		}
		for {
			LoginToken := cfg.InputStr("you must input 32 characters login token:")
			if len(LoginToken) != 32 {
				fmt.Println("Login token is 32 characters, please input again:")
			} else {
				cfg.Vars.Cat.Params.LoginToken = LoginToken
				cfg.SaveJson()
				break
			}
		}
		cfg.Vars.Cat.Params.Account = cfg.InputStr("you must input account:")
		cfg.SaveJson()
		return false
	}
	return false
}

func TestSfAccount(account string, password string) bool {
	if cfg.Vars.AppType == "sfacg" {
		if account != "" && password != "" {
			ShellLoginAccount(account, password)
		} else if cfg.Vars.Sfacg.UserName != "" || cfg.Vars.Sfacg.Password != "" {
			if cfg.Vars.Sfacg.Cookie != "" {
				if src.AccountDetailed() == "需要登录才能访问该资源" {
					fmt.Printf("cookie is Invalid,attempt to auto login!\naccount:%v\npassword:%v\n",
						cfg.Vars.Sfacg.UserName, cfg.Vars.Sfacg.Password)
					ShellLoginAccount(cfg.Vars.Sfacg.UserName, cfg.Vars.Sfacg.Password)
				}
			} else {
				fmt.Printf("cookie is empty,attempt to auto login!\naccount:%v\npassword:%v\n",
					cfg.Vars.Sfacg.UserName, cfg.Vars.Sfacg.Password)
				ShellLoginAccount(cfg.Vars.Sfacg.UserName, cfg.Vars.Sfacg.Password)
			}
			return false
		}
	}
	return true
}

func shellConsole(command map[string]string) {
	ExitProgram := false
	if TestSfAccount(command["account"], command["password"]) {
		ExitProgram = true
	}
	if TestCatAccount() {
		ExitProgram = true
	}
	if command["key_word"] != "" {
		if NovelId := src.SearchBook(command["key_word"]); NovelId != "" {
			shellBookDownload(NovelId)
		} else {
			fmt.Println("no book found")
		}
		ExitProgram = true
	}
	if command["book_id"] != "" {
		shellBookDownload(command["book_id"])
		ExitProgram = true
	}

	if ExitProgram {
		os.Exit(0) // exit the program if no error
	}

}
func init() {
	cfg.ConfigInit()
}

func main() {

	if len(os.Args) <= 1 {
		for {
			spaceRe, _ := regexp.Compile(`\s+`)
			inputs := spaceRe.Split(strings.TrimSpace(cfg.Input(">")), -1)
			if len(inputs) > 1 {
				commands := make(map[string]string)
				commands[inputs[0]] = inputs[1]
				shellConsole(commands)
			} else {
				fmt.Println("you must input command, like: sf command")
			}
		}
	} else {

		shellConsole(ParseCommand())
	}
}

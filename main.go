package main

import (
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

func shell(command map[string]any) {
	ExitProgram := false
	if TestSfAccount(command["account"].(string), command["password"].(string)) {
		ExitProgram = true
	}
	if TestCatAccount() {
		ExitProgram = true
	}
	if command["key_word"] != "" {
		if NovelId := src.SearchBook(command["key_word"].(string)); NovelId != "" {
			shellBookDownload(NovelId)
		} else {
			fmt.Println("no book found")
		}
		ExitProgram = true
	}
	if command["book"] != "" {
		shellBookDownload(command["book"].(string))
		ExitProgram = true
	}

	if ExitProgram {
		os.Exit(0) // exit the program if no error
	}

}
func shellConsole(inputs []string) {
	switch inputs[0] {
	case "exit":
		os.Exit(0)
	case "help":
		fmt.Println("help:")
	case "download":
		if len(inputs) == 2 {
			shellBookDownload(inputs[1])
		} else {
			fmt.Println("input book id or url, like:download <bookid/url>")
		}
	case "search":
		if len(inputs) == 2 {
			if NovelId := src.SearchBook(inputs[1]); NovelId != "" {
				shellBookDownload(NovelId)
			} else {
				fmt.Println("no book found")
			}
		} else {
			fmt.Println("input book id or url, like:download <bookid/url>")
		}
	case "show":
		if len(inputs) == 2 {
			if inputs[1] == "config" {
				cfg.FormatJson(cfg.ReadConfig(""))
			}
		} else {
			fmt.Println("input config, like:show config")
		}
	}

	if inputs[0] == "exit" {
		os.Exit(0)
	}

	if len(inputs) >= 2 {

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
				shellConsole(inputs)
			} else {
				fmt.Println("you must input command, like: sf command")
			}
		}
	} else {
		shell(cfg.ParseCommandLine())
	}
}

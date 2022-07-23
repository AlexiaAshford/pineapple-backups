package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"sf/cfg"
	"sf/multi"
	"sf/src"
)

func ExtractBookID(url string) string {
	if url != "" {
		bookID := regexp.MustCompile(`(\d+)`).FindStringSubmatch(url)
		if len(bookID) > 1 {
			return bookID[1]
		}
	}
	return ""
}
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
		if cfg.Vars.AppType == "sfacg" {
			start.SfacgBookInit()
			start.CataloguesInit()
		}
		if cfg.Vars.AppType == "cat" {
			start.CatBookInit()
			start.CataloguesInit()
		}
	case []string:
		Locks := multi.NewGoLimit(7)
		for BookIndex, BookId := range downloadId.([]string) {
			Locks.Add()
			start := src.BookInits{BookID: BookId, Index: BookIndex, Locks: Locks, ShowBook: true}
			if cfg.Vars.AppType == "sfacg" {
				start.SfacgBookInit()
				start.CataloguesInit()
			}
			if cfg.Vars.AppType == "cat" {
				start.CatBookInit()
				start.CataloguesInit()
			}
		}
		Locks.WaitZero() // wait for all goroutines to finish
	}
	os.Exit(0) // exit the program if no error
}

func ParseCommand() map[string]string {
	commandMap := make(map[string]string)
	download := flag.String("download", "", "input book id or url, like:download <bookid/url>")
	account := flag.String("account", "", "input account")
	password := flag.String("password", "", "input password")
	appType := flag.String("app", "sfacg", "input app type, like: app sfacg")
	search := flag.String("search", "", "input search keyword, like: search keyword")
	showConfig := flag.Bool("show", false, "show config, like: show config")
	flag.Parse()
	commandMap["book_id"] = ExtractBookID(*download)
	commandMap["account"] = *account
	commandMap["password"] = *password
	commandMap["app_type"] = *appType
	commandMap["key_word"] = *search
	if *showConfig {
		cfg.FormatJson(cfg.ReadConfig(""))
	}
	cfg.Vars.AppType = *appType
	return commandMap
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

func init() {
	cfg.ConfigInit()
	if len(os.Args) <= 1 {
		fmt.Println("please input command line parameters!")
		os.Exit(1)
	}
}
func main() {
	command, ExitProgram := ParseCommand(), false
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

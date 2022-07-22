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

func shell_book_download(downloadId any) {
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

func init() {
	cfg.ConfigInit()
	if len(os.Args) <= 1 {
		fmt.Println("please input parameters, like: sf login username password")
		fmt.Println("or: sf search keyword")
		fmt.Println("or: sf bookid")
		fmt.Println("or: sf url")
		fmt.Println("or: sf account")
		fmt.Println("or: sf password")
		os.Exit(1)
	} else {
		if cfg.Vars.Sfacg.Cookie != "" {
			if src.AccountDetailed() == "需要登录才能访问该资源" {
				fmt.Println("cookie is Invalid，please login first")
			}
		} else {
			fmt.Println("cookie is empty, please login first!")
		}
	}
}
func ParseCommand() map[string]string {
	commandMap := make(map[string]string)
	download := flag.String("download", "", "input book id or url, like:download <bookid/url>")
	account := flag.String("account", "", "input account, like: sf username")
	password := flag.String("password", "", "input password, like: sf password")
	appType := flag.String("app", "sfacg", "input app type, like: app sfacg")
	search := flag.String("search", "", "input search keyword, like: sf search keyword")
	flag.Parse()
	commandMap["book_id"] = ExtractBookID(*download)
	commandMap["account"] = *account
	commandMap["password"] = *password
	commandMap["app_type"] = *appType
	commandMap["key_word"] = *search
	cfg.Vars.AppType = *appType
	return commandMap
}
func main() {
	command, ExitProgram := ParseCommand(), false
	if command["account"] != "" || command["password"] != "" {
		ShellLoginAccount(command["account"], command["password"])
		ExitProgram = true
	} else if cfg.Vars.AppType == "sfacg" {
		if cfg.Vars.Sfacg.UserName == "" || cfg.Vars.Sfacg.Password == "" {
			fmt.Println("you must input account and password, like: sf username password")
			ExitProgram = true
		} else {
			src.LoginAccount(cfg.Vars.Sfacg.UserName, cfg.Vars.Sfacg.Password, 0)
		}
	}
	if command["key_word"] != "" {
		if NovelId := src.SearchBook(command["key_word"]); NovelId != "" {
			shell_book_download(NovelId)
		} else {
			fmt.Println("no book found")
		}
		ExitProgram = true
	}
	if command["book_id"] != "" {
		shell_book_download(command["book_id"])
		ExitProgram = true
	}

	if ExitProgram {
		os.Exit(0) // exit the program if no error
	}
}

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

func getBookID(url string) string {
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
		//src.LoginAccount(account, password, 0)
	}
}

func ShellBookByBookid(downloadId any) {
	switch downloadId.(type) {
	case string:
		start := src.BookInits{BookID: downloadId.(string), Index: 0, Locks: nil, ShowBook: true}
		start.SfacgBookInit()
	case []string:
		Locks := multi.NewGoLimit(7)
		for BookIndex, BookId := range downloadId.([]string) {
			Locks.Add()
			start := src.BookInits{BookID: BookId, Index: BookIndex, Locks: Locks, ShowBook: true}
			start.SfacgBookInit()
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

func main() {
	ExitProgram := false
	bookId := flag.String("id", "", "input book id, like: sf download bookid")
	sfacgUrl := flag.String("url", "", "input book id, like: sf url")
	account := flag.String("account", "", "input account, like: sf username")
	password := flag.String("password", "", "input password, like: sf password")
	appType := flag.String("app", "sfacg", "input app type, like: app sfacg")
	search := flag.String("search", "", "input search keyword, like: sf search keyword")
	flag.Parse()
	if *appType != "sfacg" {
		cfg.Vars.AppType = *appType
	} else {
		cfg.Vars.AppType = "sfacg"
	}
	if *account != "" || *password != "" {
		// if account and password are not empty, login
		ShellLoginAccount(*account, *password)
		ExitProgram = true
	} else {
		if cfg.Vars.Sfacg.UserName == "" || cfg.Vars.Sfacg.Password == "" {
			fmt.Println("you must input account and password, like: sf username password")
			ExitProgram = true
		} else {
			src.LoginAccount(cfg.Vars.Sfacg.UserName, cfg.Vars.Sfacg.Password, 0)
		}
	}
	if *search != "" {
		if NovelId := src.SearchBook(*search); NovelId != "" {
			ShellBookByBookid(NovelId)
		}
		ExitProgram = true
	}
	if *bookId != "" || *sfacgUrl != "" {
		if getBookID(*sfacgUrl) != "" {
			ShellBookByBookid(getBookID(*sfacgUrl))
		} else if *bookId != "" {
			ShellBookByBookid(*bookId)
		}
		ExitProgram = true
	}
	if ExitProgram {
		os.Exit(0) // exit the program if no error
	}
}

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
	if account != "" || password != "" { // if account and password are not empty, login
		if account == "" {
			fmt.Println("you must input account, like: sf username")
		} else if password == "" {
			fmt.Println("you must input password, like: sf password")
		} else {
			src.LoginAccount(account, password, 0)
		}
	} else {
		if cfg.Vars.Sfacg.UserName == "" || cfg.Vars.Sfacg.Password == "" {
			fmt.Println("you must input account and password, like: sf username password")
			os.Exit(1)
		} else {
			src.LoginAccount(cfg.Vars.Sfacg.UserName, cfg.Vars.Sfacg.Password, 0)
		}
	}
}
func ShellSearchBook(search string) {
	var input int
	if search != "" { // if search is not empty, search
		return
	}
	// if search keyword is not empty, search book and download
	if err := src.GetSearchDetailed(search); err == nil {
		fmt.Printf("please input the index of the book you want to download:")
		if _, err = fmt.Scanln(&input); err == nil {
			if input < len(cfg.Vars.BookInfoList) {
				ShellBookByBookid("", cfg.Vars.BookInfoList[input].NovelID, "")
			} else {
				fmt.Println("index out of range, please input again")
			}
		}
		os.Exit(0) // exit the program if no error
	} else {
		fmt.Println(err)
	}

}

func ShellBookByBookid(sfacgUrl, bookId string, downloadId any) {
	if getBookID(sfacgUrl) != "" {
		downloadId = getBookID(sfacgUrl)
	} else if bookId != "" {
		downloadId = bookId
	}
	if downloadId != "" {
		switch downloadId.(type) {
		case string:
			src.SfacgBookInit(downloadId.(string), 0, nil)
		case []string:
			Locks := multi.NewGoLimit(7)
			for BookIndex, BookId := range downloadId.([]string) {
				Locks.Add()
				src.SfacgBookInit(BookId, BookIndex, Locks)
			}
			Locks.WaitZero() // wait for all goroutines to finish
		}
		os.Exit(0) // exit the program if no error
	}

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
	bookId := flag.String("id", "", "input book id, like: sf download bookid")
	sfacgUrl := flag.String("url", "", "input book id, like: sf url")
	account := flag.String("account", "", "input account, like: sf username")
	password := flag.String("password", "", "input password, like: sf password")
	appType := flag.String("app", "", "input app type, like: app sfacg")
	search := flag.String("search", "", "input search keyword, like: sf search keyword")
	flag.Parse() // parse the flags from command line
	if *appType == "" {
		cfg.Vars.AppType = "sfacg"
	} else {
		cfg.Vars.AppType = *appType
	}
	cfg.SaveJson() // save the config file
	ShellLoginAccount(*account, *password)
	ShellSearchBook(*search)
	ShellBookByBookid(*sfacgUrl, *bookId, "")
}

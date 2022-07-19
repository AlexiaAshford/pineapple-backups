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
	bookID := regexp.MustCompile(`(\d+)`).FindStringSubmatch(url)
	if len(bookID) > 1 {
		return bookID[1]
	} else {
		return ""
	}
}

func Books(inputs any) {
	switch inputs.(type) {
	case string:
		src.BookInit(inputs.(string), 0, nil)
	case []string:
		Locks := multi.NewGoLimit(7)
		for BookIndex, BookId := range inputs.([]string) {
			Locks.Add()
			src.BookInit(BookId, BookIndex, Locks)
		}
		Locks.WaitZero() // wait for all goroutines to finish
	}
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
	if search != "" { // if search keyword is not empty, search book and download
		result := src.GetSearchDetailed(search)
		var input int
		fmt.Printf("please input the index of the book you want to download:")
		if _, err := fmt.Scanln(&input); err == nil {
			if input < len(result) {
				Books(result[input].NovelID)
			} else {
				fmt.Println("index out of range, please input again")
			}
		}
		os.Exit(0)
	}
}

func ShellBookByBookid(url, bookId string) {
	if url != "" || bookId != "" {
		var downloadId string
		if url != "" {
			if getBookID(url) != "" {
				downloadId = getBookID(url)
			} else {
				fmt.Println("you input url is not a book url")
			}
		} else {
			downloadId = bookId
		}
		Books(downloadId)
		os.Exit(0)
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
			if src.AccountDetailed() != "需要登录才能访问该资源" {
				fmt.Println("account is Valid，start to sf start to work, please wait...")
			} else {
				fmt.Println("account is Invalid，please login first, like: sf login username password")
			}
		} else {
			fmt.Println("account is Invalid，please login first")
		}
	}
}

func main() {
	bookId := flag.String("id", "", "input book id, like: sf download bookid")
	sfacgUrl := flag.String("url", "", "input book id, like: sf url")
	account := flag.String("account", "", "input account, like: sf username")
	password := flag.String("password", "", "input password, like: sf password")
	search := flag.String("search", "", "input search keyword, like: sf search keyword")
	flag.Parse() // parse the flags from command line

	ShellLoginAccount(*account, *password)
	ShellBookByBookid(*sfacgUrl, *bookId)
	ShellSearchBook(*search)
}

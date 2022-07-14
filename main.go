package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sf/src"
	"sf/src/config"
	"sf/src/threading"
)

func getBookID(url string) string {
	bookID := regexp.MustCompile(`(\d+)`).FindStringSubmatch(url)
	if len(bookID) > 1 {
		return bookID[1]
	} else {
		return ""
	}
}

func DownloadBookInit(inputs any) {
	switch inputs.(type) {
	case string:
		if BookData, err := src.GetBookDetailed(inputs.(string)); err == nil {
			fmt.Printf("开始下载:%s\n", BookData.NovelName)
			if err = ioutil.WriteFile(fmt.Sprintf("save/%v.txt", BookData.NovelName),
				[]byte(BookData.NovelName), 0777); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
			if src.GetCatalogue(BookData) {
				fmt.Printf("NovelName:%vdownload complete!", BookData.NovelName)
			}
		} else {
			fmt.Println("Error:", err)
		}
	case []string:
		Locks := threading.NewGoLimit(7)
		for BookIndex, BookId := range inputs.([]string) {
			Locks.Add()
			go func(bookId string, t *threading.GoLimit, BookIndex int) {
				fmt.Println(BookIndex)
				defer Locks.Done() // finish this goroutine when this function return
				if BookData, err := src.GetBookDetailed(bookId); err == nil {
					fmt.Printf("开始下载:%s\n", BookData.NovelName)
					if err = ioutil.WriteFile(fmt.Sprintf("save/%v.txt", BookData.NovelName),
						[]byte(BookData.NovelName), 0777); err != nil {
						fmt.Printf("Error: %v\n", err)
					}
					if src.GetCatalogue(BookData) {
						fmt.Printf("Index:%v\t\tNovelName:%vdownload complete!", BookIndex, BookData.NovelName)
					}
				} else {
					fmt.Println("Error:", err)
				}
			}(BookId, Locks, BookIndex)
		}
		Locks.WaitZero() // wait for all goroutines to finish
	}

}

func init() {
	config.NewMyJsonPro()
	if len(os.Args) <= 1 {
		fmt.Println("please input parameters, like: sf login username password")
		fmt.Println("or: sf search keyword")
		fmt.Println("or: sf bookid")
		fmt.Println("or: sf url")
		fmt.Println("or: sf account")
		fmt.Println("or: sf password")
		os.Exit(1)
	} else {
		if config.Var.Sfacg.Cookie != "" {
			if src.AccountDetailed() != "需要登录才能访问该资源" {
				fmt.Println("account is Valid，start to sf start to work, please wait...")
			} else {
				fmt.Println("account is Invalid，please login first")
			}
		} else {
			fmt.Println("account is Invalid，please login first")
		}
	}
}
func commandLine() (string, string, string, string, string) {
	bookId := flag.String("id", "", "input book id, like: sf download bookid")
	url := flag.String("url", "", "input book id, like: sf bookid")
	account := flag.String("account", "", "input account, like: sf username")
	password := flag.String("password", "", "input password, like: sf password")
	search := flag.String("search", "", "input search keyword, like: sf search keyword")
	flag.Parse()                                       // parse the flags from command line
	return *bookId, *url, *account, *password, *search // return command line parameters to main
}

func main() {
	bookId, url, account, password, search := commandLine() // get command line parameters
	if account != "" || password != "" {                    // if account and password are not empty, login
		if account == "" {
			fmt.Println("you must input account, like: sf username")
		} else if password == "" {
			fmt.Println("you must input password, like: sf password")
		} else {
			src.LoginAccount(account, password)
		}
	}
	if search != "" { // if search keyword is not empty, search book and download
		result := src.GetSearchDetailed(search)
		var input int
		fmt.Printf("please input the index of the book you want to download:")
		if _, err := fmt.Scanln(&input); err == nil {
			if input < len(result) {
				DownloadBookInit(result[input].NovelID)
			} else {
				fmt.Println("index out of range, please input again")
			}
		}
		os.Exit(0)
	}
	// if url or bookid is not empty, download book by url or bookid
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
		DownloadBookInit(downloadId)
		os.Exit(0)
	}
}

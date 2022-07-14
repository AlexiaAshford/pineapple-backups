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
	"strconv"
)

func getBookID(url string) string {
	bookID := regexp.MustCompile(`(\d+)`).FindStringSubmatch(url)
	if len(bookID) > 1 {
		return bookID[1]
	} else {
		return ""
	}
}

func downloadBook(input any) {
	var (
		bookId   string
		bookList []string
	)
	switch input.(type) {
	case string:
		bookId = input.(string)
	case int:
		bookId = strconv.Itoa(input.(int))
	case []string:
		bookList = input.([]string)
	}
	if bookId != "" && bookList == nil {
		BookData := src.GetBookDetailed(bookId)
		fmt.Printf("开始下载:%s\n", BookData.NovelName)
		if err := ioutil.WriteFile(fmt.Sprintf("save/%v.txt", BookData.NovelName),
			[]byte(BookData.NovelName), 0777); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		src.GetCatalogue(BookData)
	} else if bookList != nil && len(bookList) > 0 {

		ThreadLocks := threading.NewGoLimit(5)
		for _, bookId = range bookList {
			ThreadLocks.Add()
			go func(bookId string, t *threading.GoLimit) {
				defer ThreadLocks.Done()
				BookData := src.GetBookDetailed(bookId)
				fmt.Printf("开始下载:%s\n", BookData.NovelName)
				if err := ioutil.WriteFile(fmt.Sprintf("save/%v.txt", BookData.NovelName),
					[]byte(BookData.NovelName), 0777); err != nil {
					fmt.Printf("Error: %v\n", err)
				}
				src.GetCatalogue(BookData)
			}(bookId, ThreadLocks)
		}
		ThreadLocks.WaitZero()
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

//func main() {
//	inputs := os.Args[1:]
//	switch {
//	case inputs[0] == "login":
//		if len(inputs) >= 3 {
//			src.LoginAccount(inputs[1], inputs[2])
//		} else {
//			fmt.Println("parameters are not enough, please input username and password")
//		}
//	case inputs[0] == "search":
//		if len(inputs) >= 2 {
//			result := src.GetSearchDetailed(inputs[1])
//			var input int
//			fmt.Printf("please input the index of the book you want to download:")
//			if _, err := fmt.Scanln(&input); err == nil {
//				if input < len(result) {
//					downloadBook(result[input].NovelID)
//				} else {
//					fmt.Println("index out of range, please input again")
//				}
//			}
//
//		} else {
//			fmt.Println("parameters are not enough, please input keyword")
//		}
//	case inputs[0] == "download":
//		if len(inputs) >= 2 {
//			downloadBook(inputs[1])
//		} else {
//			fmt.Println("parameters are not enough, please input book id")
//		}
//	case inputs[0] == "url":
//		if len(inputs) >= 2 {
//			BookID := getBookID(inputs[1])
//			if BookID != "" {
//				downloadBook(BookID)
//			} else {
//				fmt.Println("parameters are not enough, please input url")
//			}
//		} else {
//			fmt.Println("parameters are not enough, please input book id")
//		}
//	default:
//		fmt.Println("the command is not exist, please input again")
//	}
//
//}

func commandLine() (string, string, string, string, string) {
	bookId := flag.String("id", "", "input book id, like: sf download bookid")
	url := flag.String("url", "", "input book id, like: sf bookid")
	account := flag.String("account", "", "input account, like: sf username")
	password := flag.String("password", "", "input password, like: sf password")
	search := flag.String("search", "", "input search keyword, like: sf search keyword")
	//var svar string
	//flag.StringVar(&svar, "svar", "bar", "a string var")
	flag.Parse()
	return *bookId, *url, *account, *password, *search
}

func main() {
	bookId, url, account, password, search := commandLine()
	if account != "" || password != "" {
		if account == "" {
			fmt.Println("you must input account, like: sf username")
		} else if password == "" {
			fmt.Println("you must input password, like: sf password")
		} else {
			src.LoginAccount(account, password)
		}
	}
	if search != "" {
		result := src.GetSearchDetailed(search)
		var input int
		fmt.Printf("please input the index of the book you want to download:")
		if _, err := fmt.Scanln(&input); err == nil {
			if input < len(result) {
				downloadBook(result[input].NovelID)
			} else {
				fmt.Println("index out of range, please input again")
			}
		}
		os.Exit(0)
	}
	if url != "" {
		BookID := getBookID(url)
		if BookID != "" {
			downloadBook(BookID)
		} else {
			fmt.Println("you input url is not a book url")
		}
		os.Exit(0)
	}
	if bookId != "" {
		downloadBook(bookId)
		os.Exit(0)
	}
}

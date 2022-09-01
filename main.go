package main

import (
	"fmt"
	"github.com/VeronicaAlexia/pineapple-backups/cfg"
	"github.com/VeronicaAlexia/pineapple-backups/src"
	"github.com/VeronicaAlexia/pineapple-backups/src/boluobao"
	"github.com/VeronicaAlexia/pineapple-backups/src/hbooker"
	"os"
	"strings"
)

func current_download_book(book_id string) {
	catalogue := src.SettingBooks(book_id) // get book catalogues
	if !catalogue.Test {
		fmt.Println(catalogue.BookMessage)
		os.Exit(1)
	} else {
		catalogue.GetDownloadsList()
	}
	if len(cfg.Current.DownloadList) > 0 {
		fmt.Println(len(cfg.Current.DownloadList), " chapters will be downloaded.")
		catalogue.ChapterBar = src.New(len(cfg.Current.DownloadList))
		catalogue.ChapterBar.Describe("working...")
		for _, file_name := range cfg.Current.DownloadList {
			catalogue.DownloadContent(file_name)
		}
		fmt.Printf("\nNovel:%v download complete!\n", cfg.Current.Book.NovelName)
		catalogue.MergeTextAndEpubFiles()
	} else {
		catalogue.MergeTextAndEpubFiles()
		cfg.ColorPrint(cfg.Current.Book.NovelName+" No chapter need to download!", 2|8)
	}
}

func shellUpdateLocalBook() {
	if cfg.Exist("./bookList.txt") && cfg.FileSize("./config.json") > 0 {
		LocalBookList := cfg.Write("./bookList.json", "", "r")
		for _, i := range strings.ReplaceAll(LocalBookList, "\n", "") {
			current_download_book(string(i))
		}
	} else {
		fmt.Println("bookList.txt not exist, create a new one!")
	}
}

func init() {
	if !cfg.Exist("./config.json") || cfg.FileSize("./config.json") == 0 {
		fmt.Println("config.json not exist, create a new one!")
	} else {
		fmt.Println("config.json exist, load config.json!")
		cfg.LoadJson()
	}
	if cfg.UpdateConfig() {
		cfg.SaveJson()
	}
	fmt.Println("you can input -h and --help to see the command list.")
}

func InitBookShelf() {
	var (
		bookshelf_index      int
		response_err         error
		bookshelf_book_index []int
		bookshelf_book_list  map[int][]map[string]string
	)
	if cfg.Vars.AppType == "sfacg" {
		bookshelf_book_list, response_err = boluobao.GET_BOOK_SHELF_INFORMATION()
	} else {
		bookshelf_book_list, response_err = hbooker.GET_BOOK_SHELF_INFORMATION()
	}
	if response_err != nil {
		var test_login_status bool
		fmt.Println("BookShelf Error:", response_err)
		if cfg.Vars.AppType == "sfacg" {
			test_login_status = src.AutoAccount()
		} else if cfg.Vars.AppType == "cat" {
			test_login_status = src.InputAccountToken()
		}
		if !test_login_status {
			fmt.Println("please login your account and password, like: sf account password")
			os.Exit(1)
		} else {
			InitBookShelf()
		}

	}

	fmt.Println("\nyou account is valid, start loading bookshelf information.")
	if len(bookshelf_book_list) == 1 {
		fmt.Println("you only have one bookshelf, default loading bookshelf index:1")
		bookshelf_index = 0
	} else {
		fmt.Println("please input bookshelf index:")
		bookshelf_index = cfg.InputInt(">", len(bookshelf_book_list))
	}
	for book_index, book := range bookshelf_book_list[bookshelf_index] {
		fmt.Println("book-index", book_index, "book-name:", book["novel_name"], "\t\tbook-id:", book["novel_id"])
		bookshelf_book_index = append(bookshelf_book_index, book_index)
	}
	for {
		if comment, ok := cfg.Console(); ok {
			if cfg.TestIntList(bookshelf_book_index, comment[0]) {
				shell([]string{"book", bookshelf_book_list[bookshelf_index][cfg.StrToInt(comment[0])]["novel_id"]})
			} else {
				shell(comment)
			}
		}
	}
}

func shell(inputs []string) {
	switch inputs[0] {
	case "q", "quit":
		fmt.Println("exit the program!")
		os.Exit(0)
	case "up", "update":
		shellUpdateLocalBook()
	case "a", "app":
		if cfg.TestList([]string{"sfacg", "cat"}, inputs[1]) {
			cfg.Vars.AppType = inputs[1]
		} else {
			fmt.Println("app type error, please input again.")
		}
	case "book", "download":
		if len(inputs) == 2 {
			if book_id := cfg.ExtractBookID(inputs[1]); book_id != "" {
				current_download_book(book_id)
			} else {
				fmt.Println("book id is empty, please input again.")
			}
		} else {
			fmt.Println("input book id or url, like:download <bookid/url>")
		}
	case "s", "search":
		if len(inputs) == 2 && inputs[1] != "" {
			current_download_book(src.SearchBook(inputs[1]))
		} else {
			fmt.Println("input book id or url, like:download <bookid/url>")
		}
	case "t", "token":

		if ok := src.InputAccountToken(); !ok {
			fmt.Println("you must input account and token.")
		}

	case "l", "login":

		if len(inputs) >= 3 {
			src.LoginAccount(inputs[1], inputs[2], 0)
		} else {
			fmt.Println("you must input account and password, like: sf account password")
		}
	default:
		fmt.Println("command not found,please input help to see the command list:", inputs[0])
	}

}

func main() {
	commentLine := cfg.CommandInit()
	if len(os.Args) > 1 && commentLine[0] != "console" {
		if cfg.Account != "" && cfg.Password != "" {
			shell([]string{"login", cfg.Account, cfg.Password})
		} else {
			src.TestAppTypeAndAccount()
		}
		if len(commentLine) > 0 {
			shell(commentLine)
		}
	} else {
		for _, s := range cfg.HelpMessage {
			fmt.Println("[info]", s)
		}
		InitBookShelf()

	}
}

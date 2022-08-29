package main

import (
	"fmt"
	"os"
	"sf/cfg"
	"sf/src"
	"sf/src/boluobao"
	"strconv"
	"strings"
)

func current_download_book(book_id any) {
	current_book_id := cfg.ExtractBookID(book_id.(string))
	if current_book_id == "" {
		fmt.Println("input book id or url is empty, please input again:")
		return
	}
	start := src.BookInits{BookID: current_book_id, Index: 0, Locks: nil, ShowBook: true}
	catalogue := start.SetBookInfo() // get book catalogues
	if !catalogue.TestBookResult {
		return
	}
	catalogue.GetDownloadsList()
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
	os.Exit(1)
}

func shellUpdateLocalBook() {
	if cfg.Exist("./bookList.txt") && cfg.FileSize("./config.json") > 0 {
		LocalBookList := cfg.Write("./bookList.json", "", "r")
		LocalBookList = strings.Replace(LocalBookList, "\n", "", -1)
		current_download_book(LocalBookList)
	} else {
		fmt.Println("bookList.txt not exist, create a new one!")
	}
}

func shellBookMain(inputs []string) {
	if cfg.Vars.AppType == "cat" {
		if len(inputs[1]) == 9 { // test if the input is hbooker book id
			current_download_book(inputs[1])
		} else {
			fmt.Println("hbooker bookid is 9 characters, please input again:")
		}
	} else {
		current_download_book(inputs[1])
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
	response := boluobao.GET_BOOK_SHELF_INFORMATION()
	if response.Status.HTTPCode == 200 {
		var bookshelf_index int
		fmt.Println("\nyou account is valid, start loading bookshelf information.")
		for index, value := range response.Data {
			fmt.Println("bookshelf index:", index+1, "bookshelf name:", value.Name)
		}
		if len(response.Data) == 1 {
			fmt.Println("you only have one bookshelf, default loading bookshelf index:1")
			bookshelf_index = 1
		} else {
			bookshelf_index = cfg.InputInt("please input bookshelf index:")
			if bookshelf_index > len(response.Data) {
				// if the index is valid (less than the length of the search result)
				bookshelf_index = cfg.InputInt("you input index is out of range, please input again:")
			}
		}
		var bookshelf_book_id []string
		var bookshelf_book_index []string
		for book_index, book := range response.Data[bookshelf_index].Expand.Novels {
			fmt.Println("book-index", book_index, "book-name:", book.NovelName, "\t\tbook-id:", book.NovelID)
			bookshelf_book_index = append(bookshelf_book_index, strconv.Itoa(book_index))
			bookshelf_book_id = append(bookshelf_book_id, strconv.Itoa(book.NovelID))
		}
		for {
			if comment, err := cfg.Console(); err {
				if cfg.TestList(bookshelf_book_index, comment[0]) {
					shell([]string{"book", bookshelf_book_id[cfg.StrToInt(comment[0])]})
				} else {
					shell(comment)
				}
			}
		}
	} else {
		if !src.AutoAccount() {
			fmt.Println("please login your account and password, like: sf account password")
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
			shellBookMain(inputs)
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
	if len(os.Args) > 1 {
		if cfg.Account != "" && cfg.Password != "" {
			shell([]string{"login", cfg.Account, cfg.Password})
		} else {
			cfg.Vars.AppType = cfg.App_type
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

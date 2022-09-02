package main

import (
	"fmt"
	"github.com/VeronicaAlexia/pineapple-backups/config"
	"github.com/VeronicaAlexia/pineapple-backups/src"
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
	if len(config.Current.DownloadList) > 0 {
		fmt.Println(len(config.Current.DownloadList), " chapters will be downloaded.")
		catalogue.ChapterBar = src.New(len(config.Current.DownloadList))
		catalogue.ChapterBar.Describe("working...")
		for _, file_name := range config.Current.DownloadList {
			catalogue.DownloadContent(file_name)
		}
		fmt.Printf("\nNovel:%v download complete!\n", config.Current.Book.NovelName)
		catalogue.MergeTextAndEpubFiles()
	} else {
		catalogue.MergeTextAndEpubFiles()
		config.ColorPrint(config.Current.Book.NovelName+" No chapter need to download!", 2|8)
	}
}

func shellUpdateLocalBook() {
	if config.Exist("./bookList.txt") && config.FileSize("./config.json") > 0 {
		LocalBookList := config.Write("./bookList.json", "", "r")
		for _, i := range strings.ReplaceAll(LocalBookList, "\n", "") {
			current_download_book(string(i))
		}
	} else {
		fmt.Println("bookList.txt not exist, create a new one!")
	}
}

func init() {
	if !config.Exist("./config.json") || config.FileSize("./config.json") == 0 {
		fmt.Println("config.json not exist, create a new one!")
	} else {
		fmt.Println("config.json exist, load config.json!")
		config.LoadJson()
	}
	if config.UpdateConfig() {
		config.SaveJson()
	}
	fmt.Println("you can input -h and --help to see the command list.")
}

func shell(inputs []string) {
	switch inputs[0] { // switch command
	case "up", "update":
		shellUpdateLocalBook()
	case "a", "app":
		if config.TestList([]string{"sfacg", "cat"}, inputs[1]) {
			config.Vars.AppType = inputs[1]
		} else {
			fmt.Println("app type error, please input again.")
		}
	case "book", "download":
		if len(inputs) == 2 {
			if book_id := config.ExtractBookID(inputs[1]); book_id != "" {
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
	commentLine := config.CommandInit()
	if len(os.Args) > 1 && commentLine[0] != "console" {
		if config.Account != "" && config.Password != "" {
			shell([]string{"login", config.Account, config.Password})
		} else {
			src.TestAppTypeAndAccount()
		}
		if len(commentLine) > 0 {
			shell(commentLine)
		}
	} else {
		for _, s := range config.HelpMessage {
			fmt.Println("[info]", s)
		}
		bookshelf_book_index, book_shelf_bookcase := src.InitBookShelf() // init bookshelf information
		for {
			if comment, ok := config.Console(); ok {
				if config.TestIntList(bookshelf_book_index, comment[0]) {
					shell([]string{"book", book_shelf_bookcase[config.StrToInt(comment[0])]["novel_id"]})
				} else if comment[0] == "load" || comment[0] == "bookshelf" {
					bookshelf_book_index, book_shelf_bookcase = src.InitBookShelf() // load bookshelf information
				} else if comment[0] == "quit" || comment[0] == "exit" {
					fmt.Println("exit the program!")
					os.Exit(0)
				} else {
					shell(comment)
				}
			}
		}

	}
}

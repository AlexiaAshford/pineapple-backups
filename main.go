package main

import (
	"fmt"
	"github.com/VeronicaAlexia/pineapple-backups/config"
	"github.com/VeronicaAlexia/pineapple-backups/config/book"
	"github.com/VeronicaAlexia/pineapple-backups/config/file"
	"github.com/VeronicaAlexia/pineapple-backups/config/tool"
	"github.com/VeronicaAlexia/pineapple-backups/src/app"
	"os"
	"strings"
)

func init() {
	if !config.Exist("./config.json") || config_file.SizeFile("./config.json") == 0 {
		fmt.Println("config.json not exist, create a new one!")
	} else {
		config.LoadJson()
	}
	if config.UpdateConfig() {
		config.SaveJson()
	}
}

func current_download_book_function(book_id string) {
	catalogue := app.SettingBooks(book_id) // get book catalogues
	if !catalogue.Test {
		fmt.Println(catalogue.BookMessage)
		os.Exit(1)
	} else {
		catalogue.GetDownloadsList()
	}
	if len(config.Current.DownloadList) > 0 {
		fmt.Println(len(config.Current.DownloadList), " chapters will be downloaded.")
		catalogue.ChapterBar = app.New(len(config.Current.DownloadList))
		catalogue.ChapterBar.Describe("working...")
		for _, file_name := range config.Current.DownloadList {
			catalogue.DownloadContent(file_name)
		}
		fmt.Printf("\nNovel:%v download complete!\n", config.Current.Book.NovelName)
	} else {
		fmt.Println(config.Current.Book.NovelName + " No chapter need to download!\n")
	}
	catalogue.MergeTextAndEpubFiles()
}

func update_local_booklist() {
	if config.Exist("./bookList.txt") {
		for _, i := range strings.ReplaceAll(config_file.Write("./bookList.json", "", "r"), "\n", "") {
			if !strings.Contains(string(i), "#") {
				current_download_book_function(string(i))
			}
		}
	} else {
		fmt.Println("bookList.txt not exist, create a new one!")
	}
}
func shell(inputs []string) {
	switch inputs[0] { // switch command
	case "up", "update":
		update_local_booklist()
	case "a", "app":
		if tool.TestList([]string{"sfacg", "cat"}, inputs[1]) {
			config.Vars.AppType = inputs[1]
		} else {
			fmt.Println("app type error, please input again.")
		}
	case "d", "b", "book", "download":
		if len(inputs) == 2 {
			if book_id := config_book.ExtractBookID(inputs[1]); book_id != "" {
				current_download_book_function(book_id)
			} else {
				fmt.Println("book id is empty, please input again.")
			}
		} else {
			fmt.Println("input book id or url, like:download <bookid/url>")
		}
	case "s", "search":
		if len(inputs) == 2 && inputs[1] != "" {
			s := app.Search{Keyword: inputs[1], Page: 0}
			current_download_book_function(s.SearchBook())
		} else {
			fmt.Println("input book id or url, like:download <bookid/url>")
		}
	case "l", "t", "token", "login":
		if config.Vars.AppType == "cat" {
			if ok := app.InputAccountToken(); !ok {
				fmt.Println("you must input account and token.")
			}
		} else if len(inputs) >= 3 {
			app.LoginAccount(inputs[1], inputs[2], 0)
		} else {
			fmt.Println("you must input account and password, like: sf account password")
		}
	default:
		fmt.Println("command not found,please input help to see the command list:", inputs[0])
	}
}

func shell_run_console_and_bookshelf() {
	bookshelf_book_index, book_shelf_bookcase := app.InitBookShelf() // init bookshelf information
	for {
		if comment, ok := config.ConsoleInput(); ok {
			if tool.TestIntList(bookshelf_book_index, comment[0]) {
				current_download_book_function(book_shelf_bookcase[tool.StrToInt(comment[0])]["novel_id"])
			} else if comment[0] == "load" || comment[0] == "bookshelf" {
				bookshelf_book_index, book_shelf_bookcase = app.InitBookShelf() // load bookshelf information
			} else if comment[0] == "quit" || comment[0] == "exit" {
				fmt.Println("exit the program!")
				os.Exit(0)
			} else {
				shell(comment)
			}
		}
	}
}

func main() {
	commentLine := config.InitCommand()
	config.Vars.ThreadNum = commentLine.MaxThread
	config.Vars.AppType = commentLine.AppType
	fmt.Println("current app type:", config.Vars.AppType)
	if len(os.Args) > 1 {
		if commentLine.Account != "" && commentLine.Password != "" {
			shell([]string{"login", commentLine.Account, commentLine.Password})
		} else if commentLine.Login {
			app.TestAppTypeAndAccount()
		} else if commentLine.BookId != "" {
			current_download_book_function(commentLine.BookId)
		} else if commentLine.SearchKey != "" {
			s := app.Search{Keyword: commentLine.SearchKey, Page: 0}
			current_download_book_function(s.SearchBook())
		} else if commentLine.Update {
			update_local_booklist()
		} else if commentLine.Token {
			app.InputAccountToken()
		} else {
			shell_run_console_and_bookshelf()
		}
	} else {
		for _, message := range config.HelpMessage {
			fmt.Println("[info]", message)
		}
		shell_run_console_and_bookshelf()
	}
}

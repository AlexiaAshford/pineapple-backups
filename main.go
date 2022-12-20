package main

import (
	"fmt"
	"github.com/VeronicaAlexia/BoluobaoAPI/boluobao/account"
	"github.com/VeronicaAlexia/BoluobaoAPI/boluobao/task"
	BoluobaoConfig "github.com/VeronicaAlexia/BoluobaoAPI/pkg/config"
	HbookerConfig "github.com/VeronicaAlexia/HbookerAPI/config"
	"github.com/VeronicaAlexia/pineapple-backups/config"
	"github.com/VeronicaAlexia/pineapple-backups/pkg/file"
	"github.com/VeronicaAlexia/pineapple-backups/pkg/threading"
	"github.com/VeronicaAlexia/pineapple-backups/pkg/tools"
	"github.com/VeronicaAlexia/pineapple-backups/src"
	"github.com/VeronicaAlexia/pineapple-backups/src/app/hbooker"
	"github.com/urfave/cli"
	"log"
	"os"
	"strconv"
	"strings"
)

func init() {
	if !config.Exist("./config.json") || file.SizeFile("./config.json") == 0 {
		fmt.Println("config.json not exist, create a new one!")
	} else {
		config.LoadJson()
	}
	config.UpdateConfig()
	InitApp := cli.NewApp()
	InitApp.Name = "pineapple-backups"
	InitApp.Version = "V.1.6.9"
	InitApp.Usage = "https://github.com/VeronicaAlexia/pineapple-backups"
	InitApp.Flags = config.Args

	HbookerConfig.AppConfig.AppVersion = config.Apps.Cat.Params.AppVersion
	HbookerConfig.AppConfig.Account = config.Apps.Cat.Params.Account
	HbookerConfig.AppConfig.LoginToken = config.Apps.Cat.Params.LoginToken
	HbookerConfig.AppConfig.DeviceToken = config.Apps.Cat.Params.DeviceToken

	BoluobaoConfig.AppConfig.App = true
	BoluobaoConfig.AppConfig.AppKey = config.Vars.AppKey
	BoluobaoConfig.AppConfig.DeviceId = "240a90cc-4c40-32c7-b44e-d4cf9e670605"
	BoluobaoConfig.AppConfig.Cookie = config.Apps.Sfacg.Cookie

	InitApp.Action = func(c *cli.Context) {
		fmt.Println("you can input -h and --help to see the command list.")
		if config.CommandLines.AppType != "cat" && config.CommandLines.AppType != "sfacg" {
			fmt.Println(config.CommandLines.AppType, "app type error, default app type is cat.")
			config.CommandLines.AppType = "cat" // default app type is cat
		}
	}
	if err := InitApp.Run(os.Args); err != nil {
		log.Fatal(err)
	}
	config.Vars.ThreadNum = config.CommandLines.MaxThread
	config.Vars.AppType = config.CommandLines.AppType
	config.Vars.Epub = config.CommandLines.Epub

	fmt.Println("current app type:", config.Vars.AppType)
}

func current_download_book_function(book_id string) {
	catalogue := src.SettingBooks(book_id) // get book catalogues
	if !catalogue.Test {
		fmt.Println(catalogue.BookMessage)
		os.Exit(1)
	} else {
		catalogue.GetDownloadsList()
	}
	if len(config.Current.DownloadList) > 0 {
		thread := threading.NewGoLimit(uint(config.Vars.MaxRetry))
		fmt.Println(len(config.Current.DownloadList), " chapters will be downloaded.")
		catalogue.ChapterBar = src.New(len(config.Current.DownloadList))
		catalogue.ChapterBar.Describe("working...")
		for _, file_name := range config.Current.DownloadList {
			thread.Add()
			go catalogue.DownloadContent(thread, file_name)
		}
		thread.WaitZero()
		fmt.Printf("\nNovel:%v download complete!\n", config.Current.NewBooks["novel_name"])
	} else {
		fmt.Println(config.Current.NewBooks["novel_name"] + " No chapter need to download!\n")
	}
	catalogue.MergeTextAndEpubFiles()
}

func update_local_booklist() {
	if config.Exist("./bookList.txt") {
		for _, i := range strings.ReplaceAll(file.Open("./bookList.json", "", "r"), "\n", "") {
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
		if tools.TestList([]string{"sfacg", "cat"}, inputs[1]) {
			config.Vars.AppType = inputs[1]
		} else {
			fmt.Println("app type error, please input again.")
		}
	case "d", "b", "book", "download":
		if len(inputs) == 2 {
			if book_id := config.FindID(inputs[1]); book_id != "" {
				current_download_book_function(book_id)
			} else {
				fmt.Println("book id is empty, please input again.")
			}
		} else {
			fmt.Println("input book id or url, like:download <bookid/url>")
		}
	case "s", "search":
		if len(inputs) == 2 && inputs[1] != "" {
			s := src.Search{Keyword: inputs[1], Page: 0}
			current_download_book_function(s.SearchBook())
		} else {
			fmt.Println("input search keyword, like:search <keyword>")
		}

	case "l", "login":
		if config.Vars.AppType == "cat" && len(inputs) >= 3 {
			hbooker.GET_LOGIN_TOKEN(inputs[1], inputs[2])
		} else if config.Vars.AppType == "sfacg" && len(inputs) >= 3 {
			src.LoginAccount(inputs[1], inputs[2], 0)
		} else {
			fmt.Println("you must input account and password, like: -login account password")
		}
	case "t", "token":
		if ok := src.InputAccountToken(); !ok {
			fmt.Println("you must input account and token.")
		}
	default:
		fmt.Println("command not found,please input help to see the command list:", inputs[0])
	}
}

func shell_run_console_and_bookshelf() {
	bookshelf_book_index, book_shelf_bookcase := src.InitBookShelf() // init bookshelf information
	for {
		if comment := tools.GET(">"); comment != nil {
			if tools.TestIntList(bookshelf_book_index, comment[0]) {
				current_download_book_function(book_shelf_bookcase[tools.StrToInt(comment[0])]["novel_id"])
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

func main() {
	if len(os.Args) > 1 {
		if config.CommandLines.Account != "" && config.CommandLines.Password != "" {
			shell([]string{"login", config.CommandLines.Account, config.CommandLines.Password})
		} else if config.CommandLines.Login {
			src.TestAppTypeAndAccount()
		} else if config.CommandLines.BookId != "" {
			current_download_book_function(config.CommandLines.BookId)
		} else if config.CommandLines.SearchKey != "" {
			s := src.Search{Keyword: config.CommandLines.SearchKey, Page: 0}
			current_download_book_function(s.SearchBook())
		} else if config.CommandLines.Update {
			update_local_booklist()
		} else if config.CommandLines.Token {
			src.InputAccountToken()
		} else {
			shell_run_console_and_bookshelf()
		}
	} else {
		if config.Vars.AppType == "cat" {
			// recommend list for hbooker app
			if book_id := src.NEW_RECOMMEND().GET_HBOOKER_RECOMMEND(); book_id != "" {
				current_download_book_function(book_id)
			}
		} else if config.Vars.AppType == "sfacg" {
			accounts := account.GET_ACCOUNT_INFORMATION()
			if accounts.Data.AccountID > 0 {
				Tasks := task.Task{AccountId: strconv.Itoa(accounts.Data.AccountID)}
				Tasks.COMPLETE_ALL()
			} else {
				fmt.Println("you need to login to use the automatic sign-in and task function")
			}
		}
		for _, message := range config.HelpMessage {
			fmt.Println("[info]", message)
		}
		shell_run_console_and_bookshelf()
	}
}

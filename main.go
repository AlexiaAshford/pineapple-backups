package main

import (
	"fmt"
	"github.com/VeronicaAlexia/BoluobaoAPI/boluobao"
	BoluobaoConfig "github.com/VeronicaAlexia/BoluobaoAPI/pkg/config"
	HbookerAPI "github.com/VeronicaAlexia/HbookerAPI/ciweimao/book"
	HbookerConfig "github.com/VeronicaAlexia/HbookerAPI/pkg/config"
	"github.com/VeronicaAlexia/pineapple-backups/config"
	"github.com/VeronicaAlexia/pineapple-backups/pkg/file"
	"github.com/VeronicaAlexia/pineapple-backups/pkg/threading"
	"github.com/VeronicaAlexia/pineapple-backups/pkg/tools"
	"github.com/VeronicaAlexia/pineapple-backups/src"
	"github.com/urfave/cli"
	"log"
	"os"
	"strconv"
	"strings"
)

var bookShelfList map[string]int

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
	HbookerConfig.AppConfig.AppVersion = config.Apps.Hbooker.AppVersion
	HbookerConfig.AppConfig.Account = config.Apps.Hbooker.Account
	HbookerConfig.AppConfig.LoginToken = config.Apps.Hbooker.LoginToken
	HbookerConfig.AppConfig.DeviceToken = config.Apps.Hbooker.DeviceToken

	BoluobaoConfig.AppConfig.App = true // set boluobao app mode
	BoluobaoConfig.AppConfig.AppKey = "FMLxgOdsfxmN!Dt4"
	BoluobaoConfig.AppConfig.DeviceId = "240a90cc-4c40-32c7-b44e-d4cf9e670605"
	BoluobaoConfig.AppConfig.Cookie = config.Apps.Sfacg.Cookie

	InitApp.Action = func(c *cli.Context) {
		fmt.Println("you can input -h and --help to see the command list.")
		if config.Command.AppType != "cat" && config.Command.AppType != "sfacg" {
			fmt.Println(config.Command.AppType, "app type error")
			os.Exit(1)
		}
	}
	if err := InitApp.Run(os.Args); err != nil {
		log.Fatal(err)
	}
	config.Vars.ThreadNum = config.Command.MaxThread
	config.Vars.AppType = config.Command.AppType
	//config.Vars.Epub = config.Command.Epub

	fmt.Println("current app type:", config.Vars.AppType)
}

func current_download_book_function(book_id string) {
	catalogue := src.SettingBooks(book_id) // get book catalogues
	if !catalogue.Test {
		fmt.Println(catalogue.BookMessage)
		os.Exit(1)
	}
	DownloadList := catalogue.GetDownloadsList()

	if DownloadList != nil && len(DownloadList) > 0 {
		thread := threading.NewGoLimit(uint(config.Vars.MaxRetry))
		fmt.Println(len(DownloadList), " chapters will be downloaded.")
		catalogue.ChapterBar = src.New(len(DownloadList))
		catalogue.ChapterBar.Describe("working...")
		//fmt.Println(DownloadList)
		for _, chapterID := range DownloadList {
			thread.Add()
			go catalogue.DownloadContent(thread, chapterID)
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
	case "bs", "bookshelf":
		if len(bookShelfList) > 0 {
			if len(inputs) == 2 {
				value, ok := bookShelfList[inputs[1]]
				if ok {
					current_download_book_function(strconv.Itoa(value))
				} else {
					fmt.Println(inputs[1], "key not found")
				}
			}
		} else {
			fmt.Println("bookshelf is empty, please login and update bookshelf.")
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
			//hbooker.GET_LOGIN_TOKEN(inputs[1], inputs[2])
			fmt.Println("hbooker login function is not available now.")
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
			} else if comment[0] == "q" || comment[0] == "quit" {
				shell(comment)
			} else {
				fmt.Println("input book index to download book")
				fmt.Println("input load to reload bookshelf")
				fmt.Println("input quit to exit bookshelf")
			}
		}
	}
}

func main() {

	if len(os.Args) > 1 {
		if config.Command.Account != "" && config.Command.Password != "" {
			shell([]string{"login", config.Command.Account, config.Command.Password})
		} else if config.Command.Login {
			src.TestAppTypeAndAccount()
		} else if config.Command.BookId != "" {
			current_download_book_function(config.Command.BookId)
		} else if config.Command.SearchKey != "" {
			s := src.Search{Keyword: config.Command.SearchKey, Page: 0}
			current_download_book_function(s.SearchBook())
		} else if config.Command.Update {
			update_local_booklist()
		} else if config.Command.Token {
			src.InputAccountToken()
		} else {
			shell_run_console_and_bookshelf()
		}
	} else {
		for _, message := range config.HelpMessage {
			fmt.Println("[info]", message)
		}
		if config.Vars.AppType == "cat" {
			// recommend list for hbooker app
			//if book_id := recommend.NEW_RECOMMEND().CHANGE_NEW_RECOMMEND(); book_id != "" {
			//	current_download_book_function(book_id)
			//}
			response := HbookerAPI.GET_BOOK_INFORMATION("100280239")
			//fmt.Println(response)
			if response.Code == "100000" {
				shell_run_console_and_bookshelf()
			} else {
				fmt.Println("hbooker error:", response.Tip)
			}
		} else if config.Vars.AppType == "sfacg" {
			boluobao.API.User.UserInformation()
			if bs := src.NewChoiceBookshelf(); bs != nil {
				bs.NewSfacgBookshelf()
				bookShelfList = bs.ShelfBook
			}
			shell(tools.GET(">"))

			//shell_run_console_and_bookshelf()
			//if accounts.Data.AccountID > 0 {
			//	Tasks := task.Task{AccountId: strconv.Itoa(accounts.Data.AccountID)}
			//	Tasks.NovelCompleteTas()
			//	shell_run_console_and_bookshelf()
			//} else {
			//	fmt.Println("you need to login to use the automatic sign-in and task function")
			//}
		}
	}
}

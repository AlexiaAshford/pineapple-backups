package main

import (
	"fmt"
	"github.com/AlexiaVeronica/boluobaoLib"
	"github.com/AlexiaVeronica/hbookerLib"
	"github.com/AlexiaVeronica/pineapple-backups/config"
	"github.com/AlexiaVeronica/pineapple-backups/pkg/command"
	"github.com/AlexiaVeronica/pineapple-backups/pkg/file"
	"github.com/AlexiaVeronica/pineapple-backups/pkg/threading"
	"github.com/AlexiaVeronica/pineapple-backups/pkg/tools"
	"github.com/AlexiaVeronica/pineapple-backups/src"
	"os"
	"strings"
)

var bookShelfList map[string]string

func init() {
	if !config.Exist("./config.json") || file.SizeFile("./config.json") == 0 {
		fmt.Println("config.json not exist, create a new one!")
	} else {
		config.LoadJson()
	}
	config.UpdateConfig()

	command.NewApp()
	config.APP.Hbooker = &config.Hbooker{Client: hbookerLib.NewClient(hbookerLib.WithAccountAndLoginToken(config.Apps.Hbooker.Account, config.Apps.Hbooker.LoginToken))}
	config.APP.SFacg = &config.SFacg{Client: boluobaoLib.NewClient(boluobaoLib.WithCookie(config.Apps.Sfacg.Cookie))}
	fmt.Println("current app type:", command.Command.AppType)
}

func currentDownloadBookFunction(bookId string) {
	catalogue, err := src.SettingBooks(bookId) // get book catalogues
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	downloadList, err := catalogue.GetDownloadsList()
	if err != nil {
		fmt.Println(err)
		return
	}
	thread := threading.NewGoLimit(uint(32))
	fmt.Println(len(downloadList), " chapters will be downloaded.")
	catalogue.ChapterBar = src.New(len(downloadList))
	catalogue.ChapterBar.Describe("working...")
	for _, chapterID := range downloadList {
		thread.Add()
		go catalogue.DownloadContent(thread, chapterID)
	}
	thread.WaitZero()
	catalogue.MergeTextAndEpubFiles()
}

func updateLocalBooklist() {
	if config.Exist("./bookList.txt") {
		for _, i := range strings.ReplaceAll(file.Open("./bookList.json", "", "r"), "\n", "") {
			if !strings.Contains(string(i), "#") {
				currentDownloadBookFunction(string(i))
			}
		}
	} else {
		fmt.Println("bookList.txt not exist, create a new one!")
	}
}
func shellSwitch(inputs []string) {
	switch inputs[0] { // switch command
	case "up", "update":
		updateLocalBooklist()
	case "a", "app":
		if tools.TestList([]string{"sfacg", "cat"}, inputs[1]) {
			command.Command.AppType = inputs[1]
		} else {
			fmt.Println("app type error, please input again.")
		}
	case "d", "b", "book", "download":
		if len(inputs) == 2 {
			if bookId := config.FindID(inputs[1]); bookId != "" {
				currentDownloadBookFunction(bookId)
			} else {
				fmt.Println("book id is empty, please input again.")
			}
		} else {
			fmt.Println("input book id or url, like:download <bookid/url>")
		}
	case "bs", "bookshelf":
		if len(bookShelfList) > 0 && len(inputs) == 2 {
			if value, ok := bookShelfList[inputs[1]]; ok {
				currentDownloadBookFunction(value)
			} else {
				fmt.Println(inputs[1], "key not found")
			}
		} else {
			fmt.Println("bookshelf is empty, please login and update bookshelf.")
		}
	case "s", "search":
		if len(inputs) == 2 && inputs[1] != "" {
			s := src.Search{Keyword: inputs[1], Page: 0}
			currentDownloadBookFunction(s.SearchBook())
		} else {
			fmt.Println("input search keyword, like:search <keyword>")
		}

	case "l", "login":
		if command.Command.AppType == "cat" && len(inputs) >= 3 {
			config.APP.Hbooker.Client.SetDefaultParams(inputs[1], inputs[2])
		} else if command.Command.AppType == "sfacg" && len(inputs) >= 3 {
			src.LoginAccount(inputs[1], inputs[2])
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

func shell(messageOpen bool) {
	if messageOpen {
		for _, message := range config.HelpMessage {
			fmt.Println("[info]", message)
		}
	}
	bookshelf, err := src.NewChoiceBookshelf()
	if err != nil {
		fmt.Println(err)
	} else {
		bookShelfList = bookshelf
	}

	for {
		if inputRes := tools.GET(">"); len(inputRes) > 0 {
			shellSwitch(inputRes)
		}
	}

}

func main() {

	if len(os.Args) > 1 {
		if command.Command.Account != "" && command.Command.Password != "" {
			src.LoginAccount(command.Command.Account, command.Command.Password)
		} else if command.Command.BookID != "" {
			currentDownloadBookFunction(command.Command.BookID)
		} else if command.Command.SearchKey != "" {
			s := src.Search{Keyword: command.Command.SearchKey, Page: 0}
			currentDownloadBookFunction(s.SearchBook())
		} else if command.Command.Update {
			updateLocalBooklist()
		} else if command.Command.Token {
			src.InputAccountToken()
		} else if command.Command.BookShelf {
			shell(false)
		} else {
			fmt.Println("command not found,please input help to see the command list:", os.Args[1])
		}
	} else {
		shell(true)
	}
}

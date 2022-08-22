package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"sf/cfg"
	"sf/src"
	"sf/struct"
	"strings"
)

func shellBookDownload(downloadId any) {
	switch downloadId.(type) {
	case string:
		start := src.BookInits{BookID: downloadId.(string), Index: 0, Locks: nil, ShowBook: true}
		start.InitEpubFile()
		catalogues := start.DownloadBookInit() // get book catalogues
		if catalogues.TestBookResult {
			catalogues.InitCatalogue()
		}
	case []string:
		Locks := cfg.NewGoLimit(7)
		for BookIndex, BookId := range downloadId.([]string) {
			Locks.Add()
			start := src.BookInits{BookID: BookId, Index: BookIndex, Locks: Locks, ShowBook: true}
			catalogues := start.DownloadBookInit() // get book catalogues
			start.InitEpubFile()
			if catalogues.TestBookResult {
				catalogues.InitCatalogue()
			}
		}
		Locks.WaitZero() // wait for all goroutines to finish
	}
	os.Exit(0) // exit the program if no error
}

func shellLoginAccount(inputs []string) bool {
	if cfg.Vars.AppType == "sfacg" {
		if len(inputs) >= 3 {
			src.LoginAccount(inputs[1], inputs[2], 0)
		} else {
			fmt.Println("you must input account and password, like: sf account password")
		}
	} else if cfg.Vars.AppType == "cat" {
		if ok := src.InputAccountToken(); !ok {
			fmt.Println("you must input account and token.")
		}
	}
	return true
}

func shellUpdateLocalBook() {
	if cfg.Exist("./bookList.txt") && cfg.FileSize("./config.json") > 0 {
		LocalBookList := cfg.Write("./bookList.json", "", "r")
		LocalBookList = strings.Replace(LocalBookList, "\n", "", -1)
		shellBookDownload(LocalBookList)
	} else {
		fmt.Println("bookList.txt not exist, create a new one!")
	}
}

func shellBookMain(inputs []string) {
	if len(inputs) == 2 {
		if cfg.Vars.AppType == "cat" {
			if len(inputs[1]) == 9 { // test if the input is hbooker book id
				shellBookDownload(inputs[1])
			} else {
				fmt.Println("hbooker bookid is 9 characters, please input again:")
			}
		} else {
			shellBookDownload(inputs[1])
		}
	} else {
		fmt.Println("input book id or url, like:download <bookid/url>")
	}
}

func shellSearchBookMain(inputs []string) {
	if len(inputs) == 2 {
		if NovelId := src.SearchBook(inputs[1]); NovelId != "" {
			shellBookDownload(NovelId)
		} else {
			fmt.Println("No found search book, please input again:")
		}
	} else {
		fmt.Println("input book id or url, like:download <bookid/url>")
	}
}

func ParseCommandLine() _struct.Command {
	bookid := flag.String("download", "", "input book id or url")
	account := flag.String("account", "", "input account")
	password := flag.String("password", "", "input password")
	appType := flag.String("app", "sfacg", "input app type, like: app sfacg")
	search := flag.String("search", "", "input search keyword, like: search keyword")
	thread := flag.Int("max", 0, "input thread number, like: thread 1")
	showInfo := flag.Bool("show", false, "show config, like: show config")
	update := flag.Bool("update", false, "update config, like: update config")

	flag.Parse()
	if *thread > 0 && *thread < 64 {
		cfg.Vars.MaxThreadNumber = *thread
	}
	if *account != "" && *password != "" {
		shellConsole([]string{"login", *account, *password})
	} else {
		cfg.Vars.AppType = *appType
		src.TestAppTypeAndAccount()
	}
	return _struct.Command{Download: *bookid, Search: *search, ShowConfig: *showInfo, Update: *update}
}

func shellConsole(inputs []string) {
	switch inputs[0] {
	case "a", "app":
		cfg.Vars.AppType = inputs[1]
		src.TestAppTypeAndAccount()
	case "q", "quit":
		os.Exit(0)
	case "uo", "update":
		shellUpdateLocalBook()
	case "h", "help":
		fmt.Println("help:")
	case "show", "test":
		cfg.FormatJson(cfg.ReadConfig(""))
	case "book", "download":
		shellBookMain(inputs)
	case "s", "search":
		shellSearchBookMain(inputs)
	case "l", "login":
		shellLoginAccount(inputs)
	default:
		fmt.Println("command not found,please input help to see the command list:", inputs[0])
	}

}
func init() {
	if !cfg.Exist("./config.json") || cfg.FileSize("./config.json") == 0 {
		fmt.Println("config.json not exist, create a new one!")
	} else {
		cfg.LoadJson()
	}
	if cfg.UpdateConfig() {
		cfg.SaveJson()
	}
}

func main() {
	if len(os.Args) <= 1 {
		//for _, v := range hbooker.GetChangeRecommend() {
		//	fmt.Println(v.BookName)
		//	fmt.Println(v.BookID)
		//}
		for {
			fmt.Println("input help to see the command list:")
			fmt.Println("input quit to quit")
			fmt.Println("input download <bookid/url> to download book")
			fmt.Println("input search <keyword> to search book")
			fmt.Println("input show to show config")
			fmt.Println("input update config to update config by config.json")
			fmt.Println("input login <account> <password> to login account")
			fmt.Println("input app <app app keyword> to change app type")
			fmt.Println("input max <thread> to change max thread number")
			fmt.Println("you can input command like this: download <bookid/url>")
			fmt.Println("you can input non-existent command to exit the program")

			spaceRe, _ := regexp.Compile(`\s+`)
			inputs := spaceRe.Split(strings.TrimSpace(cfg.Input(">")), -1)
			if len(inputs) > 1 {
				shellConsole(inputs)
			} else if inputs[0] != "" {
				fmt.Println("you must input command, like: sf command")
			}
			os.Exit(1)
		}
	} else {
		var CommandLine []string
		ArgsCommandLine := ParseCommandLine()
		if ArgsCommandLine.ShowConfig {
			CommandLine = []string{"show", "config"}
		}
		if ArgsCommandLine.Download != "" {
			CommandLine = []string{"download", ArgsCommandLine.Download}
		}
		if ArgsCommandLine.Search != "" {
			CommandLine = []string{"search", ArgsCommandLine.Search}
		}
		if len(CommandLine) > 0 {
			shellConsole(CommandLine)
		}
	}
}

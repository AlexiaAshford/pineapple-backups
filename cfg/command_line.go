package cfg

import (
	"flag"
	"fmt"
)

func ParseCommandLine() map[string]any {
	CommandMap, AppList := make(map[string]any), []string{"sfacg", "cat"}
	download := flag.String("download", "", "input book id or url, like:download <bookid/url>")
	account := flag.String("account", "", "input account")
	password := flag.String("password", "", "input password")
	appType := flag.String("app", "sfacg", "input app type, like: app sfacg")
	search := flag.String("search", "", "input search keyword, like: search keyword")
	Thread := flag.Int("max", 0, "input thread number, like: thread 1")
	showConfig := flag.Bool("show", false, "show config, like: show config")
	flag.Parse()
	CommandMap["book"] = ExtractBookID(*download)
	CommandMap["max"] = fmt.Sprintf("%d", *Thread)
	CommandMap["account"] = *account
	CommandMap["password"] = *password
	CommandMap["app_type"] = *appType
	CommandMap["key_word"] = *search
	if *showConfig {
		// print config json file information to console
		FormatJson(ReadConfig(""))
	}
	// check app type and cheng edit config
	if TestList(AppList, *appType) {
		Vars.AppType = *appType
	} else {
		fmt.Printf("app type %v is invalid, please input again:", *appType)
		Vars.AppType = "sfacg"
	}
	if *Thread != 0 {
		if *Thread >= 64 {
			fmt.Println("thread number is too large, please input again:")
		} else {
			Vars.MaxThreadNumber = *Thread
			fmt.Println("change thread number to:", *Thread, "thread")
		}
	}
	return CommandMap
}

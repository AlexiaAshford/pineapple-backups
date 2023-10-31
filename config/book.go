package config

import (
	"fmt"
	"github.com/AlexiaVeronica/pineapple-backups/pkg/command"
	"regexp"
)

func FindID(url string) string {
	currentBookId := regexp.MustCompile(`(\d+)`).FindStringSubmatch(url)
	if len(currentBookId) > 1 {
		if command.Command.AppType == "cat" {
			if len(currentBookId[1]) != 9 { // test if the input is hbooker book id
				fmt.Println("hbooker bookid is 9 characters, please input again:")
			} else {
				return currentBookId[1]
			}
		}
	}
	return ""
}

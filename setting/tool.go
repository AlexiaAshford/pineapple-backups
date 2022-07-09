package setting

import (
	"fmt"
	"regexp"
	"strconv"
)

func InputInt(explain string) int {
	var InputInfo string
	fmt.Printf(explain)
	_, err := fmt.Scanln(&InputInfo)
	if err != nil {
		fmt.Println(err)
		return 404
	}
	Info, err := strconv.Atoi(InputInfo)
	if err != nil {
		fmt.Println("enter value is not a number")
		return 404
	} else {
		return Info
	}
}

func RegexpName(Name string) string {
	return regexp.MustCompile(`[\\/:*?"<>|]`).ReplaceAllString(Name, "")
}

func FormatCookie(SFCommunity, sessionAPP string) string {
	return fmt.Sprintf(".SFCommunity=%v;session_APP=%v;", SFCommunity, sessionAPP)
}

package cfg

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strconv"
)

func RegexpName(Name string) string {
	return regexp.MustCompile(`[\\/:*?"<>|]`).ReplaceAllString(Name, "")
}

func Mkdir(filePath string) {
	dirPath, _ := os.Getwd()
	if err := os.MkdirAll(path.Join(dirPath, filePath), os.ModePerm); err != nil {
		fmt.Println(err)
	}
}

// input int
//func InputInt(introduction string) int {
//	var input int
//	// if search keyword is not empty, search book and download
//	fmt.Printf(introduction)
//	if _, err := fmt.Scanln(&input); err == nil {
//		return input
//	} else {
//		fmt.Println(err)
//		InputInt(introduction)
//	}
//	fmt.Println("something wrong, return 0")
//	return 0
//}

// input str
func InputStr(introduction string) string {
	var input string
	// if search keyword is not empty, search book and download
	fmt.Printf(introduction)
	if _, err := fmt.Scanln(&input); err == nil {
		if input != "" {
			return input
		}
	}
	return InputStr(">")
}
func IsNum(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// StrToInt string to int
func StrToInt(str string) int {
	if i, err := strconv.Atoi(str); err == nil {
		return i
	} else {
		return 0
	}
}

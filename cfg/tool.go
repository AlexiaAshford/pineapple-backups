package cfg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
)

func RegexpName(Name string) string {
	return regexp.MustCompile(`[\\/:*?"<>|]`).ReplaceAllString(Name, "")
}

func TestList(List []string, testString string) bool {
	for _, s := range List {
		if s == testString {
			return true
		}
	}
	return false
}

func ExtractBookID(url string) string {
	if url != "" {
		bookID := regexp.MustCompile(`(\d+)`).FindStringSubmatch(url)
		if len(bookID) > 1 {
			return bookID[1]
		}
	}
	return ""
}

func Mkdir(filePath string) {
	dirPath, _ := os.Getwd()
	if err := os.MkdirAll(path.Join(dirPath, filePath), os.ModePerm); err != nil {
		fmt.Println(err)
	}
}

// InputStr input str
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

func FormatJson(jsonString []byte) {
	var str bytes.Buffer
	if err := json.Indent(&str, jsonString, "", "    "); err == nil {
		fmt.Println(str.String())
	} else {
		log.Fatalln(err)
	}
}
func TestKeyword(Text string, keyword any) bool {
	switch keyword.(type) {
	case string:
		return strings.Contains(Text, keyword.(string))
	case int:
		return strings.Contains(Text, strconv.Itoa(keyword.(int)))
	default:
		panic("keyword type error")
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

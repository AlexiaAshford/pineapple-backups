package cfg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path"
	"regexp"
	"strconv"
	"time"
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

func DelayTime() {
	time.Sleep(time.Second * time.Duration(rand.Intn(2)))
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

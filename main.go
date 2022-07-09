package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sf/setting"
	"sf/src"
)

func downloadBook(bookId string) {
	BookData := src.GetBookDetailed(bookId)
	fmt.Printf("开始下载:%s\n", BookData.NovelName)
	if err := ioutil.WriteFile(fmt.Sprintf("save/%v.txt", BookData.NovelName),
		[]byte(BookData.NovelName), 0777); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	setting.NewFile(fmt.Sprintf("config/%v.json", BookData.NovelName))
	src.GetCatalogue(BookData)
}

func main() {
	setting.NewMyJsonPro()
	if len(os.Args) >= 2 {
		inputs := os.Args[1:]
		switch {
		case inputs[0] == "l", inputs[0] == "login":
			if len(inputs) >= 3 {
				src.LoginAccount(inputs[1], inputs[2])
			} else {
				fmt.Println("parameters are not enough")
			}
		case inputs[0] == "d", inputs[0] == "download":
			if len(inputs) >= 2 {
				downloadBook(inputs[1])
			} else {
				fmt.Println("parameters are not enough")
			}
		}
	} else {
		fmt.Println("please input parameters, like: sf login username password")
	}
}

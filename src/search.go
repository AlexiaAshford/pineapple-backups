package src

import (
	"errors"
	"fmt"
	"sf/cfg"
	"sf/src/boluobao"
)

func SearchDetailed(keyword string) error {
	response := boluobao.GetSearchDetailedByKeyword(keyword)
	if response.Status.HTTPCode != 200 || len(response.Data.Novels) == 0 {
		return errors.New(keyword + "is not found in sfacg！")
	} else {
		fmt.Println("search result length:", len(response.Data.Novels))
	}
	for index, book := range response.Data.Novels {
		fmt.Println("Index:", index, "\t\t\tBookName:", book.NovelName)
		cfg.Vars.BookInfoList = append(cfg.Vars.BookInfoList, InitBookStruct(book))
	}
	return nil
}

func SearchBook(search string) string {
	// if search keyword is not empty, search book and download
	if err := SearchDetailed(search); err == nil {
		inputs := cfg.InputInt("please input the index of the book you want to download:")
		if inputs < len(cfg.Vars.BookInfoList) {
			return cfg.Vars.BookInfoList[inputs].NovelID
		} else {
			fmt.Println("index out of range, please input again")
		}
	} else {
		fmt.Println(err)
	}
	return ""
}

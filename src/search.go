package src

import (
	"errors"
	"fmt"
	"sf/cfg"
	"sf/src/boluobao"
)

func GetSearchDetailed(keyword string) error {
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

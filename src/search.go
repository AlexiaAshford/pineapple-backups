package src

import (
	"errors"
	"fmt"
	"sf/cfg"
	"sf/src/boluobao"
	"sf/structural/sfacg_structs"
)

func AddBookInformation(book sfacg_structs.BookInfoData) {
	cfg.Vars.BookInfoList = append(cfg.Vars.BookInfoList, BookInformation(book, false))
}

func GetSearchDetailed(keyword string) error {
	response := boluobao.GetSearchDetailedByKeyword(keyword)
	if response.Status.HTTPCode != 200 || len(response.Data.Novels) == 0 {
		return errors.New(keyword + "is not found in sfacg！")
	} else {
		fmt.Println("search result length:", len(response.Data.Novels))
	}
	for index, book := range response.Data.Novels {
		fmt.Println("Index:", index, "\t\t\tBookName:", book.NovelName)
		AddBookInformation(book)
	}
	return nil
}

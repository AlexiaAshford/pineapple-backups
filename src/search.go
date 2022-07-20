package src

import (
	"errors"
	"fmt"
	"sf/cfg"
	"sf/src/boluobao"
	"sf/structural"
	"sf/structural/sfacg_structs"
	"strconv"
)

func AddBookInformation(book sfacg_structs.Novels) {
	cfg.Vars.BookInfoList = append(
		cfg.Vars.BookInfoList,
		structural.Books{
			NovelName:  cfg.RegexpName(book.NovelName),
			NovelID:    strconv.Itoa(book.NovelID),
			IsFinish:   book.IsFinish,
			MarkCount:  book.MarkCount,
			NovelCover: book.NovelCover,
			AuthorName: book.AuthorName,
			CharCount:  book.CharCount,
			SignStatus: book.SignStatus,
		},
	)
}

func GetSearchDetailed(keyword string) error {
	response := boluobao.GetSearchDetailedByKeyword(keyword)
	if response.Status.HTTPCode != 200 || len(response.Data.Novels) > 0 {
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

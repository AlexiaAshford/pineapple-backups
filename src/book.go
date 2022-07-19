package src

import (
	"errors"
	"fmt"
	"os"
	cfg "sf/cfg"
	"sf/src/boluobao"
	"strconv"
)

type Books struct {
	NovelName  string
	NovelID    string
	IsFinish   bool
	MarkCount  int
	NovelCover string
	AuthorName string
	CharCount  int
	SignStatus string
}

func GetBookDetailed(bookId string) (Books, error) {
	response := boluobao.GetBookDetailedById(bookId)
	if response.Status.HTTPCode != 200 || response.Data.NovelName == "" {
		return Books{}, errors.New(bookId + " is not a valid book number！")
	} else {
		fmt.Println("BookName:", response.Data.NovelName)
		fmt.Println("BookID:", response.Data.NovelID)
		fmt.Println("AuthorName:", response.Data.AuthorName)
		fmt.Println("CharCount:", response.Data.CharCount)
		fmt.Println("MarkCount:", response.Data.MarkCount)
	}
	return Books{
		NovelName:  cfg.RegexpName(response.Data.NovelName),
		NovelID:    strconv.Itoa(response.Data.NovelID),
		IsFinish:   response.Data.IsFinish,
		MarkCount:  response.Data.MarkCount,
		NovelCover: response.Data.NovelCover,
		AuthorName: response.Data.AuthorName,
		CharCount:  response.Data.CharCount,
		SignStatus: response.Data.SignStatus,
	}, nil

}

func GetSearchDetailed(keyword string) []Books {
	var searchList []Books
	response := boluobao.GetSearchDetailedByKeyword(keyword)
	if response.Status.HTTPCode != 200 || len(response.Data.Novels) == 0 {
		fmt.Println(keyword + "is not a valid book number！")
		os.Exit(0)
	}
	fmt.Println("search result length:", len(response.Data.Novels))
	for index, bookInfo := range response.Data.Novels {
		fmt.Println("Index:", index, "\t\t\tBookName:", bookInfo.NovelName)
		searchList = append(
			searchList, Books{
				NovelName:  cfg.RegexpName(bookInfo.NovelName),
				NovelID:    strconv.Itoa(bookInfo.NovelID),
				IsFinish:   bookInfo.IsFinish,
				MarkCount:  bookInfo.MarkCount,
				NovelCover: bookInfo.NovelCover,
				AuthorName: bookInfo.AuthorName,
				CharCount:  bookInfo.CharCount,
				SignStatus: bookInfo.SignStatus,
			},
		)
	}
	return searchList
}

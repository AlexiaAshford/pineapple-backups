package src

import (
	"fmt"
	"sf/src/boluobao"
	cfg "sf/src/config"
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

func GetBookDetailed(bookId string) Books {
	response := boluobao.Get_book_detailed_by_id(bookId)
	if response.Status.HTTPCode != 200 || response.Data.NovelName == "" {
		panic(bookId + "is not a valid book number！")
	} else {
		fmt.Println("BookName:", response.Data.NovelName)
		fmt.Println("BookID:", response.Data.NovelID)
		fmt.Println("AuthorName:", response.Data.AuthorName)
		fmt.Println("CharCount:", response.Data.CharCount)
		fmt.Println("MarkCount:", response.Data.MarkCount)
		return Books{
			NovelName:  cfg.RegexpName(response.Data.NovelName),
			NovelID:    strconv.Itoa(response.Data.NovelID),
			IsFinish:   response.Data.IsFinish,
			MarkCount:  response.Data.MarkCount,
			NovelCover: response.Data.NovelCover,
			AuthorName: response.Data.AuthorName,
			CharCount:  response.Data.CharCount,
			SignStatus: response.Data.SignStatus,
		}
	}
}

func GetSearchDetailed(keyword string) []Books {
	var searchList []Books
	response := boluobao.Get_search_detailed_by_keyword(keyword)
	if response.Status.HTTPCode != 200 || len(response.Data.Novels) == 0 {
		panic(keyword + "is not a valid book number！")
	}
	fmt.Println("search result length:", len(response.Data.Novels))
	for index, bookInfo := range response.Data.Novels {
		fmt.Println("Index:", index, "\t\t\tBookName:", bookInfo.NovelName)
		searchList = append(searchList, Books{
			NovelName:  cfg.RegexpName(bookInfo.NovelName),
			NovelID:    strconv.Itoa(bookInfo.NovelID),
			IsFinish:   bookInfo.IsFinish,
			MarkCount:  bookInfo.MarkCount,
			NovelCover: bookInfo.NovelCover,
			AuthorName: bookInfo.AuthorName,
			CharCount:  bookInfo.CharCount,
			SignStatus: bookInfo.SignStatus,
		})
	}
	return searchList
}

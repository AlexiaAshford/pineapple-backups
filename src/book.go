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

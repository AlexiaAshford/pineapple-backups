package src

import (
	"encoding/json"
	"fmt"
	cfg "sf/setting"
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
	var BookData BookInformation
	if err := json.Unmarshal(Get(fmt.Sprintf("novels/%v?expand=", bookId)), &BookData); err != nil {
		panic(err)
	}
	if BookData.Status.HTTPCode != 200 || BookData.Data.NovelName == "" {
		panic(bookId + "is not a valid book number！")
	}
	return Books{
		NovelName:  cfg.RegexpName(BookData.Data.NovelName),
		NovelID:    strconv.Itoa(BookData.Data.NovelID),
		IsFinish:   BookData.Data.IsFinish,
		MarkCount:  BookData.Data.MarkCount,
		NovelCover: BookData.Data.NovelCover,
		AuthorName: BookData.Data.AuthorName,
		CharCount:  BookData.Data.CharCount,
		SignStatus: BookData.Data.SignStatus,
	}
}

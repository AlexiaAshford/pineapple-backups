package src

import (
	"errors"
	"fmt"
	"os"
	"sf/cfg"
	"sf/multi"
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

func SfacgBookInit(bookID string, Index int, Locks *multi.GoLimit) {
	if Locks != nil {
		defer Locks.Done() // finish this goroutine when this function return
	}
	if BookData, err := GetSfacgBookDetailed(bookID); err == nil { // get book data by book id
		fmt.Println("BookName:", BookData.NovelName)
		fmt.Println("BookID:", BookData.NovelID)
		fmt.Println("AuthorName:", BookData.AuthorName)
		fmt.Println("CharCount:", BookData.CharCount)
		fmt.Println("MarkCount:", BookData.MarkCount)
		fmt.Printf("开始下载:%s\n", BookData.NovelName)
		cachepath := fmt.Sprintf("%v/%v.txt", cfg.Vars.SaveFile, BookData.NovelName)
		for i := 0; i < 5; i++ {
			if cfg.WriteFile(cachepath, BookData.NovelName+"\n", 0644) == nil {
				break
			} else {
				fmt.Println("write file error, try again...")
			}
		}
		if SfacgCatalogue(BookData) {
			if Index > 0 {
				fmt.Printf("\nIndex:%v\t\tNovelName:%vdownload complete!", Index, BookData.NovelName)
			} else {
				fmt.Printf("\nNovelName:%vdownload complete!", BookData.NovelName)
			}
		}
	} else {
		fmt.Println("Error:", err)
	}
}
func GetSfacgBookDetailed(bookId string) (Books, error) {
	response := boluobao.GetBookDetailedById(bookId)
	if response.Status.HTTPCode != 200 || response.Data.NovelName == "" {
		return Books{}, errors.New(bookId + " is not a valid book number！")
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
	} else {
		fmt.Println("search result length:", len(response.Data.Novels))
		for index, bookInfo := range response.Data.Novels {
			fmt.Println("Index:", index, "\t\t\tBookName:", bookInfo.NovelName)
			searchList = append(
				searchList,
				Books{
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
	}
	return searchList
}

package src

import (
	"errors"
	"fmt"
	"os"
	"sf/cfg"
	"sf/multi"
	"sf/src/boluobao"
	"sf/structural"
	"strconv"
)

func SfacgBookInit(bookID string, Index int, Locks *multi.GoLimit) {
	if Locks != nil {
		defer Locks.Done() // finish this goroutine when this function return
	}
	if err := GetSfacgBookDetailed(bookID); err == nil { // get book data by book id
		fmt.Println("BookName:", cfg.Vars.BookInfo.NovelName)
		fmt.Println("BookID:", cfg.Vars.BookInfo.NovelID)
		fmt.Println("AuthorName:", cfg.Vars.BookInfo.AuthorName)
		fmt.Println("CharCount:", cfg.Vars.BookInfo.CharCount)
		fmt.Println("MarkCount:", cfg.Vars.BookInfo.MarkCount)
		fmt.Printf("开始下载:%s\n", cfg.Vars.BookInfo.NovelName)
		cachepath := fmt.Sprintf("%v/%v.txt", cfg.Vars.SaveFile, cfg.Vars.BookInfo.NovelName)
		for i := 0; i < 5; i++ {
			if cfg.WriteFile(cachepath, cfg.Vars.BookInfo.NovelName+"\n", 0644) == nil {
				break
			} else {
				fmt.Println("write file error, try again...")
			}
		}
		if SfacgCatalogue() {
			if Index > 0 {
				fmt.Printf("\nIndex:%v\t\tNovelName:%vdownload complete!", Index, cfg.Vars.BookInfo.NovelName)
			} else {
				fmt.Printf("\nNovelName:%vdownload complete!", cfg.Vars.BookInfo.NovelName)
			}
		}
	} else {
		fmt.Println("Error:", err)
	}
}
func GetSfacgBookDetailed(bookId string) error {
	response := boluobao.GetBookDetailedById(bookId)
	if response.Status.HTTPCode == 200 && response.Data.NovelName != "" {
		cfg.Vars.BookInfo = structural.Books{
			NovelName:  cfg.RegexpName(response.Data.NovelName),
			NovelID:    strconv.Itoa(response.Data.NovelID),
			IsFinish:   response.Data.IsFinish,
			MarkCount:  response.Data.MarkCount,
			NovelCover: response.Data.NovelCover,
			AuthorName: response.Data.AuthorName,
			CharCount:  response.Data.CharCount,
			SignStatus: response.Data.SignStatus,
		}
		return nil
	} else {
		return errors.New(bookId + " is not a valid book number！")
	}

}

func GetSearchDetailed(keyword string) {
	response := boluobao.GetSearchDetailedByKeyword(keyword)
	if response.Status.HTTPCode != 200 || len(response.Data.Novels) == 0 {
		fmt.Println(keyword + "is not a valid book number！")
		os.Exit(0)
	} else {
		fmt.Println("search result length:", len(response.Data.Novels))
		for index, bookInfo := range response.Data.Novels {
			fmt.Println("Index:", index, "\t\t\tBookName:", bookInfo.NovelName)
			cfg.Vars.BookInfoList = append(
				cfg.Vars.BookInfoList,
				structural.Books{
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
}

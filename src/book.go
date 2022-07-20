package src

import (
	"errors"
	"fmt"
	"sf/cfg"
	"sf/multi"
	"sf/src/boluobao"
	"sf/structural"
	"sf/structural/sfacg_structs"
	"strconv"
)

func SfacgBookInit(bookID string, Index int, Locks *multi.GoLimit) {
	if Locks != nil {
		defer Locks.Done() // finish this goroutine when this function return
	}
	if err := GetSfacgBookDetailed(bookID); err == nil { // get book data by book id
		cachepath := fmt.Sprintf("%v/%v.txt", cfg.Vars.SaveFile, cfg.Vars.BookInfo.NovelName)
		for i := 0; i < 5; i++ {
			if cfg.WriteFile(cachepath, cfg.Vars.BookInfo.NovelName+"\n", 0644) == nil {
				break
			} else {
				fmt.Println("write file error, try again...")
			}
		}
		catalogues := sfacgCatalogue{}
		if catalogues.SfacgCatalogue() {
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

func BookInformation(book sfacg_structs.BookInfoData, ShowDetailed bool) structural.Books {
	if ShowDetailed {
		briefIntroduction := fmt.Sprintf(
			"Name: %v\nBookID: %v\nAuthor: %v\nCount: %v\nMark: %v\n",
			book.NovelName, book.NovelID, book.AuthorName, book.CharCount, book.MarkCount)
		fmt.Println(briefIntroduction)
	}
	return structural.Books{
		NovelName:  cfg.RegexpName(book.NovelName),
		NovelID:    strconv.Itoa(book.NovelID),
		IsFinish:   book.IsFinish,
		MarkCount:  book.MarkCount,
		NovelCover: book.NovelCover,
		AuthorName: book.AuthorName,
		CharCount:  book.CharCount,
		SignStatus: book.SignStatus,
	}
}
func GetSfacgBookDetailed(bookId string) error {
	response := boluobao.GetBookDetailedById(bookId)
	if response.Status.HTTPCode == 200 && response.Data.NovelName != "" {
		cfg.Vars.BookInfo = BookInformation(response.Data, true)
		return nil
	} else {
		return errors.New(bookId + " is not a valid book number！")
	}

}

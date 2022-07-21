package src

import (
	"fmt"
	"sf/cfg"
	"sf/multi"
	"sf/src/boluobao"
	"sf/structural"
	"sf/structural/hbooker_structs"
	"sf/structural/sfacg_structs"
	"strconv"
)

type BookInits struct {
	BookID string
	Index  int
	Locks  *multi.GoLimit
}

func (books *BookInits) SfacgBookInit() {
	if books.Locks != nil {
		defer books.Locks.Done() // finish this goroutine when this function return
	}
	response := boluobao.GetBookDetailedById(books.BookID)
	if response.Status.HTTPCode == 200 && response.Data.NovelName != "" {
		cfg.Vars.BookInfo = InitBookStruct(response.Data)
		books.ShowBookDetailed()
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
			if books.Index > 0 {
				fmt.Printf("\nIndex:%v\t\tNovelName:%vdownload complete!", books.Index, cfg.Vars.BookInfo.NovelName)
			} else {
				fmt.Printf("\nNovelName:%vdownload complete!", cfg.Vars.BookInfo.NovelName)
			}
		}
	} else {
		fmt.Println(books.BookID, "is not a valid book number！")
	}
}
func (books *BookInits) ShowBookDetailed() {
	briefIntroduction := fmt.Sprintf(
		"Name: %v\nBookID: %v\nAuthor: %v\nCount: %v\nMark: %v\n",
		cfg.Vars.BookInfo.NovelName, cfg.Vars.BookInfo.NovelID,
		cfg.Vars.BookInfo.AuthorName, cfg.Vars.BookInfo.CharCount,
		cfg.Vars.BookInfo.MarkCount,
	)
	fmt.Println(briefIntroduction)
}

func InitBookStruct(bookAny any) structural.Books {
	switch bookAny.(type) {
	case sfacg_structs.BookInfoData:
		book := bookAny.(sfacg_structs.BookInfoData)
		return structural.Books{
			NovelName:  cfg.RegexpName(book.NovelName),
			NovelID:    strconv.Itoa(book.NovelID),
			IsFinish:   book.IsFinish,
			MarkCount:  strconv.Itoa(book.MarkCount),
			NovelCover: book.NovelCover,
			AuthorName: book.AuthorName,
			CharCount:  strconv.Itoa(book.CharCount),
			SignStatus: book.SignStatus,
		}
	case hbooker_structs.BookInfo:
		book := bookAny.(hbooker_structs.BookInfo)
		return structural.Books{
			NovelName:  cfg.RegexpName(book.BookName),
			NovelID:    book.BookID,
			NovelCover: book.Cover,
			AuthorName: book.AuthorName,
			CharCount:  book.TotalWordCount,
			SignStatus: book.Discount,
		}
	}
	return structural.Books{}
}

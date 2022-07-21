package src

import (
	"fmt"
	"sf/cfg"
	"sf/multi"
	"sf/src/boluobao"
	HbookerAPI "sf/src/hbooker"
	"sf/structural"
	"sf/structural/hbooker_structs"
	"sf/structural/sfacg_structs"
	"strconv"
)

type BookInits struct {
	BookID        string
	Index         int
	ShowBook      bool
	Locks         *multi.GoLimit
	SfacgBookData sfacg_structs.BookInfoData
	CatBookData   hbooker_structs.BookInfo
}

func (books *BookInits) CatBookInit() {
	if books.Locks != nil {
		defer books.Locks.Done() // finish this goroutine when this function return
	}
	response := HbookerAPI.GetBookDetailById(books.BookID)
	fmt.Println(response)
}
func (books *BookInits) SfacgBookInit() {
	if books.Locks != nil {
		defer books.Locks.Done() // finish this goroutine when this function return
	}
	response := boluobao.GetBookDetailedById(books.BookID)
	if response.Status.HTTPCode == 200 && response.Data.NovelName != "" {
		books.SfacgBookData = response.Data
		cfg.Vars.BookInfo = books.InitBookStruct()
		if books.ShowBook {
			books.ShowBookDetailed()
		}
		cfg.EncapsulationWrite(
			fmt.Sprintf("%v/%v.txt", cfg.Vars.SaveFile, cfg.Vars.BookInfo.NovelName),
			cfg.Vars.BookInfo.NovelName+"\n", 5, 0644,
		) // write novel name to file for later use
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

func (books *BookInits) InitBookStruct() structural.Books {
	if cfg.Vars.AppType == "sfacg" {
		return structural.Books{
			NovelName:  cfg.RegexpName(books.SfacgBookData.NovelName),
			NovelID:    strconv.Itoa(books.SfacgBookData.NovelID),
			IsFinish:   books.SfacgBookData.IsFinish,
			MarkCount:  strconv.Itoa(books.SfacgBookData.MarkCount),
			NovelCover: books.SfacgBookData.NovelCover,
			AuthorName: books.SfacgBookData.AuthorName,
			CharCount:  strconv.Itoa(books.SfacgBookData.CharCount),
			SignStatus: books.SfacgBookData.SignStatus,
		}
	}
	if cfg.Vars.AppType == "sfacg" {
		return structural.Books{
			NovelName:  cfg.RegexpName(books.CatBookData.BookName),
			NovelID:    books.CatBookData.BookID,
			NovelCover: books.CatBookData.Cover,
			AuthorName: books.CatBookData.AuthorName,
			CharCount:  books.CatBookData.TotalWordCount,
			//SignStatus: books.CatBookData.SignStatus,
		}
	}
	return structural.Books{}
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

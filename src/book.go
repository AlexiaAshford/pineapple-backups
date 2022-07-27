package src

import (
	"fmt"
	"sf/cfg"
	"sf/multi"
	"sf/src/boluobao"
	"sf/src/hbooker"
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

func (books *BookInits) DownloadBookInit() Catalogue {
	if cfg.Vars.AppType == "sfacg" {
		response := boluobao.GetBookDetailedById(books.BookID)
		if response.Status.HTTPCode == 200 && response.Data.NovelName != "" {
			books.SfacgBookData = response.Data
		} else {
			fmt.Println(books.BookID, "is not a valid book number！\n", response.Status.Msg)
		}
	} else if cfg.Vars.AppType == "cat" {
		response := hbooker.GetBookDetailById(books.BookID)
		if response.Code == "100000" {
			books.CatBookData = response.Data.BookInfo
		} else {
			fmt.Println(books.BookID, "is not a valid book number！")
		}
	} else {
		panic("app type is not valid!")
	}

	cfg.BookConfig.BookInfo = books.InitBookStruct()
	savePath := fmt.Sprintf("%v/%v.txt", cfg.Vars.SaveFile, cfg.BookConfig.BookInfo.NovelName)
	if !cfg.CheckFileExist(savePath) {
		cfg.EncapsulationWrite(savePath, books.ShowBookDetailed()+"\n\n", 5, "w")
	} else {
		books.ShowBookDetailed()
	}
	return Catalogue{}
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
	if cfg.Vars.AppType == "cat" {
		return structural.Books{
			NovelName:  cfg.RegexpName(books.CatBookData.BookName),
			NovelID:    books.CatBookData.BookID,
			NovelCover: books.CatBookData.Cover,
			AuthorName: books.CatBookData.AuthorName,
			CharCount:  books.CatBookData.TotalWordCount,
			MarkCount:  books.CatBookData.UpdateStatus,
			//SignStatus: books.CatBookData.SignStatus,
		}
	}
	return structural.Books{}
}

func (books *BookInits) ShowBookDetailed() string {
	briefIntroduction := fmt.Sprintf(
		"Name: %v\nBookID: %v\nAuthor: %v\nCount: %v\nMark: %v\n",
		cfg.BookConfig.BookInfo.NovelName, cfg.BookConfig.BookInfo.NovelID,
		cfg.BookConfig.BookInfo.AuthorName, cfg.BookConfig.BookInfo.CharCount,
		cfg.BookConfig.BookInfo.MarkCount,
	)
	if books.ShowBook {
		fmt.Println(briefIntroduction)
	}
	return briefIntroduction
}

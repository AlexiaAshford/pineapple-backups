package src

import (
	"fmt"
	"path"
	"sf/cfg"
	"sf/src/boluobao"
	"sf/src/hbooker"
	"sf/struct"
	"sf/struct/hbooker_structs"
	"sf/struct/sfacg_structs"
	"strconv"
)

type BookInits struct {
	BookID   string
	Index    int
	ShowBook bool
	Locks    *cfg.GoLimit
	BookData any
}

func (books *BookInits) DownloadBookInit() Catalogue {
	if cfg.Vars.AppType == "sfacg" {
		response := boluobao.GetBookDetailedById(books.BookID)
		if response.Status.HTTPCode == 200 && response.Data.NovelName != "" {
			books.BookData = response
		} else {
			fmt.Println(books.BookID, "is not a valid book number！\nmessage:", response.Status.Msg)
			return Catalogue{TestBookResult: false}
		}
	} else if cfg.Vars.AppType == "cat" {
		response := hbooker.GetBookDetailById(books.BookID)
		if response.Code == "100000" {
			books.BookData = response.Data.BookInfo
		} else {
			fmt.Println(books.BookID, "is not a valid book number！")
		}
	} else {
		panic("app type" + cfg.Vars.AppType + " is not valid!")
	}
	cfg.CurrentBook.BookInfo = books.InitBookStruct()

	savePath := path.Join(cfg.Vars.SaveFile, cfg.CurrentBook.BookInfo.NovelName+".txt")
	if !cfg.Exist(savePath) {
		cfg.Write(savePath, books.ShowBookDetailed()+"\n\n", "w")
	} else {
		books.ShowBookDetailed()
	}
	return Catalogue{SaveTextPath: savePath, TestBookResult: true}

}
func (books *BookInits) InitBookStruct() _struct.Books {
	switch books.BookData.(type) {
	case *sfacg_structs.BookInfo:
		result := books.BookData.(*sfacg_structs.BookInfo).Data
		return _struct.Books{
			NovelName:  cfg.RegexpName(result.NovelName),
			NovelID:    strconv.Itoa(result.NovelID),
			IsFinish:   result.IsFinish,
			MarkCount:  strconv.Itoa(result.MarkCount),
			NovelCover: result.NovelCover,
			AuthorName: result.AuthorName,
			CharCount:  strconv.Itoa(result.CharCount),
			SignStatus: result.SignStatus,
		}
	case hbooker_structs.BookInfo:
		result := books.BookData.(hbooker_structs.BookInfo)
		return _struct.Books{
			NovelName:  cfg.RegexpName(result.BookName),
			NovelID:    result.BookID,
			NovelCover: result.Cover,
			AuthorName: result.AuthorName,
			CharCount:  result.TotalWordCount,
			MarkCount:  result.UpdateStatus,
			//SignStatus: result.SignStatus,
		}
	}
	return _struct.Books{}
}

func (books *BookInits) ShowBookDetailed() string {
	briefIntroduction := fmt.Sprintf(
		"Name: %v\nBookID: %v\nAuthor: %v\nCount: %v\nMark: %v\n",
		cfg.CurrentBook.BookInfo.NovelName, cfg.CurrentBook.BookInfo.NovelID,
		cfg.CurrentBook.BookInfo.AuthorName, cfg.CurrentBook.BookInfo.CharCount,
		cfg.CurrentBook.BookInfo.MarkCount,
	)
	if books.ShowBook {
		fmt.Println(briefIntroduction)
	}
	return briefIntroduction
}

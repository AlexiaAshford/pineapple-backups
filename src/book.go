package src

import (
	"fmt"
	"os"
	"path"
	"sf/cfg"
	"sf/epub"
	"sf/src/boluobao"
	"sf/src/hbooker"
	"sf/src/https"
	"sf/struct"
	"sf/struct/hbooker_structs"
	"sf/struct/sfacg_structs"
	"strconv"
	"strings"
)

type BookInits struct {
	BookID      string
	Index       int
	ShowBook    bool
	Locks       *cfg.GoLimit
	BookData    any
	EpubSetting *epub.Epub
}

func (books *BookInits) InitEpubFile() {
	// set epub setting and add section
	AddImage := true
	books.EpubSetting = epub.NewEpub(cfg.Current.BookInfo.NovelName)
	books.EpubSetting.SetAuthor(cfg.Current.BookInfo.AuthorName) // set author
	coverPath := path.Join("cover", cfg.Current.BookInfo.NovelName+".jpg")
	if !cfg.Exist(coverPath) {
		if reader := https.GetCover(cfg.Current.BookInfo.NovelCover); reader == nil {
			fmt.Println("download cover failed!")
			AddImage = false
		} else {
			_ = os.WriteFile(coverPath, reader, 0666)
		}
	}
	if AddImage {
		_, _ = books.EpubSetting.AddImage(coverPath, "")
		books.EpubSetting.SetCover(strings.Replace(coverPath, "cover", "../images", -1), "")
	}

}

func (books *BookInits) DownloadBookInit() Catalogue {
	switch cfg.Vars.AppType {
	case "sfacg":
		response := boluobao.GetBookDetailedById(books.BookID)
		if response.Status.HTTPCode == 200 && response.Data.NovelName != "" {
			books.BookData = response
		} else {
			fmt.Println(books.BookID, "is not a valid book number！\nmessage:", response.Status.Msg)
			return Catalogue{TestBookResult: false}
		}
	case "cat":
		response := hbooker.GetBookDetailById(books.BookID)
		if response.Code == "100000" {
			books.BookData = response.Data.BookInfo
		} else {
			fmt.Println(books.BookID, "is not a valid book number！")
			return Catalogue{TestBookResult: false}
		}
	default:
		panic("app type" + cfg.Vars.AppType + " is not valid!")

	}
	cfg.Current.BookInfo = books.InitBookStruct()
	cfg.Current.ConfigPath = path.Join(cfg.Vars.ConfigName, cfg.Current.BookInfo.NovelName+".conf")
	cfg.Current.OutputPath = path.Join(cfg.Vars.OutputName, cfg.Current.BookInfo.NovelName+".txt")
	books.InitEpubFile()
	if !cfg.Exist(cfg.Current.OutputPath) {
		cfg.Write(cfg.Current.OutputPath, books.ShowBookDetailed()+"\n\n", "w")
	} else {
		books.ShowBookDetailed()
	}
	return Catalogue{TestBookResult: true, EpubSetting: books.EpubSetting}

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
		cfg.Current.BookInfo.NovelName, cfg.Current.BookInfo.NovelID,
		cfg.Current.BookInfo.AuthorName, cfg.Current.BookInfo.CharCount,
		cfg.Current.BookInfo.MarkCount,
	)
	if books.ShowBook {
		fmt.Println(briefIntroduction)
	}
	return briefIntroduction
}

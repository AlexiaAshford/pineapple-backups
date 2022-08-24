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
	"strconv"
	"strings"
)

type BookInits struct {
	BookID      string
	Index       int
	ShowBook    bool
	Locks       *cfg.GoLimit
	EpubSetting *epub.Epub
}

func (books *BookInits) InitEpubFile() {
	AddImage := true                                             // add image to epub file
	books.EpubSetting = epub.NewEpub(cfg.Current.Book.NovelName) // set epub setting and add section
	books.EpubSetting.SetAuthor(cfg.Current.Book.AuthorName)     // set author
	if !cfg.Exist(cfg.Current.CoverPath) {
		if reader := https.GetCover(cfg.Current.Book.NovelCover); reader == nil {
			fmt.Println("download cover failed!")
			AddImage = false
		} else {
			_ = os.WriteFile(cfg.Current.CoverPath, reader, 0666)
		}
	}
	if AddImage {
		_, _ = books.EpubSetting.AddImage(cfg.Current.CoverPath, "")
		books.EpubSetting.SetCover(strings.Replace(cfg.Current.CoverPath, "cover", "../images", -1), "")
	}

}

func (books *BookInits) SetBookInfo() Catalogue {
	switch cfg.Vars.AppType {
	case "sfacg":
		response := boluobao.GET_BOOK_INFORMATION(books.BookID)
		if response.Status.HTTPCode == 200 && response.Data.NovelName != "" {
			cfg.Current.Book = _struct.Books{
				NovelName:  cfg.RegexpName(response.Data.NovelName),
				NovelID:    strconv.Itoa(response.Data.NovelID),
				NovelCover: response.Data.NovelCover,
				AuthorName: response.Data.AuthorName,
				CharCount:  strconv.Itoa(response.Data.CharCount),
				MarkCount:  strconv.Itoa(response.Data.MarkCount),
				SignStatus: response.Data.SignStatus,
			}
		} else {
			fmt.Println(books.BookID, "is not a valid book number！\nmessage:", response.Status.Msg)
			return Catalogue{TestBookResult: false}
		}
	case "cat":
		response := hbooker.GET_BOOK_INFORMATION(books.BookID)
		if response.Code == "100000" {
			cfg.Current.Book = _struct.Books{
				NovelName:  cfg.RegexpName(response.Data.BookInfo.BookName),
				NovelID:    response.Data.BookInfo.BookID,
				NovelCover: response.Data.BookInfo.Cover,
				AuthorName: response.Data.BookInfo.AuthorName,
				CharCount:  response.Data.BookInfo.TotalWordCount,
				MarkCount:  response.Data.BookInfo.UpdateStatus,
				//SignStatus: result.SignStatus,
			}
		} else {
			fmt.Println(books.BookID, "is not a valid book number！")
			return Catalogue{TestBookResult: false}
		}
	default:
		panic("app type" + cfg.Vars.AppType + " is not valid!")
	}
	cfg.Current.ConfigPath = path.Join(cfg.Vars.ConfigName, cfg.Current.Book.NovelName)
	cfg.Current.OutputPath = path.Join(cfg.Vars.OutputName, cfg.Current.Book.NovelName+".txt")
	cfg.Current.CoverPath = path.Join("cover", cfg.Current.Book.NovelName+".jpg")
	books.InitEpubFile()
	cfg.Write(cfg.Current.OutputPath, books.ShowBookDetailed()+"\n\n", "w")
	return Catalogue{TestBookResult: true, EpubSetting: books.EpubSetting}

}

func (books *BookInits) ShowBookDetailed() string {
	briefIntroduction := fmt.Sprintf(
		"Name: %v\nBookID: %v\nAuthor: %v\nCount: %v\n",
		cfg.Current.Book.NovelName, cfg.Current.Book.NovelID,
		cfg.Current.Book.AuthorName, cfg.Current.Book.CharCount,
	)
	if books.ShowBook {
		fmt.Println(briefIntroduction)
	}
	return briefIntroduction
}

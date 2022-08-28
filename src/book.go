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
	_struct "sf/struct"
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
	var err error
	var result _struct.Books
	switch cfg.Vars.AppType {
	case "sfacg":
		result, err = boluobao.GET_BOOK_INFORMATION(books.BookID)
	case "cat":
		result, err = hbooker.GET_BOOK_INFORMATION(books.BookID)
	}

	if err == nil {
		cfg.Current.Book = result
		return books.BookDetailed()
	} else {
		fmt.Println(books.BookID, "is not a valid book number！\nmessage:", err)
		return Catalogue{TestBookResult: false}
	}
}

func (books *BookInits) BookDetailed() Catalogue {
	cfg.Current.ConfigPath = path.Join(cfg.Vars.ConfigName, cfg.Current.Book.NovelName)
	cfg.Current.OutputPath = path.Join(cfg.Vars.OutputName, cfg.Current.Book.NovelName+".txt")
	cfg.Current.CoverPath = path.Join("cover", cfg.Current.Book.NovelName+".jpg")
	books.InitEpubFile()
	briefIntroduction := fmt.Sprintf("Name: %v\nBookID: %v\nAuthor: %v\nCount: %v\n\n\n",
		cfg.Current.Book.NovelName, cfg.Current.Book.NovelID, cfg.Current.Book.AuthorName, cfg.Current.Book.CharCount,
	)
	if books.ShowBook {
		fmt.Println(briefIntroduction)
	}
	cfg.Write(cfg.Current.OutputPath, briefIntroduction, "w")
	return Catalogue{TestBookResult: true, EpubSetting: books.EpubSetting}
}

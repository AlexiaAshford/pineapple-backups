package src

import (
	"fmt"
	"github.com/VeronicaAlexia/pineapple-backups/cfg"
	"github.com/VeronicaAlexia/pineapple-backups/epub"
	"github.com/VeronicaAlexia/pineapple-backups/src/boluobao"
	"github.com/VeronicaAlexia/pineapple-backups/src/hbooker"
	"github.com/VeronicaAlexia/pineapple-backups/src/https"
	_struct "github.com/VeronicaAlexia/pineapple-backups/struct"
	"os"
	"path"
	"strings"
)

type BookInits struct {
	BookID      string
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
		books.EpubSetting.SetCover(strings.ReplaceAll(cfg.Current.CoverPath, "cover", "../images"), "")
	}

}

func SettingBooks(book_id string) Catalogue {
	var err error
	var result _struct.Books
	switch cfg.Vars.AppType {
	case "sfacg":
		result, err = boluobao.GET_BOOK_INFORMATION(book_id)
	case "cat":
		result, err = hbooker.GET_BOOK_INFORMATION(book_id)
	}
	if err == nil {
		cfg.Current.Book = result
		cfg.Current.ConfigPath = path.Join(cfg.Vars.ConfigName, cfg.Current.Book.NovelName)
		cfg.Current.OutputPath = path.Join(cfg.Vars.OutputName, cfg.Current.Book.NovelName+".txt")
		cfg.Current.CoverPath = path.Join("cover", cfg.Current.Book.NovelName+".jpg")
		books := BookInits{BookID: book_id, Locks: nil, ShowBook: true}
		return books.BookDetailed()
	} else {
		return Catalogue{Test: false, BookMessage: fmt.Sprintf("book_id:%v is invalid:%v", book_id, err)}
	}

}

func (books *BookInits) BookDetailed() Catalogue {
	books.InitEpubFile()
	briefIntroduction := fmt.Sprintf("Name: %v\nBookID: %v\nAuthor: %v\nCount: %v\n\n\n",
		cfg.Current.Book.NovelName, cfg.Current.Book.NovelID, cfg.Current.Book.AuthorName, cfg.Current.Book.CharCount,
	)
	if books.ShowBook {
		fmt.Println(briefIntroduction)
	}
	cfg.Write(cfg.Current.OutputPath, briefIntroduction, "w")
	return Catalogue{Test: true, EpubSetting: books.EpubSetting}
}

package src

import (
	"fmt"
	"github.com/VeronicaAlexia/BoluobaoAPI/boluobao"
	"github.com/VeronicaAlexia/ciweimaoapiLib"
	"github.com/VeronicaAlexia/pineapple-backups/config"
	"github.com/VeronicaAlexia/pineapple-backups/pkg/command"
	"github.com/VeronicaAlexia/pineapple-backups/pkg/epub"
	"github.com/VeronicaAlexia/pineapple-backups/pkg/file"
	"github.com/VeronicaAlexia/pineapple-backups/pkg/request"
	"github.com/VeronicaAlexia/pineapple-backups/pkg/threading"
	"github.com/VeronicaAlexia/pineapple-backups/pkg/tools"
	"os"
	"path"
	"strconv"
)

type BookInits struct {
	BookID      string
	ShowBook    bool
	Locks       *threading.GoLimit
	EpubSetting *epub.Epub
}
type Books struct {
	NovelName  string
	NovelID    string
	IsFinish   bool
	MarkCount  string
	NovelCover string
	AuthorName string
	CharCount  string
	SignStatus string
}

func (books *BookInits) InitEpubFile() {
	AddImage := true                                                        // add image to epub file
	books.EpubSetting = epub.NewEpub(config.Current.NewBooks["novel_name"]) // set epub setting and add section
	books.EpubSetting.SetAuthor(config.Current.NewBooks["author_name"])     // set author
	if !config.Exist(config.Current.CoverPath) {
		if reader := request.Request(config.Current.NewBooks["novel_cover"]); reader == nil {
			fmt.Println("download cover failed!")
			AddImage = false
		} else {
			_ = os.WriteFile(config.Current.CoverPath, reader, 0666)
		}
	}
	if AddImage {
		_, _ = books.EpubSetting.AddImage(config.Current.CoverPath, "")
		books.EpubSetting.SetCover(path.Join("../images", config.Current.NewBooks["novel_name"]+".jpg"), "")
	}

}

func SettingBooks(book_id string) Catalogue {
	var err error
	switch command.Command.AppType {
	case "sfacg":
		BookInfo := boluobao.API.Book.NovelInfo(book_id)
		if BookInfo != nil {
			config.Current.NewBooks = map[string]string{
				"novel_name":  tools.RegexpName(BookInfo.NovelName),
				"novel_id":    strconv.Itoa(BookInfo.NovelID),
				"novel_cover": BookInfo.NovelCover,
				"author_name": BookInfo.AuthorName,
				"char_count":  strconv.Itoa(BookInfo.CharCount),
				"mark_count":  strconv.Itoa(BookInfo.MarkCount),
			}
		} else {
			return Catalogue{Test: false, BookMessage: fmt.Sprintf("book_id:%v is invalid:%v", book_id, err)}
		}
	case "cat":
		BookInfo := ciweimaoapi.GetBookInfo(book_id)
		config.Current.NewBooks = map[string]string{
			"novel_name":  BookInfo.Data.BookInfo.BookName,
			"novel_id":    BookInfo.Data.BookInfo.BookId,
			"novel_cover": BookInfo.Data.BookInfo.Cover,
			"author_name": BookInfo.Data.BookInfo.AuthorName,
			"char_count":  BookInfo.Data.BookInfo.TotalWordCount,
		}
		if BookInfo.Data.BookInfo.BookName == "" {
			return Catalogue{Test: false, BookMessage: fmt.Sprintf("book_id:%v is invalid:%v", book_id, err)}
		}
	}
	//fmt.Println(config.Current.NewBooks)
	tools.Mkdir(path.Join(config.Vars.OutputName, config.Current.NewBooks["novel_name"]))
	config.Current.ConfigPath = path.Join(config.Vars.ConfigName, config.Current.NewBooks["novel_name"])
	//config.Current.OutputPath = path.Join(OutputPath, config.Current.NewBooks["novel_name"]+".txt")
	config.Current.CoverPath = path.Join(config.Vars.ConfigName, config.Vars.CoverFile, config.Current.NewBooks["novel_name"]+".jpg")
	books := BookInits{BookID: book_id, Locks: nil, ShowBook: true}
	return books.BookDetailed()

}

func (books *BookInits) BookDetailed() Catalogue {
	if config.Current.NewBooks["novel_name"] == "" {
		fmt.Println("下载失败")
		return Catalogue{Test: false}
	}
	books.InitEpubFile()
	briefIntroduction := fmt.Sprintf("Name: %v\nBookID: %v\nAuthor: %v\nCount: %v\n\n\n",
		config.Current.NewBooks["novel_name"], config.Current.NewBooks["novel_id"], config.Current.NewBooks["author_name"],
		config.Current.NewBooks["char_count"],
	)
	if books.ShowBook {
		fmt.Println(briefIntroduction)
	}

	file.Open(
		path.Join(config.Vars.OutputName, config.Current.NewBooks["novel_name"], config.Current.NewBooks["novel_name"]+".txt"),
		briefIntroduction,
		"w",
	)
	return Catalogue{Test: true, EpubSetting: books.EpubSetting}
}

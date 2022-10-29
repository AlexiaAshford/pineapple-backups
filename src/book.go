package src

import (
	"encoding/json"
	"fmt"
	"github.com/VeronicaAlexia/pineapple-backups/config"
	"github.com/VeronicaAlexia/pineapple-backups/pkg/epub"
	"github.com/VeronicaAlexia/pineapple-backups/pkg/file"
	"github.com/VeronicaAlexia/pineapple-backups/pkg/request"
	"github.com/VeronicaAlexia/pineapple-backups/pkg/threading"
	"github.com/VeronicaAlexia/pineapple-backups/pkg/tools"
	"github.com/VeronicaAlexia/pineapple-backups/src/app/boluobao"
	"github.com/VeronicaAlexia/pineapple-backups/src/app/hbooker"
	"os"
	"path"
	"strings"
)

type BookInits struct {
	BookID      string
	ShowBook    bool
	Locks       *threading.GoLimit
	EpubSetting *epub.Epub
}

func (books *BookInits) InitEpubFile() {
	AddImage := true                                                // add image to epub file
	books.EpubSetting = epub.NewEpub(config.Current.Book.NovelName) // set epub setting and add section
	books.EpubSetting.SetAuthor(config.Current.Book.AuthorName)     // set author
	if !config.Exist(config.Current.CoverPath) {
		if reader := request.Request(config.Current.Book.NovelCover); reader == nil {
			fmt.Println("download cover failed!")
			AddImage = false
		} else {
			_ = os.WriteFile(config.Current.CoverPath, reader, 0666)
		}
	}
	if AddImage {
		_, _ = books.EpubSetting.AddImage(config.Current.CoverPath, "")
		books.EpubSetting.SetCover(strings.ReplaceAll(config.Current.CoverPath, "cover", "../images"), "")
	}

}

func SettingBooks(book_id string) Catalogue {
	var err error
	config.Current.BackupsPath = path.Join("backups", book_id+".json")
	if !config.Exist(config.Current.BackupsPath) {
		fmt.Println("book info is not exist, request book info...")
		tools.Mkdir("backups")
		switch config.Vars.AppType {
		case "sfacg":
			err = boluobao.GET_BOOK_INFORMATION(book_id)
		case "cat":
			err = hbooker.GET_BOOK_INFORMATION(book_id)
		}
		if err == nil {
			file.Open(config.Current.BackupsPath, tools.JsonString(config.Current.Book), "w")
		} else {
			return Catalogue{Test: false, BookMessage: fmt.Sprintf("book_id:%v is invalid:%v", book_id, err)}
		}
	}
	_ = json.Unmarshal([]byte(file.ReadFile(config.Current.BackupsPath)), &config.Current.Book)
	OutputPath := tools.Mkdir(path.Join(config.Vars.OutputName, config.Current.Book.NovelName))
	config.Current.ConfigPath = path.Join(config.Vars.ConfigName, config.Current.Book.NovelName)
	config.Current.OutputPath = path.Join(OutputPath, config.Current.Book.NovelName+".txt")
	config.Current.CoverPath = path.Join("cover", config.Current.Book.NovelName+".jpg")
	books := BookInits{BookID: book_id, Locks: nil, ShowBook: true}
	return books.BookDetailed()

}

func (books *BookInits) BookDetailed() Catalogue {
	books.InitEpubFile()
	briefIntroduction := fmt.Sprintf("Name: %v\nBookID: %v\nAuthor: %v\nCount: %v\n\n\n",
		config.Current.Book.NovelName, config.Current.Book.NovelID, config.Current.Book.AuthorName, config.Current.Book.CharCount,
	)
	if books.ShowBook {
		fmt.Println(briefIntroduction)
	}
	file.Open(config.Current.OutputPath, briefIntroduction, "w")
	return Catalogue{Test: true, EpubSetting: books.EpubSetting}
}

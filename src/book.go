package src

import (
	"fmt"
	"github.com/AlexiaVeronica/pineapple-backups/config"
	"github.com/AlexiaVeronica/pineapple-backups/pkg/command"
	"github.com/AlexiaVeronica/pineapple-backups/pkg/epub"
	"github.com/AlexiaVeronica/pineapple-backups/pkg/file"
	"github.com/AlexiaVeronica/pineapple-backups/pkg/tools"
	"github.com/imroc/req/v3"
	"os"
	"path"
)

func InitEpubFile() *epub.Epub {
	AddImage := true
	var epubSetting *epub.Epub
	switch command.Command.AppType {
	case "sfacg":
		epubSetting = epub.NewEpub(config.APP.SFacg.BookInfo.NovelName)
		config.Apps.Config.CoverFile = path.Join(config.Vars.ConfigName, "cover", config.APP.SFacg.BookInfo.NovelName+".jpg")
		epubSetting.SetAuthor(config.APP.SFacg.BookInfo.AuthorName)
		if !config.Exist(config.Apps.Config.CoverFile) {
			if reader, err := req.C().R().Get(config.APP.SFacg.BookInfo.NovelCover); err != nil {
				fmt.Println("download cover failed, err:", err)
				AddImage = false
			} else {
				_ = os.WriteFile(config.Apps.Config.CoverFile, reader.Bytes(), 0666)
			}
		}
		if AddImage {
			epubSetting.AddImage(config.Apps.Config.CoverFile, "")
			epubSetting.SetCover(path.Join("../images", config.APP.SFacg.BookInfo.NovelName+".jpg"), "")
		}
	case "cat":
		epubSetting = epub.NewEpub(config.APP.Hbooker.BookInfo.BookName) // set epub setting and add section
		config.Apps.Config.CoverFile = path.Join(config.Vars.ConfigName, "cover", config.APP.SFacg.BookInfo.NovelName+".jpg")
		epubSetting.SetAuthor(config.APP.Hbooker.BookInfo.AuthorName) // set author
		if !config.Exist(config.Apps.Config.CoverFile) {
			if reader, err := req.C().R().Get(config.APP.Hbooker.BookInfo.Cover); err != nil {
				fmt.Println("download cover failed, err:", err)
				AddImage = false
			} else {
				_ = os.WriteFile(config.Apps.Config.CoverFile, reader.Bytes(), 0666)
			}
		}
		if AddImage {
			_, _ = epubSetting.AddImage(config.Apps.Config.CoverFile, "")
			epubSetting.SetCover(path.Join("../images", config.APP.Hbooker.BookInfo.BookName+".jpg"), "")
		}
	}
	return epubSetting
}

func SettingBooks(bookId string) (*Catalogue, error) {
	var err error
	switch command.Command.AppType {
	case "sfacg":
		config.APP.SFacg.BookInfo, err = config.APP.SFacg.Client.API.GetBookInfo(bookId)
		if err != nil {
			return nil, err
		}
		tools.Mkdir(path.Join(config.Vars.OutputName, config.APP.SFacg.BookInfo.NovelName))
	case "cat":
		config.APP.Hbooker.BookInfo, err = config.APP.Hbooker.Client.API.GetBookInfo(bookId)
		if err != nil {
			return nil, err
		}
		tools.Mkdir(path.Join(config.Vars.OutputName, config.APP.Hbooker.BookInfo.BookName))
	}
	return BookDetailed(), nil

}

func BookDetailed() *Catalogue {

	var briefIntroduction string
	if command.Command.AppType == "sfacg" {
		briefIntroduction = fmt.Sprintf("Name: %v\nBookID: %v\nAuthor: %v\nCount: %v\n\n\n",
			config.APP.SFacg.BookInfo.NovelName, config.APP.SFacg.BookInfo.NovelId, config.APP.SFacg.BookInfo.AuthorName,
			config.APP.SFacg.BookInfo.CharCount,
		)
		file.Open(
			path.Join(config.Vars.OutputName, config.APP.SFacg.BookInfo.NovelName, config.APP.SFacg.BookInfo.NovelName+".txt"),
			briefIntroduction,
			"w",
		)
	} else {
		briefIntroduction = fmt.Sprintf("Name: %v\nBookID: %v\nAuthor: %v\nCount: %v\n\n\n",
			config.APP.Hbooker.BookInfo.BookName, config.APP.Hbooker.BookInfo.BookID, config.APP.Hbooker.BookInfo.AuthorName,
			config.APP.Hbooker.BookInfo.TotalWordCount,
		)

		file.Open(
			path.Join(config.Vars.OutputName, config.APP.Hbooker.BookInfo.BookName, config.APP.Hbooker.BookInfo.BookName+".txt"),
			briefIntroduction,
			"w",
		)
	}
	fmt.Println(briefIntroduction)
	return &Catalogue{EpubSetting: InitEpubFile()}
}

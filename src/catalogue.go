package src

import (
	"fmt"
	"github.com/VeronicaAlexia/BoluobaoAPI/boluobao"
	"github.com/VeronicaAlexia/pineapple-backups/config"
	"github.com/VeronicaAlexia/pineapple-backups/pkg/command"
	"github.com/VeronicaAlexia/pineapple-backups/pkg/epub"
	"github.com/VeronicaAlexia/pineapple-backups/pkg/file"
	"github.com/VeronicaAlexia/pineapple-backups/pkg/threading"
	"github.com/VeronicaAlexia/pineapple-backups/pkg/tools"
	"path"
	"strings"
	"time"
)

type Catalogue struct {
	ChapterBar  *ProgressBar
	ChapterCfg  []string
	EpubSetting *epub.Epub
}

func (catalogue *Catalogue) ReadChapterConfig() {
	if !config.Exist(config.Current.ConfigPath) {
		tools.Mkdir(config.Current.ConfigPath)
		catalogue.ChapterCfg = []string{}
	} else {
		catalogue.ChapterCfg = tools.GetFileName(config.Current.ConfigPath)
	}
}

func (catalogue *Catalogue) GetDownloadsList() []string {
	catalogue.ReadChapterConfig()
	if command.Command.AppType == "sfacg" {
		return boluobao.API.Book.NovelCatalogue(config.Current.NewBooks["novel_id"])
	} else if command.Command.AppType == "cat" {
		var DownloadList []string
		divisionList, err := config.APP.Hbooker.Client.API.GetDivisionListByBookId(config.APP.Hbooker.BookInfo.BookID)
		if err != nil {
			fmt.Println("get division list error:", err)
		} else {
			for _, value := range divisionList {
				for _, value2 := range value.ChapterList {
					if !tools.TestList(catalogue.ChapterCfg, value2.ChapterID) {
						if value2.AuthAccess == "0" || value2.AuthAccess == "1" {
							DownloadList = append(DownloadList, value2.ChapterID)
						} else {
							fmt.Println(value2.ChapterTitle, " is vip chapter, You need to subscribe it")
						}
					}
				}
			}
		}
		fmt.Println(DownloadList)
		return DownloadList
	}
	return nil
}

func (catalogue *Catalogue) DownloadContent(threading *threading.GoLimit, chapterID string) {
	defer threading.Done()
	catalogue.speed_progress()
	if config.Exist(path.Join(config.Current.ConfigPath, chapterID+".txt")) {
		return
	}
	var contentText string
	if command.Command.AppType == "sfacg" {
		response := boluobao.API.Book.NovelContent(chapterID)
		if response != nil {
			contentText = response.Data.Title + "\n" + response.Data.Expand.Content
		} else {
			return
		}
	} else if command.Command.AppType == "cat" {
		chapterKey, err := config.APP.Hbooker.Client.API.GetChapterKey(chapterID)
		if err != nil {
			return
		}
		chapterInfo, err := config.APP.Hbooker.Client.API.GetChapterContentAPI(chapterID, chapterKey)
		if err != nil {
			fmt.Println("get chapter content error:", err)
			return
		}
		contentText = chapterInfo.ChapterTitle + "\n" + chapterInfo.TxtContent
	}
	if contentText != "" {
		file.Open(path.Join(config.Current.ConfigPath, chapterID+".txt"), contentText, "w")
	}
}

func (catalogue *Catalogue) MergeTextAndEpubFiles() {
	var savePath string
	if command.Command.AppType == "sfacg" {
		savePath = path.Join(config.Vars.OutputName, config.Current.NewBooks["novel_name"], config.Current.NewBooks["novel_name"])
	} else {
		savePath = path.Join(config.Vars.OutputName, config.APP.Hbooker.BookInfo.BookName, config.APP.Hbooker.BookInfo.BookName)
	}
	var NovelCatalogue []string
	if command.Command.AppType == "sfacg" {
		NovelCatalogue = boluobao.API.Book.NovelCatalogue(config.Current.NewBooks["novel_id"])
	} else {
		divisionList, err := config.APP.Hbooker.Client.API.GetDivisionListByBookId(config.APP.Hbooker.BookInfo.BookID)
		if err != nil {
			fmt.Println("get division list error:", err)
			return
		}
		for _, i := range divisionList {
			for _, chapterInfo := range i.ChapterList {
				if !tools.TestList(catalogue.ChapterCfg, chapterInfo.ChapterID) {
					if chapterInfo.IsPaid == "0" || chapterInfo.AuthAccess == "1" {
						NovelCatalogue = append(NovelCatalogue, chapterInfo.ChapterID)
					} else {
						fmt.Println(chapterInfo.ChapterTitle, " is vip chapter, You need to subscribe it")
					}
				}
			}
		}
	}

	for _, chapterId := range NovelCatalogue {
		if config.Exist(path.Join(config.Current.ConfigPath, chapterId+".txt")) {
			content := file.Open(path.Join(config.Current.ConfigPath, chapterId+".txt"), "", "r")
			//if config.Vars.Epub {
			file.Open(savePath+".txt", "\n\n\n"+content, "a")
			catalogue.add_chapter_in_epub_file(strings.Split(content, "\n")[0], content+"</p>")
			//} // save to epub file if epub is true
		}

	}

	outPutEpubNow := time.Now() // start time
	if err := catalogue.EpubSetting.Write(savePath + ".epub"); err != nil {
		fmt.Println("output epub error:", err)
	}
	fmt.Println("output epub file success, time:", time.Since(outPutEpubNow))
	//}
}

func (catalogue *Catalogue) add_chapter_in_epub_file(title, content string) {
	xmlContent := "<h1>" + title + "</h1>\n<p>" + strings.ReplaceAll(content, "\n", "</p>\n<p>")
	if _, err := catalogue.EpubSetting.AddSection(xmlContent, title, "", ""); err != nil {
		fmt.Printf("epub add chapter:%v\t\terror message:%v", title, err)
	}
}
func (catalogue *Catalogue) speed_progress() {
	if err := catalogue.ChapterBar.Add(1); err != nil {
		fmt.Println("bar error:", err)
	} else {
		//time.Sleep(time.Second * time.Duration(2))
	}

}

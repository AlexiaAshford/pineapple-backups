package src

import (
	"fmt"
	"github.com/VeronicaAlexia/BoluobaoAPI/boluobao"
	"github.com/VeronicaAlexia/ciweimaoapiLib"
	"github.com/VeronicaAlexia/pineapple-backups/config"
	"github.com/VeronicaAlexia/pineapple-backups/pkg/command"
	"github.com/VeronicaAlexia/pineapple-backups/pkg/epub"
	"github.com/VeronicaAlexia/pineapple-backups/pkg/file"
	"github.com/VeronicaAlexia/pineapple-backups/pkg/threading"
	"github.com/VeronicaAlexia/pineapple-backups/pkg/tools"
	"log"
	"path"
	"strings"
	"time"
)

type Catalogue struct {
	ChapterBar  *ProgressBar
	ChapterCfg  []string
	Test        bool
	BookMessage string
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
		response := ciweimaoapi.GetCatalog(config.Current.NewBooks["novel_id"])
		if response != nil && response.Code == "100000" {
			for _, value := range response.Data.ChapterList {
				for _, value2 := range value.ChapterList {
					if !tools.TestList(catalogue.ChapterCfg, value2.ChapterId) {
						if value2.AuthAccess == "0" || value2.AuthAccess == "1" {
							DownloadList = append(DownloadList, value2.ChapterId)
						} else {
							fmt.Println(value2.ChapterTitle, " is vip chapter, You need to subscribe it")
						}
					}
				}
			}
		} else {
			log.Printf("GetCatalog Error: %v", response.Tip.(string))
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
	var content_text string
	for i := 0; i < 5; i++ {
		if command.Command.AppType == "sfacg" {
			response := boluobao.API.Book.NovelContent(chapterID)
			if response != nil {
				content_text = response.Data.Expand.Content
			} else {
				break
			}
		} else if command.Command.AppType == "cat" {
			content_text = ciweimaoapi.GetContent(chapterID).Data.ChapterInfo.TxtContent
		}
		if content_text != "" {
			file.Open(path.Join(config.Current.ConfigPath, chapterID+".txt"), content_text, "w")
			break
		}
	}
}

func (catalogue *Catalogue) MergeTextAndEpubFiles() {
	savePath := path.Join(config.Vars.OutputName, config.Current.NewBooks["novel_name"], config.Current.NewBooks["novel_name"])
	var NovelCatalogue []string
	if command.Command.AppType == "sfacg" {
		NovelCatalogue = boluobao.API.Book.NovelCatalogue(config.Current.NewBooks["novel_id"])
	} else {
		for _, i := range ciweimaoapi.GetCatalog(config.Current.NewBooks["novel_id"]).Data.ChapterList {
			for _, chapterInfo := range i.ChapterList {
				if !tools.TestList(catalogue.ChapterCfg, chapterInfo.ChapterId) {
					if chapterInfo.AuthAccess == "0" || chapterInfo.AuthAccess == "1" {
						NovelCatalogue = append(NovelCatalogue, chapterInfo.ChapterId)
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

	//for _, local_file_name := range tools.GetFileName(config.Current.ConfigPath) {
	//	content := file.Open(config.Current.ConfigPath+"/"+local_file_name, "", "r")
	//	file.Open(savePath+".txt", "\n\n\n"+content, "a")
	//	if config.Vars.Epub {
	//		catalogue.add_chapter_in_epub_file(strings.Split(content, "\n")[0], content+"</p>")
	//	} // save to epub file if epub is true
	//}
	//if config.Vars.Epub { // output epub file
	out_put_epub_now := time.Now() // start time
	if err := catalogue.EpubSetting.Write(savePath + ".epub"); err != nil {
		fmt.Println("output epub error:", err)
	}
	fmt.Println("output epub file success, time:", time.Since(out_put_epub_now))
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

//func (catalogue *Catalogue) AddChapterConfig(chapId any) {
//	switch chapId.(type) {
//	case string:
//		cfg.Write(cfg.Current.ConfigPath, chapId.(string)+",", "a")
//	case int:
//		cfg.Write(cfg.Current.ConfigPath, strconv.Itoa(chapId.(int))+",", "a")
//	}
//}

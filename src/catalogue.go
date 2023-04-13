package src

import (
	"fmt"
	"github.com/VeronicaAlexia/BoluobaoAPI/boluobao"
	"github.com/VeronicaAlexia/HbookerAPI/ciweimao/book"
	"github.com/VeronicaAlexia/pineapple-backups/config"
	"github.com/VeronicaAlexia/pineapple-backups/pkg/epub"
	"github.com/VeronicaAlexia/pineapple-backups/pkg/file"
	"github.com/VeronicaAlexia/pineapple-backups/pkg/threading"
	"github.com/VeronicaAlexia/pineapple-backups/pkg/tools"
	"path"
	"strconv"
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

func GET_DIVISION(BookId string) []map[string]string {
	var chapter_index int
	var division_info_list []map[string]string
	VolumeList := book.GET_DIVISION_LIST_BY_BOOKID(BookId)
	for division_index, division_info := range VolumeList.Data.ChapterList {
		fmt.Printf("第%v卷\t\t%v\n", division_index+1, division_info.DivisionName)
		for _, chapter := range division_info.ChapterList {
			chapter_index += 1
			division_info_list = append(division_info_list, map[string]string{
				"is_valid":       chapter.IsValid,
				"chapter_id":     chapter.ChapterID,
				"money":          chapter.AuthAccess,
				"chapter_name":   chapter.ChapterTitle,
				"division_name":  division_info.DivisionName,
				"division_id":    division_info.DivisionID,
				"division_index": strconv.Itoa(division_index),
				"chapter_index":  strconv.Itoa(chapter_index),
				"file_name":      file.FileCacheName(division_index, chapter_index, chapter.ChapterID),
			})
		}
	}
	return division_info_list
}

func (catalogue *Catalogue) GetDownloadsList() []string {
	catalogue.ReadChapterConfig()
	var chapter_info_list []map[string]string
	if config.Vars.AppType == "sfacg" {
		return boluobao.API.Book.NovelCatalogue(config.Current.NewBooks["novel_id"])
	} else if config.Vars.AppType == "cat" {
		var DownloadList []string
		chapter_info_list = GET_DIVISION(config.Current.NewBooks["novel_id"])
		for _, chapter_info := range chapter_info_list {
			if !tools.TestList(catalogue.ChapterCfg, chapter_info["file_name"]) {
				if chapter_info["money"] == "0" || chapter_info["money"] == "1" {
					DownloadList = append(DownloadList, chapter_info["chapter_id"])
				} else {
					fmt.Println(chapter_info["chapter_name"], " is vip chapter, You need to subscribe it")
				}
			}
		}
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
		if config.Vars.AppType == "sfacg" {
			response := boluobao.API.Book.NovelContent(chapterID)
			if response != nil {
				content_text = response.Data.Expand.Content
			} else {
				break
			}
		} else if config.Vars.AppType == "cat" {
			content_text = book.GetContent(chapterID)
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
	if config.Vars.AppType == "sfacg" {
		NovelCatalogue = boluobao.API.Book.NovelCatalogue(config.Current.NewBooks["novel_id"])
	} else {
		for _, chapter_info := range GET_DIVISION(config.Current.NewBooks["novel_id"]) {
			if !tools.TestList(catalogue.ChapterCfg, chapter_info["file_name"]) {
				if chapter_info["money"] == "0" || chapter_info["money"] == "1" {
					NovelCatalogue = append(NovelCatalogue, chapter_info["chapter_id"])
				} else {
					fmt.Println(chapter_info["chapter_name"], " is vip chapter, You need to subscribe it")
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

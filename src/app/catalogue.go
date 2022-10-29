package app

import (
	"fmt"
	"github.com/VeronicaAlexia/pineapple-backups/config"
	"github.com/VeronicaAlexia/pineapple-backups/epub"
	"github.com/VeronicaAlexia/pineapple-backups/pkg/file"
	"github.com/VeronicaAlexia/pineapple-backups/pkg/tools"
	"github.com/VeronicaAlexia/pineapple-backups/src/app/boluobao"
	"github.com/VeronicaAlexia/pineapple-backups/src/app/hbooker"
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

func (catalogue *Catalogue) GetDownloadsList() {
	catalogue.ReadChapterConfig()
	var chapter_info_list []map[string]string
	if config.Vars.AppType == "sfacg" {
		chapter_info_list = boluobao.GET_CATALOGUE(config.Current.Book.NovelID)
	} else if config.Vars.AppType == "cat" {
		chapter_info_list = hbooker.GET_DIVISION(config.Current.Book.NovelID)
	}
	for _, chapter_info := range chapter_info_list {
		if !tools.TestList(catalogue.ChapterCfg, chapter_info["file_name"]) {
			if chapter_info["money"] == "0" || chapter_info["money"] == "1" {
				config.Current.DownloadList = append(config.Current.DownloadList, chapter_info["file_name"])
			} else {
				fmt.Println(chapter_info["chapter_name"], " is vip chapter, You need to subscribe it")
			}
		}
	}
}

func (catalogue *Catalogue) DownloadContent(threading *config.GoLimit, file_name string) {
	defer threading.Done()
	chapter_id := catalogue.speed_progress(file_name)
	var content_text string
	for i := 0; i < 5; i++ {
		if config.Vars.AppType == "sfacg" {
			content_text = boluobao.GET_CHAPTER_CONTENT(chapter_id)
		} else if config.Vars.AppType == "cat" {
			content_text = hbooker.GET_CHAPTER_CONTENT(chapter_id, hbooker.GET_KET_BY_CHAPTER_ID(chapter_id))
		}
		if content_text != "" {
			file.Open(path.Join(config.Current.ConfigPath, file_name), content_text, "w")
			break
		}
	}
}

func (catalogue *Catalogue) MergeTextAndEpubFiles() {
	for _, local_file_name := range tools.GetFileName(config.Current.ConfigPath) {
		content := file.Open(config.Current.ConfigPath+"/"+local_file_name, "", "r")
		file.Open(config.Current.OutputPath, "\n\n\n"+content, "a")
		if config.Vars.Epub {
			catalogue.add_chapter_in_epub_file(strings.Split(content, "\n")[0], content+"</p>")
		} // save to epub file if epub is true
	}
	if config.Vars.Epub { // output epub file
		out_put_epub_now := time.Now() // start time
		if err := catalogue.EpubSetting.Write(strings.ReplaceAll(config.Current.OutputPath, ".txt", ".epub")); err != nil {
			fmt.Println("output epub error:", err)
		}
		fmt.Println("output epub file success, time:", time.Since(out_put_epub_now))
	}
	config.Current.DownloadList = nil
}

func (catalogue *Catalogue) add_chapter_in_epub_file(title, content string) {
	xmlContent := "<h1>" + title + "</h1>\n<p>" + strings.ReplaceAll(content, "\n", "</p>\n<p>")
	if _, err := catalogue.EpubSetting.AddSection(xmlContent, title, "", ""); err != nil {
		fmt.Printf("epub add chapter:%v\t\terror message:%v", title, err)
	}
}
func (catalogue *Catalogue) speed_progress(file_name string) string {
	if err := catalogue.ChapterBar.Add(1); err != nil {
		fmt.Println("bar error:", err)
	} else {
		//time.Sleep(time.Second * time.Duration(2))
	}
	return strings.ReplaceAll(strings.Split(file_name, "-")[2], ".txt", "")

}

//func (catalogue *Catalogue) AddChapterConfig(chapId any) {
//	switch chapId.(type) {
//	case string:
//		cfg.Write(cfg.Current.ConfigPath, chapId.(string)+",", "a")
//	case int:
//		cfg.Write(cfg.Current.ConfigPath, strconv.Itoa(chapId.(int))+",", "a")
//	}
//}

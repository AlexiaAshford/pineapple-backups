package src

import (
	"fmt"
	"github.com/VeronicaAlexia/BoluobaoAPI/boluobao/book"
	HbookerAPI "github.com/VeronicaAlexia/HbookerAPI/ciweimao/book"
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

func GET_CATALOGUE(NovelID string) []map[string]string {
	var chapter_index int
	var division_info []map[string]string
	response := book.GET_CATALOGUE(NovelID)
	for division_index, division := range response.Data.VolumeList {
		fmt.Printf("第%v卷\t\t%v\n", division_index+1, division.Title)
		for _, chapter := range division.ChapterList {
			chapter_index += 1
			division_info = append(division_info, map[string]string{
				"division_name":  division.Title,
				"division_id":    strconv.Itoa(division.VolumeID),
				"division_index": strconv.Itoa(division_index),
				"chapter_name":   chapter.Title,
				"chapter_id":     strconv.Itoa(chapter.ChapID),
				"chapter_index":  strconv.Itoa(chapter_index),
				"money":          strconv.Itoa(chapter.OriginNeedFireMoney),
				"file_name":      file.NameSetting(chapter.VolumeID, chapter.ChapOrder, chapter.ChapID),
			})
		}
	}
	return division_info
}

func GET_DIVISION_Hbooker(BookId string) []map[string]string {
	var chapter_index int
	var division_info_list []map[string]string
	VolumeList := HbookerAPI.GET_DIVISION_LIST_BY_BOOKID(BookId)
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

func (catalogue *Catalogue) GetDownloadsList() {
	catalogue.ReadChapterConfig()
	var chapter_info_list []map[string]string
	if config.Vars.AppType == "sfacg" {
		chapter_info_list = GET_CATALOGUE(config.Current.NewBooks["novel_id"])
	} else if config.Vars.AppType == "cat" {
		chapter_info_list = GET_DIVISION_Hbooker(config.Current.NewBooks["novel_id"])
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

func (catalogue *Catalogue) DownloadContent(threading *threading.GoLimit, file_name string) {
	defer threading.Done()
	chapter_id := catalogue.speed_progress(file_name)
	var content_text string
	for i := 0; i < 5; i++ {
		if config.Vars.AppType == "sfacg" {
			content_text = book.Content(chapter_id).Data.Expand.Content
		} else if config.Vars.AppType == "cat" {
			content_text = HbookerAPI.GetContent(chapter_id)
		}
		if content_text != "" {
			file.Open(path.Join(config.Current.ConfigPath, file_name), content_text, "w")
			break
		}
	}
}

func (catalogue *Catalogue) MergeTextAndEpubFiles() {
	savePath := path.Join(config.Vars.OutputName, config.Current.NewBooks["novel_name"], config.Current.NewBooks["novel_name"])
	for _, local_file_name := range tools.GetFileName(config.Current.ConfigPath) {
		content := file.Open(config.Current.ConfigPath+"/"+local_file_name, "", "r")
		file.Open(savePath+".txt", "\n\n\n"+content, "a")
		if config.Vars.Epub {
			catalogue.add_chapter_in_epub_file(strings.Split(content, "\n")[0], content+"</p>")
		} // save to epub file if epub is true
	}
	if config.Vars.Epub { // output epub file
		out_put_epub_now := time.Now() // start time
		if err := catalogue.EpubSetting.Write(savePath + ".epub"); err != nil {
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

package src

import (
	"fmt"
	"io/ioutil"
	"path"
	"sf/cfg"
	"sf/epub"
	"sf/src/boluobao"
	"sf/src/hbooker"
	"sf/src/hbooker/Encrypt"
	"strconv"
	"strings"
	"time"
)

type Catalogue struct {
	ChapterBar     *ProgressBar
	ChapterCfg     string
	TestBookResult bool
	EpubSetting    *epub.Epub
}

func (catalogue *Catalogue) ReadChapterConfig() {
	if !cfg.Exist(cfg.Current.ConfigPath) {
		cfg.Mkdir(cfg.Current.ConfigPath)
		catalogue.ChapterCfg = ""
	} else {
		FileInfo, _ := ioutil.ReadDir(path.Join(cfg.Vars.ConfigName, cfg.Current.Book.NovelName))
		for _, v := range FileInfo {
			catalogue.ChapterCfg += v.Name() + ","
		}
	}
}

func (catalogue *Catalogue) GetDownloadsList() {
	catalogue.ReadChapterConfig()
	switch cfg.Vars.AppType {
	case "sfacg":
		for divisionIndex, division := range boluobao.GET_CATALOGUE(cfg.Current.Book.NovelID).Data.VolumeList {
			fmt.Printf("第%v卷\t\t%v\n", divisionIndex+1, division.Title)
			for _, chapter := range division.ChapterList {
				//if cfg.TestKeyword(catalogue.ChapterCfg, chapter.ChapID) {
				//	fmt.Println("已下载：", chapter.Title)
				//}
				if chapter.OriginNeedFireMoney == 0 && !cfg.TestKeyword(catalogue.ChapterCfg, chapter.ChapID) {
					cfg.Current.DownloadList = append(cfg.Current.DownloadList, strconv.Itoa(chapter.ChapID))
				}
			}
		}
	case "cat":
		for index, division := range hbooker.GET_DIVISION(cfg.Current.Book.NovelID) {
			fmt.Printf("第%v卷\t\t%v\n", index+1, division.DivisionName)
			for _, chapter := range hbooker.GET_CATALOGUE(division.DivisionID) {
				if chapter.IsValid == "1" && !cfg.TestKeyword(catalogue.ChapterCfg, chapter.ChapterID) {
					cfg.Current.DownloadList = append(cfg.Current.DownloadList, chapter.ChapterID)
				}
			}
		}

	}

}

func config_file_name(index, VolumeID, ChapID any) string {
	index = cfg.StrToInt(fmt.Sprintf("%d", index))
	return fmt.Sprintf("%05d", index) + "-" + fmt.Sprintf("%v", VolumeID) + "-" + fmt.Sprintf("%v", ChapID) + ".txt"
}

func (catalogue *Catalogue) DownloadContent(chapter_id string) {
	catalogue.SpeedProgressAndDelayTime()
	if cfg.Vars.AppType == "sfacg" {
		response := boluobao.GET_CONTENT(chapter_id)
		if response.Status.HTTPCode == 200 {
			result := response.Data // get content data
			file_name := config_file_name(result.ChapOrder, result.VolumeID, result.ChapID)
			content_text := fmt.Sprintf("%v:%v\n%v", result.Title, result.AddTime, result.Expand.Content)
			cfg.Write(path.Join(cfg.Current.ConfigPath, file_name), content_text, "w")
		} else {
			fmt.Println("Error:", response.Status.Msg)
		}

	} else if cfg.Vars.AppType == "cat" {
		chapter_key := hbooker.GetKeyByCid(chapter_id)
		response := hbooker.GetContent(chapter_id, chapter_key)
		if response.Code == "100000" {
			TxtContent := Encrypt.Decode(response.Data.ChapterInfo.TxtContent, chapter_key)
			result := response.Data.ChapterInfo
			file_name := config_file_name(result.ChapterIndex, result.DivisionID, result.ChapterID)
			content_text := fmt.Sprintf("%v:%v\n%v\n\n\n", result.ChapterTitle, result.Uptime, string(TxtContent))
			cfg.Write(path.Join(cfg.Current.ConfigPath, file_name), content_text, "w")
		}

	} else {
		panic("app type error please check config file, the app type is:" + cfg.Vars.AppType)
	}

}

func (catalogue *Catalogue) MergeFiles() {
	//cfg.Write(cfg.Current.OutputPath, catalogue.ContentList["cache"], "w")
	//for _, ChapterId := range cfg.Current.DownloadList {
	//	cfg.Write(cfg.Current.OutputPath, catalogue.ContentList[ChapterId], "a")
	//}
	FileInfo, _ := ioutil.ReadDir(path.Join(cfg.Vars.ConfigName, cfg.Current.Book.NovelName))
	for _, v := range FileInfo {
		content := cfg.Write(cfg.Current.ConfigPath+"/"+v.Name(), "", "r")
		catalogue.AddChapterInEpubFile(strings.Split(content, "\n")[0], content)
		cfg.Write(cfg.Current.OutputPath, content, "a")
	}
	bT := time.Now() // 开始时间
	// save epub file
	epub_file_name := strings.Replace(cfg.Current.OutputPath, ".txt", ".epub", -1)
	if err := catalogue.EpubSetting.Write(epub_file_name); err != nil {
		fmt.Println(epub_file_name, " epub error:", err)
	}
	fmt.Println("write epub: ", time.Since(bT)) // 从开始到当前所消耗的时间
}

func (catalogue *Catalogue) AddChapterInEpubFile(title, content string) {
	xmlContent := "<h1>" + title + "</h1>\n<p>" + strings.Replace(content, "\n", "</p>\n<p>", -1)
	if _, err := catalogue.EpubSetting.AddSection(xmlContent, title, "", ""); err != nil {
		fmt.Printf("epub add chapter:%v\t\terror message:%v", title, err)
	}
}
func (catalogue *Catalogue) SpeedProgressAndDelayTime() {
	if err := catalogue.ChapterBar.Add(1); err != nil {
		fmt.Println("bar error:", err)
	} else {
		//time.Sleep(time.Second * time.Duration(2))
	}
}

//func (catalogue *Catalogue) makeContentInformation(response any) {
//	cfg.FileLock.Lock()         // lock file to avoid file write conflict
//	defer cfg.FileLock.Unlock() // unlock file after write
//	var writeContent string
//	switch response.(type) {
//	case *sfacg_structs.Content:
//		result := response.(*sfacg_structs.Content).Data
//		writeContent = fmt.Sprintf("%v:%v\n%v\n\n\n", result.Title, result.AddTime, result.Expand.Content)
//		catalogue.AddChapterConfig(result.ChapID)
//		catalogue.AddChapterInEpubFile(result.Title, result.Expand.Content)
//		catalogue.ContentList[strconv.Itoa(result.ChapID)] = writeContent
//	case *hbooker_structs.ContentStruct:
//		result := response.(*hbooker_structs.ContentStruct).Data.ChapterInfo
//		writeContent = fmt.Sprintf("%v:%v\n%v\n\n\n", result.ChapterTitle, result.Uptime, result.TxtContent)
//		catalogue.AddChapterConfig(result.ChapterID)
//		catalogue.AddChapterInEpubFile(result.ChapterTitle, result.TxtContent)
//		catalogue.ContentList[result.ChapterID] = writeContent
//	}
//	catalogue.SpeedProgressAndDelayTime()
//
//}

//func (catalogue *Catalogue) AddChapterConfig(chapId any) {
//	switch chapId.(type) {
//	case string:
//		cfg.Write(cfg.Current.ConfigPath, chapId.(string)+",", "a")
//	case int:
//		cfg.Write(cfg.Current.ConfigPath, strconv.Itoa(chapId.(int))+",", "a")
//	}
//}

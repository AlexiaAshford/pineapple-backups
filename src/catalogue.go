package src

import (
	"fmt"
	"sf/cfg"
	"sf/epub"
	"sf/src/boluobao"
	"sf/src/hbooker"
	"sf/struct/hbooker_structs"
	"sf/struct/sfacg_structs"
	"strconv"
	"strings"
)

type Catalogue struct {
	ChapterBar     *ProgressBar
	ChapterList    []string
	ChapterCfg     string
	TestBookResult bool
	contentList    map[string]string
	EpubSetting    *epub.Epub
}

func (catalogue *Catalogue) ReadChapterConfig() string {
	catalogue.contentList = make(map[string]string)
	if !cfg.Exist(cfg.Current.ConfigPath) {
		cfg.Write(cfg.Current.ConfigPath, "", "w")
		catalogue.ChapterCfg = ""
	} else { // read config file
		catalogue.ChapterCfg = cfg.Write(cfg.Current.ConfigPath, "", "r")
	}
	catalogue.contentList["cache"] = cfg.Write(cfg.Current.OutputPath, "", "r")

	return cfg.Vars.AppType
}
func (catalogue *Catalogue) AddChapterConfig(chapId any) {
	switch chapId.(type) {
	case string:
		cfg.Write(cfg.Current.ConfigPath, chapId.(string)+",", "a")
	case int:
		cfg.Write(cfg.Current.ConfigPath, strconv.Itoa(chapId.(int))+",", "a")
	}
}

func (catalogue *Catalogue) InitCatalogue() {
	switch catalogue.ReadChapterConfig() {
	case "sfacg":
		for divisionIndex, division := range boluobao.GET_CATALOGUE(cfg.Current.Book.NovelID).Data.VolumeList {
			fmt.Printf("第%v卷\t\t%v\n", divisionIndex+1, division.Title)
			for _, chapter := range division.ChapterList {
				if chapter.OriginNeedFireMoney == 0 && !cfg.TestKeyword(catalogue.ChapterCfg, chapter.ChapID) {
					catalogue.ChapterList = append(catalogue.ChapterList, strconv.Itoa(chapter.ChapID))
				}
			}
		}
	case "cat":
		for index, division := range hbooker.GetDivisionIdByBookId(cfg.Current.Book.NovelID) {
			fmt.Printf("第%v卷\t\t%v\n", index+1, division.DivisionName)
			for _, chapter := range hbooker.GetCatalogueByDivisionId(division.DivisionID) {
				if chapter.IsValid == "1" && !cfg.TestKeyword(catalogue.ChapterCfg, chapter.ChapterID) {
					catalogue.ChapterList = append(catalogue.ChapterList, chapter.ChapterID)
				}
			}
		}

	}
	if len(catalogue.ChapterList) > 0 {
		catalogue.ChapterBar = New(len(catalogue.ChapterList))
		catalogue.ChapterBar.Describe("working...")
		catalogue.DownloadContent()
		fmt.Printf("\nNovel:%v download complete!\n", cfg.Current.Book.NovelName)
		fmt.Println("The txt file is out put to:", cfg.Current.OutputPath)
	} else {
		fmt.Println("No chapter need to download!")
	}
}

func (catalogue *Catalogue) DownloadContent() {
	for _, ChapterId := range catalogue.ChapterList {

		func(ChapterId string) {
			if cfg.Vars.AppType == "sfacg" {
				for retry := 0; retry < cfg.Vars.MaxRetry; retry++ {
					response := boluobao.GET_CONTENT(ChapterId)
					if response.Status.HTTPCode == 200 {
						catalogue.makeContentInformation(response)
					} else {
						fmt.Println("Error:", response.Status.Msg)
					}
				}

			} else if cfg.Vars.AppType == "cat" {
				chapter_key := hbooker.GetKeyByCid(ChapterId)
				if response, ok := hbooker.GetContent(ChapterId, chapter_key); ok {
					catalogue.makeContentInformation(response)
				}
			} else {
				panic("app type error please check config file, the app type is:" + cfg.Vars.AppType)
			}
		}(ChapterId)

	}
	cfg.Write(cfg.Current.OutputPath, catalogue.contentList["cache"], "w")
	for _, ChapterId := range catalogue.ChapterList {
		cfg.Write(cfg.Current.OutputPath, catalogue.contentList[ChapterId], "a")
	}

	if err := catalogue.EpubSetting.Write(strings.Replace(cfg.Current.OutputPath, ".txt", ".epub", -1)); err != nil {
		fmt.Println("epub error:", err)
	}
	catalogue.ChapterList = nil
}

func (catalogue *Catalogue) AddChapterInEpubFile(title, content string) {
	xmlContent := "<h1>" + title + "</h1>\n<p>" + strings.Replace(content, "\n", "</p>\n<p>", -1)
	if _, err := catalogue.EpubSetting.AddSection(xmlContent, title, "", ""); err != nil {
		fmt.Printf("epub add chapter:%v\t\terror message:%v", title, err)
	}
}
func (catalogue *Catalogue) makeContentInformation(response any) {
	cfg.FileLock.Lock()         // lock file to avoid file write conflict
	defer cfg.FileLock.Unlock() // unlock file after write
	var writeContent string
	switch response.(type) {
	case *sfacg_structs.Content:
		result := response.(*sfacg_structs.Content).Data
		writeContent = fmt.Sprintf("%v:%v\n%v\n\n\n", result.Title, result.AddTime, result.Expand.Content)
		catalogue.AddChapterConfig(result.ChapID)
		catalogue.AddChapterInEpubFile(result.Title, result.Expand.Content)
		catalogue.contentList[strconv.Itoa(result.ChapID)] = writeContent
	case *hbooker_structs.ContentStruct:
		result := response.(*hbooker_structs.ContentStruct).Data.ChapterInfo
		writeContent = fmt.Sprintf("%v:%v\n%v\n\n\n", result.ChapterTitle, result.Uptime, result.TxtContent)
		catalogue.AddChapterConfig(result.ChapterID)
		catalogue.AddChapterInEpubFile(result.ChapterTitle, result.TxtContent)
		catalogue.contentList[result.ChapterID] = writeContent
	}
	catalogue.SpeedProgressAndDelayTime()

}

func (catalogue *Catalogue) SpeedProgressAndDelayTime() {
	if err := catalogue.ChapterBar.Add(1); err != nil {
		fmt.Println("bar error:", err)
	} else {
		//time.Sleep(time.Second * time.Duration(2))
	}
}

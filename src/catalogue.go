package src

import (
	"fmt"
	"path"
	"sf/cfg"
	"sf/src/boluobao"
	"sf/src/hbooker"
	"sf/structural/hbooker_structs"
	"sf/structural/sfacg_structs"
	"strconv"
)

type Catalogue struct {
	ChapterBar     *ProgressBar
	ChapterList    []string
	ConfigPath     string
	SaveTextPath   string
	ChapterCfg     string
	TestBookResult bool
	contentList    map[string]string
}

func (catalogue *Catalogue) ReadChapterConfig() string {
	catalogue.ConfigPath = path.Join(cfg.Vars.ConfigFile, cfg.CurrentBook.BookInfo.NovelName+".conf")
	catalogue.contentList = make(map[string]string)
	if !cfg.Exist(catalogue.ConfigPath) {
		cfg.Write(catalogue.ConfigPath, "", "w")
		catalogue.ChapterCfg = ""
	} else { // read config file
		catalogue.ChapterCfg = cfg.Write(catalogue.ConfigPath, "", "r")
	}
	catalogue.contentList["cache"] = cfg.Write(catalogue.SaveTextPath, "", "r")

	return cfg.Vars.AppType
}
func (catalogue *Catalogue) AddChapterConfig(chapId any) {
	switch chapId.(type) {
	case string:
		cfg.Write(catalogue.ConfigPath, chapId.(string)+",", "a")
	case int:
		cfg.Write(catalogue.ConfigPath, strconv.Itoa(chapId.(int))+",", "a")
	}
}

func (catalogue *Catalogue) InitCatalogue() {
	switch catalogue.ReadChapterConfig() {
	case "sfacg":
		for divisionIndex, division := range boluobao.GetCatalogue(cfg.CurrentBook.BookInfo.NovelID).Data.VolumeList {
			fmt.Printf("第%v卷\t\t%v\n", divisionIndex+1, division.Title)
			for _, chapter := range division.ChapterList {
				if chapter.OriginNeedFireMoney == 0 && !cfg.TestKeyword(catalogue.ChapterCfg, chapter.ChapID) {
					catalogue.ChapterList = append(catalogue.ChapterList, strconv.Itoa(chapter.ChapID))
				}
			}
		}
	case "cat":
		for index, division := range hbooker.GetDivisionIdByBookId(cfg.CurrentBook.BookInfo.NovelID) {
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
		fmt.Printf("\nNovel:%v download complete!\n", cfg.CurrentBook.BookInfo.NovelName)
		fmt.Println("The txt file is out put to:", catalogue.SaveTextPath)
	} else {
		fmt.Println("No chapter need to download!")
	}
}

func (catalogue *Catalogue) DownloadContent() {
	for _, ChapterId := range catalogue.ChapterList {

		func(ChapterId string) {
			if cfg.Vars.AppType == "sfacg" {
				if response, ok := boluobao.GetContentDetailedByCid(ChapterId); ok {
					catalogue.makeContentInformation(response)
				}
			} else if cfg.Vars.AppType == "cat" {
				if response, ok := hbooker.GetContent(ChapterId); ok {
					catalogue.makeContentInformation(response)
				}
			} else {
				panic("app type error please check config file, the app type is:" + cfg.Vars.AppType)
			}
		}(ChapterId)

	}
	cfg.Write(catalogue.SaveTextPath, catalogue.contentList["cache"], "w")
	for _, ChapterId := range catalogue.ChapterList {
		cfg.Write(catalogue.SaveTextPath, catalogue.contentList[ChapterId], "a")
	}
	catalogue.ChapterList = nil
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
		catalogue.contentList[strconv.Itoa(result.ChapID)] = writeContent
	case *hbooker_structs.ContentStruct:
		result := response.(*hbooker_structs.ContentStruct).Data.ChapterInfo
		writeContent = fmt.Sprintf("%v:%v\n%v\n\n\n", result.ChapterTitle, result.Uptime, result.TxtContent)
		catalogue.AddChapterConfig(result.ChapterID)
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

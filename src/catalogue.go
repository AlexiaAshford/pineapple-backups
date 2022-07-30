package src

import (
	"fmt"
	"math/rand"
	"path"
	"sf/cfg"
	"sf/src/boluobao"
	"sf/src/hbooker"
	"sf/structural/hbooker_structs"
	"sf/structural/sfacg_structs"
	"strconv"
	"time"
)

type Catalogue struct {
	ChapterBar   *ProgressBar
	ChapterList  []string
	ConfigPath   string
	SaveTextPath string
	ChapterCfg   string
}

func (catalogue *Catalogue) ReadChapterConfig() string {
	catalogue.ConfigPath = path.Join(cfg.Vars.ConfigFile, cfg.BookConfig.BookInfo.NovelName+".conf")
	if !cfg.CheckFileExist(catalogue.ConfigPath) {
		cfg.EncapsulationWrite(catalogue.ConfigPath, "", 5, "w")
		catalogue.ChapterCfg = ""
	} else { // read config file
		catalogue.ChapterCfg = cfg.EncapsulationWrite(catalogue.ConfigPath, "", 5, "r")
	}
	return cfg.Vars.AppType
}
func (catalogue *Catalogue) AddChapterConfig(chapId any) {
	switch chapId.(type) {
	case string:
		cfg.EncapsulationWrite(catalogue.ConfigPath, chapId.(string)+",", 5, "a")
	case int:
		cfg.EncapsulationWrite(catalogue.ConfigPath, strconv.Itoa(chapId.(int))+",", 5, "a")
	}
}

func (catalogue *Catalogue) InitCatalogue() {
	switch catalogue.ReadChapterConfig() {
	case "sfacg":
		for divisionIndex, division := range boluobao.GetCatalogue(cfg.BookConfig.BookInfo.NovelID).Data.VolumeList {
			fmt.Printf("第%v卷\t\t%v\n", divisionIndex+1, division.Title)
			for _, chapter := range division.ChapterList {
				if chapter.OriginNeedFireMoney == 0 && !cfg.TestKeyword(catalogue.ChapterCfg, chapter.ChapID) {
					catalogue.ChapterList = append(catalogue.ChapterList, strconv.Itoa(chapter.ChapID))
				}
			}
		}
	case "cat":
		for index, division := range hbooker.GetDivisionIdByBookId(cfg.BookConfig.BookInfo.NovelID) {
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
		catalogue.ChapterBar.Describe("进度:")
		catalogue.DownloadContent()
		fmt.Printf("\nNovelName:%vdownload complete!", cfg.BookConfig.BookInfo.NovelName)
	} else {
		fmt.Println("No chapter need to download!")
	}
}

func (catalogue *Catalogue) DownloadContent() {
	for _, ChapterId := range catalogue.ChapterList {
		for i := 0; i < 5; i++ {
			if cfg.Vars.AppType == "sfacg" {
				response := boluobao.GetContentDetailedByCid(ChapterId)
				if response.Status.HTTPCode == 200 {
					catalogue.makeContentInformation(response)
					break
				} else {
					fmt.Println("error message:", response.Status.Msg)
				}
			} else if cfg.Vars.AppType == "cat" {
				response := hbooker.GetContent(ChapterId)
				if response.Code == "100000" {
					catalogue.makeContentInformation(response)
					break
				}
			} else {
				panic("app type error please check config file")
			}
		}
	}
	catalogue.ChapterList = nil
}

func (catalogue *Catalogue) makeContentInformation(response any) {
	var writeContent string
	switch response.(type) {
	case sfacg_structs.Content:
		result := response.(sfacg_structs.Content).Data
		writeContent = fmt.Sprintf("%v:%v\n%v\n\n\n", result.Title, result.AddTime, result.Expand.Content)
		catalogue.AddChapterConfig(result.ChapID)
	case hbooker_structs.ContentStruct:
		result := response.(hbooker_structs.ContentStruct).Data.ChapterInfo
		writeContent = fmt.Sprintf("%v:%v\n%v\n\n\n", result.ChapterTitle, result.Uptime, result.TxtContent)
		catalogue.AddChapterConfig(result.ChapterID)
	}
	cfg.EncapsulationWrite(catalogue.SaveTextPath, writeContent, 5, "a")
	catalogue.SpeedProgressAndDelayTime()

}

func (catalogue *Catalogue) SpeedProgressAndDelayTime() {
	if err := catalogue.ChapterBar.Add(1); err != nil {
		fmt.Println("bar error:", err)
	} else {
		time.Sleep(time.Second * time.Duration(rand.Intn(2)))
	}
}

package src

import (
	"fmt"
	"math/rand"
	"os"
	"path"
	"sf/cfg"
	"sf/src/boluobao"
	"sf/src/hbooker"
	"sf/structural/sfacg_structs"
	"strconv"
	"time"
)

type Catalogue struct {
	ChapterBar  *ProgressBar
	ChapterList []string
}

func (catalogue *Catalogue) SfacgCatalogue() bool {
	response := boluobao.GetCatalogueDetailedById(cfg.Vars.BookInfo.NovelID)
	for divisionIndex, division := range response.Data.VolumeList {
		fmt.Println("index:", divisionIndex, "\t\tdivision:", division.Title)
		for _, chapter := range division.ChapterList {
			if chapter.OriginNeedFireMoney == 0 {
				catalogue.ChapterList = append(catalogue.ChapterList, strconv.Itoa(chapter.ChapID))
			} else {
				if chapter.OriginNeedFireMoney == 0 {
					fmt.Println(chapter.Title, "已经下载过了")
				}
			}

		}
	}
	catalogue.ChapterBar = New(len(catalogue.ChapterList))
	for _, ChapID := range catalogue.ChapterList {
		catalogue.SfacgContent(ChapID)
	}
	catalogue.ChapterList = nil
	return true
}

func (catalogue *Catalogue) SfacgContent(ChapterId string) {
	catalogue.DelayTime()
	response := boluobao.GetContentDetailedByCid(ChapterId)
	for i := 0; i < 5; i++ {
		if response.Status.HTTPCode == 200 {
			makeContentInformation(response)
			break
		} else {
			if response.Status.Msg == "接口校验失败,请尽快把APP升级到最新版哦~" || i == 4 {
				fmt.Println(response.Status.Msg)
				os.Exit(0)
			} else {
				fmt.Println(response.Status.Msg)
			}
		}
	}
}

func makeContentInformation(response sfacg_structs.Content) {
	writeContent := fmt.Sprintf("%v:%v\n%v\n\n\n",
		response.Data.Title, response.Data.AddTime, response.Data.Expand.Content,
	)
	configPath := path.Join(cfg.Vars.ConfigFile, cfg.Vars.BookInfo.NovelName+".conf")
	cfg.EncapsulationWrite(
		path.Join(cfg.Vars.SaveFile, cfg.Vars.BookInfo.NovelName+".txt"), writeContent, 5, "a",
	)
	cfg.EncapsulationWrite(configPath, strconv.Itoa(response.Data.ChapID)+",", 5, "a")

}

func (catalogue *Catalogue) CatCatalogue() bool {
	for index, division := range hbooker.GetDivisionIdByBookId(cfg.Vars.BookInfo.NovelID) {
		fmt.Println("index:", index, "\t\tdivision:", division.DivisionName)
		for _, chapter := range hbooker.GetCatalogueByDivisionId(division.DivisionID) {
			if chapter.IsValid == "1" {
				catalogue.ChapterList = append(catalogue.ChapterList, chapter.ChapterID)
			}
		}
	}
	catalogue.ChapterBar = New(len(catalogue.ChapterList))
	for _, ChapterID := range catalogue.ChapterList {
		catalogue.CatContent(ChapterID)
	}
	catalogue.ChapterList = nil
	return true
}

func (catalogue *Catalogue) CatContent(ChapterId string) {
	catalogue.DelayTime()
	response := hbooker.GetContent(ChapterId)
	writeContent, SavePath := fmt.Sprintf("%v:%v\n%v\n\n\n",
		response.ChapterTitle,
		response.Uptime,
		response.TxtContent,
	), path.Join(cfg.Vars.SaveFile, cfg.Vars.BookInfo.NovelName+".txt")
	cfg.EncapsulationWrite(SavePath, writeContent, 5, "a")

}
func (catalogue *Catalogue) DelayTime() {
	if err := catalogue.ChapterBar.Add(1); err != nil {
		fmt.Println("bar error:", err)
	} else {
		time.Sleep(time.Second * time.Duration(rand.Intn(2)))
	}
}

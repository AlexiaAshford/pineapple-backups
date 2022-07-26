package src

import (
	"fmt"
	"os"
	"path"
	"sf/cfg"
	"sf/src/boluobao"
	"sf/src/hbooker"
	structs "sf/structural/hbooker_structs"
	"sf/structural/sfacg_structs"
	"strconv"
)

type Catalogue struct {
	ChapterBar    *ProgressBar
	ChapterList   []structs.ChapterList
	ChapterConfig string
}

func (catalogue *Catalogue) SfacgCatalogue() bool {
	response := boluobao.GetCatalogueDetailedById(cfg.Vars.BookInfo.NovelID)
	for _, data := range response.Data.VolumeList {
		fmt.Println("\nstart download volume: ", data.Title)
		catalogue.ChapterBar = New(len(data.ChapterList))
		for _, Chapter := range data.ChapterList {
			if Chapter.OriginNeedFireMoney == 0 {

				catalogue.SfacgContent(strconv.Itoa(Chapter.ChapID))

			} else {
				fmt.Println("this chapter is VIP and need fire money, skip it")
			}
		}
	}
	return true
}

func (catalogue *Catalogue) SfacgContent(ChapterId string) {
	if err := catalogue.ChapterBar.Add(1); err != nil {
		fmt.Println("bar error:", err)
	} else {
		cfg.DelayTime()
	}
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
	writeContent := fmt.Sprintf("%v:%v\n%v\n%v\n\n",
		response.Data.Title, response.Data.AddTime, response.Data.Expand.Content, cfg.Vars.BookInfo.AuthorName,
	)
	cfg.EncapsulationWrite(
		path.Join(cfg.Vars.SaveFile, cfg.Vars.BookInfo.NovelName+".txt"), writeContent, 5, "a",
	)
}

func (catalogue *Catalogue) CatCatalogue() bool {
	for index, division := range HbookerAPI.GetDivisionIdByBookId(cfg.Vars.BookInfo.NovelID) {
		fmt.Println("index:", index, "\t\tdivision:", division.DivisionName)
		for _, chapter := range HbookerAPI.GetCatalogueByDivisionId(division.DivisionID) {
			if chapter.IsValid == "1" {
				catalogue.ChapterList = append(catalogue.ChapterList, chapter)
			}
		}
	}
	catalogue.ChapterBar = New(len(catalogue.ChapterList))
	for _, chapter := range catalogue.ChapterList {
		catalogue.CatContent(chapter.ChapterID)
	}
	return true
}

func (catalogue *Catalogue) CatContent(ChapterId string) {
	if err := catalogue.ChapterBar.Add(1); err != nil {
		fmt.Println("bar error:", err)
	} else {
		cfg.DelayTime()
	}
	response := HbookerAPI.GetContent(ChapterId)
	writeContent, SavePath := fmt.Sprintf("%v:%v\n%v\n%v\n\n",
		response.ChapterTitle,
		response.Uptime,
		response.TxtContent,
		cfg.Vars.BookInfo.AuthorName,
	), path.Join(cfg.Vars.SaveFile, cfg.Vars.BookInfo.NovelName+".txt")
	cfg.EncapsulationWrite(SavePath, writeContent, 5, "a")

}

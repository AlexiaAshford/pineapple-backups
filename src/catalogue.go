package src

import (
	"fmt"
	"math/rand"
	"os"
	"path"
	"sf/cfg"
	"sf/src/boluobao"
	"sf/src/hbooker"
	structs "sf/structural/hbooker_structs"
	"strconv"
	"time"
)

type Catalogue struct {
	ChapterBar  *ProgressBar
	ChapterList []structs.ChapterList
}

func (is *Catalogue) SfacgCatalogue() bool {
	response := boluobao.GetCatalogueDetailedById(cfg.Vars.BookInfo.NovelID)
	for _, data := range response.Data.VolumeList {
		fmt.Println("\nstart download volume: ", data.Title)
		is.ChapterBar = New(len(data.ChapterList))
		for _, Chapter := range data.ChapterList {
			if Chapter.OriginNeedFireMoney == 0 {
				is.SfacgContent(strconv.Itoa(Chapter.ChapID))
			} else {
				fmt.Println("this chapter is VIP and need fire money, skip it")
			}
		}
	}
	return true
}

func (is *Catalogue) SfacgContent(ChapterId string) {
	if err := is.ChapterBar.Add(1); err != nil {
		fmt.Println("bar error:", err)
	} else {
		time.Sleep(time.Second * time.Duration(rand.Intn(4)))
	}
	response := boluobao.GetContentDetailedByCid(ChapterId)
	if response.Status.HTTPCode != 200 {
		if response.Status.Msg == "接口校验失败,请尽快把APP升级到最新版哦~" {
			fmt.Println(response.Status.Msg)
			os.Exit(0)
		} else {
			fmt.Println(response.Status.Msg)
			is.SfacgContent(ChapterId)
		}
	} else {
		writeContent, SavePath := fmt.Sprintf("%v:%v\n%v\n%v\n\n\n",
			response.Data.Title,
			response.Data.AddTime,
			response.Data.Expand.Content,
			cfg.Vars.BookInfo.AuthorName,
		), path.Join(cfg.Vars.SaveFile, cfg.Vars.BookInfo.NovelName+".txt")
		cfg.EncapsulationWrite(SavePath, writeContent, 5, 0666)
	}
}

func (is *Catalogue) CatCatalogue() bool {
	for index, division := range HbookerAPI.GetDivisionIdByBookId(cfg.Vars.BookInfo.NovelID) {
		fmt.Println("index:", index, "\t\tdivision:", division.DivisionName)
		for _, chapter := range HbookerAPI.GetCatalogueByDivisionId(division.DivisionID) {
			if chapter.IsValid == "1" {
				is.ChapterList = append(is.ChapterList, chapter)
			}
		}
	}
	is.ChapterBar = New(len(is.ChapterList))
	for _, chapter := range is.ChapterList {
		is.CatContent(chapter.ChapterID)
	}
	return true
}

func (is *Catalogue) CatContent(ChapterId string) {
	if err := is.ChapterBar.Add(1); err != nil {
		fmt.Println("bar error:", err)
	} else {
		time.Sleep(time.Second * time.Duration(rand.Intn(4)))
	}
	response := HbookerAPI.GetContent(ChapterId)
	writeContent, SavePath := fmt.Sprintf("%v:%v\n%v\n%v\n\n\n",
		response.ChapterTitle,
		response.Uptime,
		response.TxtContent,
		cfg.Vars.BookInfo.AuthorName,
	), path.Join(cfg.Vars.SaveFile, cfg.Vars.BookInfo.NovelName+".txt")
	cfg.EncapsulationWrite(SavePath, writeContent, 5, 0666)

}

package src

import (
	"fmt"
	"math/rand"
	"os"
	"path"
	"sf/cfg"
	"sf/src/boluobao"
	"strconv"
	"time"
)

func SfacgCatalogue() bool {
	response := boluobao.GetCatalogueDetailedById(cfg.Vars.BookInfo.NovelID)
	for _, data := range response.Data.VolumeList {
		fmt.Println("\nstart download volume: ", data.Title)
		bar := New(len(data.ChapterList))
		for _, Chapter := range data.ChapterList {
			if Chapter.OriginNeedFireMoney == 0 {
				SfacgContent(len(data.ChapterList), strconv.Itoa(Chapter.ChapID), bar)
			} else {
				fmt.Println("this chapter is VIP and need fire money, skip it")
			}
		}
	}
	return true
}

func SfacgContent(ChapLength int, ChapterId string, bar *ProgressBar) {
	if err := bar.Add(1); err != nil {
		fmt.Println("bar error:", err)
	} else {
		time.Sleep(time.Second * time.Duration(rand.Intn(5)))
	}
	response := boluobao.GetContentDetailedByCid(ChapterId)
	if response.Status.HTTPCode != 200 {
		if response.Status.Msg == "接口校验失败,请尽快把APP升级到最新版哦~" {
			fmt.Println(response.Status.Msg)
			os.Exit(0)
		} else {
			fmt.Println(response.Status.Msg)
			SfacgContent(ChapLength, ChapterId, bar)
		}
	} else {

		writeContent, SavePath := fmt.Sprintf("%v:%v\n%v\n%v\n\n\n",
			response.Data.Title,
			response.Data.AddTime,
			response.Data.Expand.Content,
			cfg.Vars.BookInfo.AuthorName,
		), path.Join(cfg.Vars.SaveFile, cfg.Vars.BookInfo.NovelName+".txt")

		for i := 0; i < 5; i++ {
			if cfg.WriteFile(SavePath, writeContent, 0666) == nil {
				break
			} else {
				fmt.Println("write file error, try again", i)
			}
		}
	}
}

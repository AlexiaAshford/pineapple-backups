package app

import (
	"fmt"
	"os"
	"path"

	"github.com/AlexiaVeronica/boluobaoLib/boluobaomodel"
	"github.com/AlexiaVeronica/hbookerLib/hbookermodel"
	"github.com/AlexiaVeronica/pineapple-backups/pkg/file"
)

func sfContinueFunction(chapter boluobaomodel.ChapterList) bool {
	savePath := path.Join("cache", BoluobaoLibAPP, fmt.Sprintf("%v.txt", chapter.ChapID))
	if _, err := os.Stat(savePath); err == nil {
		return false
	} else if os.IsNotExist(err) && chapter.OriginNeedFireMoney == 0 {
		return true
	}
	return false
}

func sfContentFunction(chapter *boluobaomodel.ContentData) {
	fmt.Println("========", chapter.Title)
	savePath := path.Join("cache", BoluobaoLibAPP, fmt.Sprintf("%v.txt", chapter.ChapId))
	file.Open(savePath, "\n\n"+chapter.Title+"\n"+chapter.Expand.Content, "w")
}

func hbookerContinueFunction(chapter hbookermodel.ChapterList) bool {
	savePath := path.Join("cache", CiweimaoLibAPP, fmt.Sprintf("%v.txt", chapter.ChapterID))
	if _, err := os.Stat(savePath); err == nil {
		return false
	} else if os.IsNotExist(err) && chapter.AuthAccess == "1" {
		return true
	}
	return true
}

func hbookerContentFunction(chapter *hbookermodel.ChapterInfo) {
	savePath := path.Join("cache", CiweimaoLibAPP, fmt.Sprintf("%v.txt", chapter.ChapterID))
	file.Open(savePath, "\n\n"+chapter.ChapterTitle+"\n"+chapter.TxtContent, "w")
}

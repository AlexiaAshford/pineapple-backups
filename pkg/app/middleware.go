package app

import (
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/AlexiaVeronica/boluobaoLib/boluobaomodel"
	"github.com/AlexiaVeronica/hbookerLib/hbookermodel"
)

func sfContinueFunction(chapter boluobaomodel.ChapterList) bool {
	return shouldContinue(BoluobaoLibAPP, strconv.Itoa(chapter.ChapID), chapter.OriginNeedFireMoney == 0)
}

func sfContentFunction(chapter *boluobaomodel.ContentData) {
	writeContentToFile(BoluobaoLibAPP, strconv.Itoa(chapter.ChapId), chapter.Title, chapter.Expand.Content)
}

func hbookerContinueFunction(chapter hbookermodel.ChapterList) bool {
	return shouldContinue(CiweimaoLibAPP, chapter.ChapterID, chapter.AuthAccess == "1")
}

func hbookerContentFunction(chapter *hbookermodel.ChapterInfo) {
	writeContentToFile(CiweimaoLibAPP, chapter.ChapterID, chapter.ChapterTitle, chapter.TxtContent)
}

func shouldContinue(app, chapID string, condition bool) bool {
	savePath := path.Join("cache", app, fmt.Sprintf("%v.txt", chapID))
	if _, err := os.Stat(savePath); err == nil {
		return false
	} else if os.IsNotExist(err) && condition {
		return true
	}
	return false
}

func writeContentToFile(app, chapID, title, content string) {
	savePath := path.Join("cache", app, fmt.Sprintf("%v.txt", chapID))
	file, err := os.Create(savePath)
	if err != nil {
		fmt.Printf("Failed to create file %s: %v\n", savePath, err)
		return
	}
	defer file.Close()
	file.WriteString("\n\n" + title + "\n" + content)
}

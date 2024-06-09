package app

import (
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/AlexiaVeronica/boluobaoLib/boluobaomodel"
	"github.com/AlexiaVeronica/hbookerLib/hbookermodel"
)

func sfContinueFunction(bookInfo *boluobaomodel.BookInfoData, chapter boluobaomodel.ChapterList) bool {
	return shouldContinue(BoluobaoLibAPP, strconv.Itoa(chapter.ChapID), chapter.OriginNeedFireMoney == 0)
}

func sfContentFunction(bookInfo *boluobaomodel.BookInfoData, chapter *boluobaomodel.ContentData) {
	fmt.Println(chapter.Title + "  下载完毕")
	writeContentToFile(BoluobaoLibAPP, strconv.Itoa(chapter.ChapId), chapter.Title, chapter.Expand.Content)
}
func sfMergeText(file *os.File) func(chapter boluobaomodel.ChapterList) {
	return func(chapter boluobaomodel.ChapterList) {
		savePath := path.Join("cache", BoluobaoLibAPP, fmt.Sprintf("%v.txt", chapter.ChapID))
		if _, err := os.Stat(savePath); err == nil {
			content, ok := os.ReadFile(savePath)
			if ok == nil {
				file.WriteString(fmt.Sprintf("\n\n%s\n%s", chapter.Title, content))
			}
		}
	}
}
func hbookerMergeText(file *os.File) func(chapter hbookermodel.ChapterList) {
	return func(chapter hbookermodel.ChapterList) {
		savePath := path.Join("cache", BoluobaoLibAPP, fmt.Sprintf("%v.txt", chapter.ChapterID))
		if _, err := os.Stat(savePath); err == nil {
			content, ok := os.ReadFile(savePath)
			if ok == nil {
				file.WriteString(fmt.Sprintf("\n\n%s\n%s", chapter.ChapterTitle, content))
			}
		}
	}
}

func hbookerContinueFunction(bookInfo *hbookermodel.BookInfo, chapter hbookermodel.ChapterList) bool {
	return shouldContinue(CiweimaoLibAPP, chapter.ChapterID, chapter.AuthAccess == "1")
}

func hbookerContentFunction(bookInfo *hbookermodel.BookInfo, chapter *hbookermodel.ChapterInfo) {
	writeContentToFile(CiweimaoLibAPP, chapter.ChapterID, chapter.ChapterTitle, chapter.TxtContent)
}

func shouldContinue(app, chapID string, condition bool) bool {
	savePath := path.Join("cache", app, fmt.Sprintf("%v.txt", chapID))
	if _, err := os.Stat(savePath); err == nil || (os.IsNotExist(err) && !condition) {
		return false
	}
	return true
}

func writeContentToFile(app, chapID, title, content string) {
	file, err := os.Create(path.Join("cache", app, fmt.Sprintf("%v.txt", chapID)))
	if err != nil {
		fmt.Printf("writeContentToFile file %s: %v\n", title, err)
		return
	}
	defer file.Close()

	file.WriteString(fmt.Sprintf("\n\n%s\n%s", title, content))
}

package app

import (
	"fmt"
	"github.com/AlexiaVeronica/boluobaoLib/boluobaomodel"
	"github.com/AlexiaVeronica/hbookerLib/hbookermodel"
	"os"
	"path"

	"github.com/AlexiaVeronica/boluobaoLib"
	"github.com/AlexiaVeronica/hbookerLib"
)

type APP struct {
	CurrentApp string
	BookId     string
	Boluobao   *boluobaoLib.Client
	Ciweimao   *hbookerLib.Client
}

const (
	BoluobaoLibAPP = "boluobao"
	CiweimaoLibAPP = "ciweimao"
)

func NewApp() *APP {
	a := &APP{
		CurrentApp: BoluobaoLibAPP,
		Boluobao:   boluobaoLib.NewClient(),
		Ciweimao:   hbookerLib.NewClient(),
	}
	a.initDirectories()
	return a
}

func (a *APP) initDirectories() {
	directories := []string{
		"cache",
		"cache/" + CiweimaoLibAPP,
		"cache/" + BoluobaoLibAPP,
		"save/",
		"cover/",
		"save/" + CiweimaoLibAPP,
		"save/" + BoluobaoLibAPP,
	}
	for _, dir := range directories {
		createCacheDirectory(dir)
	}
}
func mergeTextToFile[T any](bookName string) *os.File {
	data := new(T)
	var savePath string
	switch any(data).(type) {
	case *boluobaomodel.ChapterList:
		savePath = path.Join("save", BoluobaoLibAPP, fmt.Sprintf("%v.txt", bookName))
	case *hbookermodel.ChapterList:
		savePath = path.Join("save", CiweimaoLibAPP, fmt.Sprintf("%v.txt", bookName))
	}
	file, err := os.OpenFile(savePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Printf("Failed to create file %s: %v\n", savePath, err)
		return nil
	}
	return file
}

func createCacheDirectory(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err = os.MkdirAll(path, 0777); err != nil {
			fmt.Printf("Failed to create directory %s: %v\n", path, err)
		}
	}
}

func (a *APP) SetCurrentApp(app string) {
	switch app {
	case BoluobaoLibAPP, CiweimaoLibAPP:
		a.CurrentApp = app
	default:
		panic("Invalid app type: " + app)
	}
}

func (a *APP) GetCurrentApp() string {
	return a.CurrentApp
}

func (a *APP) SearchDetailed(keyword string) *APP {
	switch a.CurrentApp {
	case BoluobaoLibAPP:
		newApp := a.Boluobao.APP()
		newApp.Search(keyword, sfContinueFunction, sfContentFunction).MergeText(
			sfMergeText(mergeTextToFile[boluobaomodel.ChapterList](newApp.GetBookInfo().NovelName)))
	case CiweimaoLibAPP:
		newApp := a.Ciweimao.APP().Search(keyword, hbookerContinueFunction, hbookerContentFunction)
		newApp.MergeText(hbookerMergeText(mergeTextToFile[hbookermodel.ChapterList](newApp.GetBookInfo().BookName)))
	}
	return a
}

func (a *APP) DownloadBookByBookId(bookId string) *APP {
	switch a.CurrentApp {
	case BoluobaoLibAPP:
		bookInfo, err := a.Boluobao.API().GetBookInfo(bookId)
		if err != nil {
			fmt.Println("Failed to get book info:", err)
			return a
		}
		a.Boluobao.APP().SetBookInfo(&bookInfo.Data).Download(sfContinueFunction, sfContentFunction).MergeText(
			sfMergeText(mergeTextToFile[boluobaomodel.ChapterList](bookInfo.Data.NovelName)))
	case CiweimaoLibAPP:
		bookInfo, err := a.Ciweimao.API().GetBookInfo(bookId)
		if err != nil {
			fmt.Println("Failed to get book info:", err)
			return a
		}
		a.Ciweimao.APP().SetBookInfo(&bookInfo.Data.BookInfo).Download(hbookerContinueFunction, hbookerContentFunction).
			MergeText(hbookerMergeText(mergeTextToFile[hbookermodel.ChapterList](bookInfo.Data.BookInfo.BookName)))
	}
	return a
}

func (a *APP) Bookshelf() *APP {
	switch a.CurrentApp {
	case BoluobaoLibAPP:
		newApp := a.Boluobao.APP().Bookshelf(sfContinueFunction, sfContentFunction)
		newApp.MergeText(sfMergeText(mergeTextToFile[boluobaomodel.ChapterList](newApp.GetBookInfo().NovelName)))
	case CiweimaoLibAPP:
		newApp := a.Ciweimao.APP().Bookshelf(hbookerContinueFunction, hbookerContentFunction)
		newApp.MergeText(hbookerMergeText(mergeTextToFile[hbookermodel.ChapterList](newApp.GetBookInfo().BookName)))
	}
	return a
}

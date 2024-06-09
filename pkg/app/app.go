package app

import (
	"fmt"
	"github.com/AlexiaVeronica/pineapple-backups/pkg/file"
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
	createCacheDirectory("cache")
	createCacheDirectory("cache/" + CiweimaoLibAPP)
	createCacheDirectory("cache/" + BoluobaoLibAPP)
	return a
}

func createCacheDirectory(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0777)
	}
}

func (a *APP) SetCurrentApp(app string) {
	switch app {
	case BoluobaoLibAPP, CiweimaoLibAPP:
		a.CurrentApp = app
	default:
		panic("app type is not correct: " + app)
	}
}

func (a *APP) GetCurrentApp() string {
	return a.CurrentApp
}

func (a *APP) SearchDetailed(keyword string) *APP {
	switch a.CurrentApp {
	case BoluobaoLibAPP:
		a.Boluobao.APP().Search(keyword, sfContinueFunction, sfContentFunction)
	case CiweimaoLibAPP:
		a.Ciweimao.APP().Search(keyword, hbookerContinueFunction, hbookerContentFunction)
	}
	return a
}

func (a *APP) DownloadBookByBookId(bookId string) *APP {
	switch a.CurrentApp {
	case BoluobaoLibAPP:
		a.Boluobao.APP().Download(bookId, sfContinueFunction, sfContentFunction)
	case CiweimaoLibAPP:
		a.Ciweimao.APP().Download(bookId, hbookerContinueFunction, hbookerContentFunction)
	}
	return a
}

func (a *APP) Bookshelf() *APP {
	switch a.CurrentApp {
	case BoluobaoLibAPP:
		a.Boluobao.APP().Bookshelf(sfContinueFunction, sfContentFunction)
	case CiweimaoLibAPP:
		a.Ciweimao.APP().Bookshelf(hbookerContinueFunction, hbookerContentFunction)
	}
	return a
}

func (a *APP) MergingText() *APP {
	switch a.CurrentApp {
	case BoluobaoLibAPP:
		fmt.Println("boluobao app does not support merging text")
		chapterList, err := a.Boluobao.API().GetCatalogue(a.BookId)
		if err != nil {
			fmt.Println("get chapter list error:", err)
			return a
		}
		for _, chapter := range chapterList.Data.VolumeList {
			for _, chapterInfo := range chapter.ChapterList {
				cachePath := path.Join("cache", BoluobaoLibAPP, fmt.Sprintf("%v.txt", chapterInfo.ChapID))
				savePath := path.Join("save", fmt.Sprintf("%v.txt", a.BookId))
				if _, err = os.Stat(cachePath); err == nil {
					file.Open(savePath, cachePath, "a")
				}
			}
		}
	case CiweimaoLibAPP:
		chapterList, err := a.Ciweimao.API().GetDivisionListByBookId(a.BookId)
		if err != nil {
			fmt.Println("get division list error:", err)
			return a
		}
		for _, division := range chapterList.Data.ChapterList {
			for _, chapter := range division.ChapterList {
				cachePath := path.Join("cache", CiweimaoLibAPP, fmt.Sprintf("%v.txt", chapter.ChapterID))
				savePath := path.Join("save", fmt.Sprintf("%v.txt", a.BookId))
				if _, err = os.Stat(cachePath); err == nil {
					file.Open(savePath, cachePath, "a")
				}
			}
		}
	}
	return a
}

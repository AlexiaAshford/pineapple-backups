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

func createCacheDirectory(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0777); err != nil {
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
		a.Boluobao.APP().Search(keyword, sfContinueFunction, sfContentFunction)
	case CiweimaoLibAPP:
		a.Ciweimao.APP().Search(keyword, hbookerContinueFunction, hbookerContentFunction)
	}
	return a
}

func (a *APP) DownloadBookByBookId(bookId string) *APP {
	fmt.Println("bookId", bookId)
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
		a.mergeTextBoluobao()
	case CiweimaoLibAPP:
		a.mergeTextCiweimao()
	}
	return a
}

func (a *APP) mergeTextBoluobao() {
	bookInfo, err := a.Boluobao.API().GetBookInfo(a.BookId)
	if err != nil {
		fmt.Println("Failed to get book info:", err)
		return
	}
	chapterList, err := a.Boluobao.API().GetCatalogue(bookInfo.Data.NovelId)
	if err != nil {
		fmt.Println("Failed to get chapter list:", err)
		return
	}
	a.writeChaptersToFile(bookInfo.Data.NovelName, chapterList.Data.VolumeList)
}

func (a *APP) mergeTextCiweimao() {
	bookInfo, err := a.Ciweimao.API().GetBookInfo(a.BookId)
	if err != nil {
		fmt.Println("Failed to get book info:", err)
		return
	}
	chapterList, err := a.Ciweimao.API().GetDivisionListByBookId(bookInfo.Data.BookInfo.BookID)
	if err != nil {
		fmt.Println("Failed to get division list:", err)
		return
	}
	a.writeChaptersToFile(bookInfo.Data.BookInfo.BookName, chapterList.Data.ChapterList)
}

func (a *APP) writeChaptersToFile(bookName string, chapterList interface{}) {
	filePath := path.Join("save", a.CurrentApp, fmt.Sprintf("%v.txt", bookName))
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Failed to create file %s: %v\n", filePath, err)
		return
	}
	defer file.Close()

	chapters := extractChapters(chapterList)
	for _, chapterPath := range chapters {
		if _, err = os.Stat(chapterPath); err == nil {
			content, err := os.ReadFile(chapterPath)
			if err == nil {
				file.Write(content)
			}
		}
	}
}

func extractChapters(chapterList interface{}) []string {
	var chapters []string
	switch list := chapterList.(type) {
	case []boluobaomodel.VolumeList:
		for _, volume := range list {
			for _, chapter := range volume.ChapterList {
				chapters = append(chapters, path.Join("cache", BoluobaoLibAPP, fmt.Sprintf("%v.txt", chapter.ChapID)))
			}
		}
	case []hbookermodel.DivisionList:
		for _, division := range list {
			for _, chapter := range division.ChapterList {
				chapters = append(chapters, path.Join("cache", CiweimaoLibAPP, fmt.Sprintf("%v.txt", chapter.ChapterID)))
			}
		}
	}

	return chapters
}

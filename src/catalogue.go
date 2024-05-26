package src

import (
	"fmt"
	"github.com/AlexiaVeronica/pineapple-backups/config"
	"github.com/AlexiaVeronica/pineapple-backups/pkg/command"
	"github.com/AlexiaVeronica/pineapple-backups/pkg/epub"
	"github.com/AlexiaVeronica/pineapple-backups/pkg/file"
	"github.com/AlexiaVeronica/pineapple-backups/pkg/threading"
	"github.com/AlexiaVeronica/pineapple-backups/pkg/tools"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type Catalogue struct {
	ChapterBar  *ProgressBar
	ChapterCfg  []string
	EpubSetting *epub.Epub
}

func (catalogue *Catalogue) ReadChapterConfig() {
	var bookName string
	if command.Command.AppType == "cat" {
		bookName = config.APP.Hbooker.BookInfo.BookName
	} else {
		bookName = config.APP.SFacg.BookInfo.NovelName
	}
	if !config.Exist(path.Join("cache", bookName)) {
		os.Mkdir(path.Join("cache", bookName), 0777)
		catalogue.ChapterCfg = []string{}
	} else {
		catalogue.ChapterCfg = tools.GetFileName("cache")
	}
}

func (catalogue *Catalogue) GetDownloadsList() ([]string, error) {
	catalogue.ReadChapterConfig()
	var downloadList []string
	switch command.Command.AppType {
	case "sfacg":
		divisionList, err := config.APP.SFacg.Client.API().GetCatalogue(config.APP.SFacg.BookInfo.NovelId)
		if err != nil {
			return nil, err
		}
		for _, division := range divisionList {
			for _, chapter := range division.ChapterList {
				if !tools.TestList(catalogue.ChapterCfg, strconv.Itoa(chapter.ChapID)) {
					if chapter.OriginNeedFireMoney == 0 {
						downloadList = append(downloadList, strconv.Itoa(chapter.ChapID))
					}
				}
			}
		}
	case "cat":
		divisionList, err := config.APP.Hbooker.Client.API.GetDivisionListByBookId(config.APP.Hbooker.BookInfo.BookID)
		if err != nil {
			return nil, err
		}
		for _, value := range divisionList {
			for _, value2 := range value.ChapterList {
				if !tools.TestList(catalogue.ChapterCfg, value2.ChapterID) {
					if value2.AuthAccess == "0" || value2.AuthAccess == "1" {
						downloadList = append(downloadList, value2.ChapterID)
					} else {
						fmt.Println(value2.ChapterTitle, " is vip chapter, You need to subscribe it")
					}
				}
			}
		}
	}
	if len(downloadList) == 0 {
		return nil, fmt.Errorf("no new chapters")
	}
	return downloadList, nil
}

func (catalogue *Catalogue) DownloadContent(threading *threading.GoLimit, chapterID string) {
	defer threading.Done()
	catalogue.speed_progress()

	if config.Exist(path.Join(config.Vars.ConfigName, config.APP.SFacg.BookInfo.NovelName, chapterID+".txt")) {
		return
	}
	var contentText string
	if command.Command.AppType == "sfacg" {
		chapterInfo, err := config.APP.SFacg.Client.API().GetChapterContent(chapterID)
		if err != nil {
			fmt.Println("get chapter content error:", err)
			return
		}
		contentText = chapterInfo.Expand.Content
	} else if command.Command.AppType == "cat" {
		chapterKey, err := config.APP.Hbooker.Client.API.GetChapterKey(chapterID)
		if err != nil {
			return
		}
		chapterInfo, err := config.APP.Hbooker.Client.API.GetChapterContentAPI(chapterID, chapterKey)
		if err != nil {
			fmt.Println("get chapter content error:", err)
			return
		}
		contentText = chapterInfo.TxtContent
	}
	if contentText != "" {
		file.Open(path.Join("cache", config.APP.SFacg.BookInfo.NovelName, chapterID+".txt"), contentText, "w")
	}
}

func (catalogue *Catalogue) MergeTextAndEpubFiles() {
	var savePath string
	if command.Command.AppType == "sfacg" {
		savePath = path.Join(config.Vars.OutputName, config.APP.SFacg.BookInfo.NovelName, config.APP.SFacg.BookInfo.NovelName)
		divisionList, err := config.APP.SFacg.Client.API().GetCatalogue(config.APP.SFacg.BookInfo.NovelId)
		if err != nil {
			fmt.Println("get division list error:", err)
			return
		}
		for _, i := range divisionList {
			for _, chapterInfo := range i.ChapterList {
				if config.Exist(path.Join(config.Vars.ConfigName, config.APP.SFacg.BookInfo.NovelName, strconv.Itoa(chapterInfo.ChapID)+".txt")) {
					content := file.Open(path.Join(config.Vars.ConfigName, config.APP.SFacg.BookInfo.NovelName, strconv.Itoa(chapterInfo.ChapID)+".txt"), "", "r")
					file.Open(savePath+".txt", "\n\n\n"+chapterInfo.Title+"\n"+content, "a")
					catalogue.addChapterInEpubFile(chapterInfo.Title, content+"</p>")
				}
			}
		}
	} else {
		savePath = path.Join(config.Vars.OutputName, config.APP.Hbooker.BookInfo.BookName, config.APP.Hbooker.BookInfo.BookName)
		divisionList, err := config.APP.Hbooker.Client.API.GetDivisionListByBookId(config.APP.Hbooker.BookInfo.BookID)
		if err != nil {
			fmt.Println("get division list error:", err)
			return
		}
		for _, i := range divisionList {
			for _, chapterInfo := range i.ChapterList {
				if config.Exist(path.Join(config.Vars.ConfigName, config.APP.SFacg.BookInfo.NovelName, chapterInfo.ChapterID+".txt")) {
					content := file.Open(path.Join(config.Vars.ConfigName, config.APP.SFacg.BookInfo.NovelName, chapterInfo.ChapterID+".txt"), "", "r")
					file.Open(savePath+".txt", "\n\n\n"+chapterInfo.ChapterTitle+"\n"+content, "a")
					catalogue.addChapterInEpubFile(chapterInfo.ChapterTitle, content+"</p>")
				}
			}
		}
	}

	outPutEpubNow := time.Now() // start time
	if err := catalogue.EpubSetting.Write(savePath + ".epub"); err != nil {
		fmt.Println("output epub error:", err)
	}
	fmt.Println("output epub file success, time:", time.Since(outPutEpubNow))
}

func (catalogue *Catalogue) addChapterInEpubFile(title, content string) {
	xmlContent := "<h1>" + title + "</h1>\n<p>" + strings.ReplaceAll(content, "\n", "</p>\n<p>")
	if _, err := catalogue.EpubSetting.AddSection(xmlContent, title, "", ""); err != nil {
		fmt.Printf("epub add chapter:%v\t\terror message:%v", title, err)
	}
}
func (catalogue *Catalogue) speed_progress() {
	if err := catalogue.ChapterBar.Add(1); err != nil {
		fmt.Println("bar error:", err)
	} else {
		//time.Sleep(time.Second * time.Duration(2))
	}

}

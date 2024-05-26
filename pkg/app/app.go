package app

import (
	"fmt"
	"github.com/AlexiaVeronica/boluobaoLib"
	"github.com/AlexiaVeronica/pineapple-backups/config"
	"strconv"
)

type APP struct {
	CurrentApp   string
	searchResult []Book
	Boluobao     *boluobaoLib.Client

	//Ciweimao * .Client
}

type Book struct {
	currentApp string
	index      int
	BookId     string
	BookName   string
	Author     string
}

const boluobaoLibAPP = "boluobao"
const ciweimaoLibAPP = "ciweimao"

func NewApp() *APP {
	return &APP{

		CurrentApp:   boluobaoLibAPP,
		searchResult: make([]Book, 0),
		Boluobao:     boluobaoLib.NewClient(),
		//Ciweimao: .NewClient(),
	}
}

func (a *APP) SetCurrentApp(app string) {
	if app == boluobaoLibAPP {
		a.CurrentApp = boluobaoLibAPP
	} else if app == ciweimaoLibAPP {
		a.CurrentApp = ciweimaoLibAPP
	} else {
		panic("app type is not correct" + app)
	}
}
func (a *APP) GetCurrentApp() string {
	return a.CurrentApp
}

func (a *APP) SearchDetailed(keyword string) *APP {
	if a.CurrentApp == boluobaoLibAPP {
		searchInfo, err := config.APP.SFacg.Client.API().GetSearch(keyword, 0)
		if err != nil {
			fmt.Println("search failed! " + err.Error())
			return nil
		}
		for index, book := range searchInfo {
			fmt.Printf("Index: %d\t\t\tBookName: %s\n", index, book.NovelName)
			a.searchResult = append(a.searchResult, Book{
				index:      index,
				BookId:     strconv.Itoa(book.NovelId),
				BookName:   book.NovelName,
				Author:     book.AuthorName,
				currentApp: a.CurrentApp,
			})
		}
	} else if a.CurrentApp == ciweimaoLibAPP {

	}
	return a
}
func (book *Book) DownloadBookByBookId(client any, bookId string) *Book {
	if book.currentApp == boluobaoLibAPP {
		bookInfo, err := client.(*boluobaoLib.Client).API().GetBookInfo(bookId)
		if err != nil {
			fmt.Println(err)
		} else {
			book.BookId = strconv.Itoa(bookInfo.NovelId)
			book.BookName = bookInfo.NovelName
			book.Author = bookInfo.AuthorName
			book.currentApp = boluobaoLibAPP
		}

	} else if book.currentApp == ciweimaoLibAPP {

	}
	return book
}

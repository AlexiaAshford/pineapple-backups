package src

import (
	"fmt"
	"github.com/VeronicaAlexia/pineapple-backups/config"
	"github.com/VeronicaAlexia/pineapple-backups/pkg/command"
	"github.com/VeronicaAlexia/pineapple-backups/pkg/tools"
	"strconv"
)

type Search struct {
	Keyword      string
	Page         int
	SearchResult []string
}

func (s *Search) CatSearchDetailed() []string {
	var searchResult []string
	searchInfo, err := config.APP.Hbooker.Client.API.GetSearchBooksAPI(s.Keyword, s.Page)
	if err != nil {
		fmt.Println("search failed!" + err.Error())
		return nil
	}
	for index, bookInfo := range searchInfo {
		fmt.Println("Index:", index, "\t\t\tBookName:", bookInfo.BookName)
		searchResult = append(searchResult, bookInfo.BookID)
	}
	return searchResult
}

func (s *Search) SfacgSearchDetailed() []string {
	var searchResult []string
	searchInfo, err := config.APP.SFacg.Client.API.GetSearch(s.Keyword, s.Page)
	if err != nil {
		fmt.Println("search failed!" + err.Error())
		return nil // if the search result is empty
	}
	for index, book := range searchInfo {
		fmt.Println("Index:", index, "\t\t\tBookName:", book.NovelName)
		searchResult = append(searchResult, strconv.Itoa(book.NovelId))
	}
	return searchResult
}

func (s *Search) load_search_list() {
	if command.Command.AppType == "cat" {
		s.SearchResult = s.CatSearchDetailed()
	} else if command.Command.AppType == "sfacg" {
		s.SearchResult = s.SfacgSearchDetailed()
	} else {
		panic("app type is not correct" + command.Command.AppType)
	}
}

func (s *Search) subtraction() {
	if s.Page > 0 {
		s.Page -= 1 // previous page
	} else {
		fmt.Println("page is 0, cannot go previous")
	}
	s.load_search_list()
}

func (s *Search) add() {
	s.Page += 1 // next page
	s.load_search_list()
}

func (s *Search) SearchBook() string {
	s.load_search_list()
	for {
		keyword := tools.InputStr("lease input search keyword:")
		if keyword == "next" || keyword == "n" {
			s.add()
		} else if keyword == "previous" || keyword == "p" {
			s.subtraction()
		}
		if tools.IsNum(keyword) && tools.StrToInt(keyword) < len(s.SearchResult) {
			if BookID := s.SearchResult[tools.StrToInt(keyword)]; BookID != "" {
				return BookID // if the input is a number (book id)
			} else {
				fmt.Println("No found search book, please input again:")
			}
		}
	}
}

//func (s *Search) ReturnBookID(keyword string) string {
//	if tool.IsNum(keyword) {
//		return s.SearchResult[tool.InputInt(keyword, len(s.SearchResult))]
//	} else { // if the input is not a number (not a book id)
//		fmt.Println("input is not a number, please input search index or book id")
//	}
//	return "" // if the input is not a number (not a book id) or the index is out of range
//}

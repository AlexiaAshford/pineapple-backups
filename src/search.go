package src

import (
	"fmt"
	"github.com/VeronicaAlexia/BoluobaoAPI/boluobao/search"
	HbookerAPI "github.com/VeronicaAlexia/HbookerAPI/ciweimao/search"
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
	response := HbookerAPI.GET_SEARCH(s.Keyword, s.Page) // init search struct
	if response.Code != "100000" || len(response.Data.BookList) == 0 {
		fmt.Println("search failed, code:", response.Code)
		return nil
	} else {
		fmt.Println("this page has", len(response.Data.BookList), "books")
	}
	for index, book_info := range response.Data.BookList {
		fmt.Println("Index:", index, "\t\t\tBookName:", book_info.BookName)
		searchResult = append(searchResult, book_info.BookID)
	}
	return searchResult
}

func (s *Search) SfacgSearchDetailed() []string {
	var searchResult []string
	response := search.GET_SEARCH(s.Keyword, s.Page)
	if response.Status.HTTPCode != 200 || len(response.Data.Novels) == 0 {
		return nil // if the search result is empty
	} else {
		fmt.Println("this page has", len(response.Data.Novels), "novels")
	}
	for index, book := range response.Data.Novels {
		fmt.Println("Index:", index, "\t\t\tBookName:", book.NovelName)
		searchResult = append(searchResult, strconv.Itoa(book.NovelID))
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

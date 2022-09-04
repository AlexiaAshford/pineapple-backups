package app

import (
	"fmt"
	"github.com/VeronicaAlexia/pineapple-backups/config"
	"github.com/VeronicaAlexia/pineapple-backups/config/tool"
	"github.com/VeronicaAlexia/pineapple-backups/src/app/boluobao"
	"github.com/VeronicaAlexia/pineapple-backups/src/app/hbooker"
	"strconv"
)

type Search struct {
	Keyword      string
	Page         int
	SearchResult []string
}

func (s *Search) CatSearchDetailed() []string {
	var searchResult []string
	response := hbooker.GET_SEARCH(s.Keyword, s.Page)
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
	response := boluobao.GET_SEARCH(s.Keyword, s.Page)
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

func (s *Search) TestApp() {
	if config.Vars.AppType == "cat" {
		s.SearchResult = s.CatSearchDetailed()
	} else if config.Vars.AppType == "sfacg" {
		s.SearchResult = s.SfacgSearchDetailed()
	} else {
		panic("app type is not correct" + config.Vars.AppType)
	}
}

func (s *Search) SearchBook() string {
	s.TestApp()
	for {
		keyword := tool.InputStr("Please input search keyword:")
		if keyword == "next" || keyword == "n" {
			s.Page += 1 // next page
			s.TestApp()
		} else if keyword == "previous" || keyword == "p" {
			if s.Page > 0 {
				s.Page -= 1 // previous page
				s.TestApp()
			} else {
				fmt.Println("page is 0, cannot go previous")
				s.TestApp()
			}
		} else {
			if BookID := s.SearchResult[tool.InputInt(keyword, len(s.SearchResult))]; BookID != "" {
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

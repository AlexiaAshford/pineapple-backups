package src

import (
	"fmt"
	"sf/cfg"
	"sf/src/boluobao"
	"sf/src/hbooker"
	"strconv"
)

func CatSearchDetailed(searchName string, page int) []string {
	var searchResult []string
	response := hbooker.Search(searchName, page)
	if response.Code != "100000" || len(response.Data.BookList) == 0 {
		fmt.Println("search failed, code:", response.Code)
		return nil
	} else {
		fmt.Println("this page has", len(response.Data.BookList), "books")
	}
	for index, book := range response.Data.BookList {
		fmt.Println("Index:", index, "\t\t\tBookName:", book.BookName)
		searchResult = append(searchResult, book.BookID)
	}
	return searchResult
}

func SfacgSearchDetailed(keyword string, page int) []string {
	var searchResult []string
	response := boluobao.GetSearchDetailedByKeyword(keyword, page)
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

func TestApp(searchName string, page int) []string {
	if cfg.Vars.AppType == "cat" {
		return CatSearchDetailed(searchName, page)
	} else if cfg.Vars.AppType == "sfacg" {
		return SfacgSearchDetailed(searchName, page)
	} else {
		panic("app type is not correct" + cfg.Vars.AppType)
	}
}

func SearchBook(searchName string) string {
	var page int
	searchResult := TestApp(searchName, 0)
	for {
		keyword := cfg.InputStr("Please input search keyword:")
		if keyword == "next" || keyword == "n" {
			page += 1 // next page
			searchResult = TestApp(searchName, page+1)
		} else if keyword == "previous" || keyword == "p" {
			if page > 0 {
				page -= 1 // previous page
				searchResult = TestApp(searchName, page)
			} else {
				fmt.Println("page is 0, cannot go previous")
				searchResult = TestApp(searchName, 0)
			}
		} else {
			if BookID := ReturnBookID(keyword, searchResult); BookID != "" {
				return BookID // if the input is a number (book id)
			}
		}
	}
}

func ReturnBookID(keyword string, searchResult []string) string {
	if cfg.IsNum(keyword) {
		inputInt := cfg.StrToInt(keyword)
		if len(searchResult) > 0 { // if search result is not empty and input is number
			//inputs := cfg.InputInt("please input the index of the book you want to download:")
			if inputInt < len(searchResult) { // if the index is valid (less than the length of the search result)
				return searchResult[inputInt]
			} else {
				fmt.Println("index out of range, please input again")
			}
		} else { // if the search result is empty (no result)
			fmt.Println("search result is empty, please input again")
		}
	} else { // if the input is not a number (not a book id)
		fmt.Println("input is not a number, please input again")
	}
	return "" // if the input is not a number (not a book id) or the index is out of range
}

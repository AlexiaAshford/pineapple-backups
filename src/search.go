package src

import (
	"fmt"
	"github.com/VeronicaAlexia/pineapple-backups/config"
	"github.com/VeronicaAlexia/pineapple-backups/config/tool"
	"github.com/VeronicaAlexia/pineapple-backups/src/boluobao"
	"github.com/VeronicaAlexia/pineapple-backups/src/hbooker"
	"strconv"
)

func CatSearchDetailed(searchName string, page int) []string {
	var searchResult []string
	response := hbooker.GET_SEARCH(searchName, page)
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

func SfacgSearchDetailed(keyword string, page int) []string {
	var searchResult []string
	response := boluobao.GET_SEARCH(keyword, page)
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
	if config.Vars.AppType == "cat" {
		return CatSearchDetailed(searchName, page)
	} else if config.Vars.AppType == "sfacg" {
		return SfacgSearchDetailed(searchName, page)
	} else {
		panic("app type is not correct" + config.Vars.AppType)
	}
}

func SearchBook(searchName string) string {
	var page int
	searchResult := TestApp(searchName, 0)
	for {
		keyword := tool.InputStr("Please input search keyword:")
		if keyword == "next" || keyword == "n" {
			page += 1 // next page
			searchResult = TestApp(searchName, page)
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
			} else {
				fmt.Println("No found search book, please input again:")
			}
		}
	}
}

func ReturnBookID(keyword string, searchResult []string) string {
	if tool.IsNum(keyword) {
		inputInt := tool.StrToInt(keyword)
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

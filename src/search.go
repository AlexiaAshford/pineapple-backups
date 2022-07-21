package src

import (
	"fmt"
	"sf/cfg"
	"sf/src/boluobao"
	"strconv"
)

func SearchDetailed(keyword string, page int) []string {
	var searchResult []string
	response := boluobao.GetSearchDetailedByKeyword(keyword, page)
	if response.Status.HTTPCode != 200 || len(response.Data.Novels) == 0 {
		return nil
	} else {
		fmt.Println("search result length:", len(response.Data.Novels))
	}
	for index, book := range response.Data.Novels {
		fmt.Println("Index:", index, "\t\t\tBookName:", book.NovelName)
		searchResult = append(searchResult, strconv.Itoa(book.NovelID))
	}
	return searchResult
}

func SearchBook(searchName string) string {
	// if search keyword is not empty, search book and download it
	var page int
	searchResult := SearchDetailed(searchName, 0)
	for {
		keyword := cfg.InputStr("Please input search keyword:")
		fmt.Println(keyword)
		if keyword == "next" || keyword == "n" {
			page += 1 // next page
			searchResult = SearchDetailed(searchName, page+1)
		} else if keyword == "previous" || keyword == "p" {
			if page > 0 {
				page -= 1 // previous page
				searchResult = SearchDetailed(searchName, page)
			} else {
				fmt.Println("page is 0, cannot go previous")
				searchResult = SearchDetailed(searchName, 0)
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
		inputInt := StrToInt(keyword)
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

// StrToInt string to int
func StrToInt(str string) int {
	if i, err := strconv.Atoi(str); err == nil {
		return i
	} else {
		return 0
	}
}

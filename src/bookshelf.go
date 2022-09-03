package src

import (
	"fmt"
	"github.com/VeronicaAlexia/pineapple-backups/config"
	"github.com/VeronicaAlexia/pineapple-backups/src/boluobao"
	"github.com/VeronicaAlexia/pineapple-backups/src/hbooker"
	"strings"
)

func request_bookshelf_book_list() (map[int][]map[string]string, error) {
	if config.Vars.AppType == "sfacg" {
		return boluobao.GET_BOOK_SHELF_INFORMATION()
	} else if config.Vars.AppType == "cat" {
		return hbooker.GET_BOOK_SHELF_INFORMATION()
	} else {
		return nil, fmt.Errorf("app type error")
	}
}

func InitBookShelf() ([]int, []map[string]string) {
	bookshelf_book_list, response_err := request_bookshelf_book_list()
	if response_err != nil || bookshelf_book_list == nil {
		var test_login_status bool
		fmt.Println("BookShelf Error:", response_err)
		if config.Vars.AppType == "sfacg" {
			test_login_status = AutoAccount()
		} else if config.Vars.AppType == "cat" {
			test_login_status = InputAccountToken()
		}
		if !test_login_status {
			fmt.Println("please login your sfacg account and password!")
			account := config.InputStr("please input your account:")
			password := config.InputStr("please input your password:")
			LoginAccount(strings.TrimSpace(account), strings.TrimSpace(password), 0)
		}
		return InitBookShelf()

	}
	fmt.Println("\nyou account is valid, start loading bookshelf information.")
	return select_bookcase(bookshelf_book_list)

}

func select_bookcase(bookshelf_book_list map[int][]map[string]string) ([]int, []map[string]string) {
	var bookshelf_index int
	if len(bookshelf_book_list) == 1 {
		fmt.Println("you only have one bookshelf, default loading bookshelf index:1")
		bookshelf_index = 0
	} else {
		fmt.Println("please input bookshelf index:")
		bookshelf_index = config.InputInt(">", len(bookshelf_book_list))
	}
	book_shelf_bookcase := bookshelf_book_list[bookshelf_index]
	var bookshelf_book_index []int
	for book_index, book := range book_shelf_bookcase {
		fmt.Println("book-index", book_index, "book-name:", book["novel_name"], "\t\tbook-id:", book["novel_id"])
		bookshelf_book_index = append(bookshelf_book_index, book_index)
	}
	return bookshelf_book_index, book_shelf_bookcase
}

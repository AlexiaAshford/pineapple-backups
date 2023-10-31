package src

import (
	"fmt"
	"github.com/VeronicaAlexia/pineapple-backups/config"
	"github.com/VeronicaAlexia/pineapple-backups/pkg/command"
	"github.com/VeronicaAlexia/pineapple-backups/pkg/tools"
	"github.com/olekukonko/tablewriter"
	"os"
	"strconv"
)

func NewChoiceBookshelf() (map[string]string, error) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"序号", "小说名", "小说ID"})
	var bookList = make(map[string]string)
	switch command.Command.AppType {
	case "sfacg":
		var shelfIndex int
		bookShelf, err := config.APP.SFacg.Client.API.GetBookShelfInfo()
		if err != nil {
			return nil, err
		}
		if len(bookShelf) == 1 {
			fmt.Println("检测到只有一个书架，无需选择书架")
			shelfIndex = 0
		} else {
			fmt.Println("检测到多个书架，需要选择书架")
			for index, value := range bookShelf {
				fmt.Println("书架号:", index, "\t\t\t书架名:", value.Name)
			}
			shelfIndex = tools.InputInt(">", len(bookShelf))
		}
		for index, bookInfo := range bookShelf[shelfIndex].Expand.Novels {
			table.Append([]string{strconv.Itoa(index), bookInfo.NovelName, strconv.Itoa(bookInfo.NovelId)})
			bookList[strconv.Itoa(index)] = strconv.Itoa(bookInfo.NovelId)
		}
		table.Render()
	case "cat":
		var shelfIndex string
		bookShelfInfo, err := config.APP.Hbooker.Client.API.GetBookShelfInfoAPI()
		if err != nil {
			return nil, err
		}
		if len(bookShelfInfo) == 1 {
			fmt.Println("检测到只有一个书架，无需选择书架")
			shelfIndex = bookShelfInfo[0].ShelfID
		} else {
			fmt.Println("检测到多个书架，需要选择书架")
			for index, value := range bookShelfInfo {
				fmt.Println("书架号:", index, "\t\t\t书架名:", value.ShelfName)
			}
			shelfIndex = bookShelfInfo[tools.InputInt(">", len(bookShelfInfo))].ShelfID
		}

		bookInfoList, err := config.APP.Hbooker.Client.API.GetBookShelfIndexesInfoAPI(shelfIndex)
		if err != nil {
			return nil, err
		}
		for index, bookInfo := range bookInfoList {
			table.Append([]string{strconv.Itoa(index), bookInfo.BookInfo.BookName, bookInfo.BookInfo.BookID})
			bookList[strconv.Itoa(index)] = bookInfo.BookInfo.BookID
		}
	}
	return bookList, nil
}

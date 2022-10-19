package app

import (
	"fmt"
	"github.com/VeronicaAlexia/pineapple-backups/config/tool"
	"github.com/VeronicaAlexia/pineapple-backups/src/app/hbooker"
	req "github.com/VeronicaAlexia/pineapple-backups/src/https"
	"github.com/VeronicaAlexia/pineapple-backups/struct/hbooker_structs"
	"strings"
)

type RECOMMEND struct {
	book_list        []string
	recommend_list   [][]string
	book_list_string string
}

func NEW_RECOMMEND() *RECOMMEND {
	var recommend_list [][]string
	recommend := new(hbooker_structs.RecommendStruct)
	req.Get(new(req.Context).Init(hbooker.BOOKCITY_RECOMMEND_DATA).
		Query("theme_type", "NORMAL").Query("tab_type", "200").QueryToString(), recommend)
	if recommend.Code != "100000" {
		fmt.Println(recommend.Tip.(string))
	} else {
		for _, data := range recommend.Data.ModuleList {
			if data.ModuleType == "1" {
				for _, book := range data.BossModule.DesBookList {
					recommend_list = append(recommend_list, []string{book.BookName, book.BookID})
				}
			}
		}
	}
	return &RECOMMEND{recommend_list: recommend_list}

}

func (is *RECOMMEND) InitBookIdList() {
	is.book_list = nil
	for index, value := range is.recommend_list {
		fmt.Println("index:", index, "\t\tbook id:", value[1], "\t\tbook name:", value[0])
		is.book_list = append(is.book_list, value[1])
	}
	is.book_list_string = strings.Join(is.book_list, ",")
}

func (is *RECOMMEND) CHANGE_NEW_RECOMMEND() {
	s := new(hbooker_structs.ChangeRecommendStruct)
	req.Get(new(req.Context).Init(hbooker.GET_CHANGE_RECOMMEND).
		Query("book_id", is.book_list_string).Query("from_module_name", "长篇好书").QueryToString(), s)
	is.recommend_list = nil
	if s.Code != "100000" {
		fmt.Println(s.Tip)
	} else {
		for _, book := range s.Data.BookList {
			is.recommend_list = append(is.recommend_list, []string{book.BookName, book.BookID})
		}
	}
}

func (is *RECOMMEND) GET_HBOOKER_RECOMMEND() string {
	is.InitBookIdList() // init book_list_string and print recommend_list
	fmt.Println("y is next item recommendation\nd is download recommend book\npress any key to exit..")
	InputChoice := tool.InputStr(">")
	if InputChoice == "y" {
		is.CHANGE_NEW_RECOMMEND() // change recommend_list
		return is.GET_HBOOKER_RECOMMEND()
	} else if InputChoice == "d" {
		return is.book_list[tool.InputInt("input index:", len(is.book_list))]
	}
	return ""

}

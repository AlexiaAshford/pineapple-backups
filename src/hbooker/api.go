package HbookerAPI

import (
	"encoding/json"
	"fmt"
	"sf/cfg"
	req "sf/src/hbooker/request"
	structs "sf/structural/hbooker_structs"
)

func GetDivisionIdByBookId(BookId string) []structs.DivisionList {
	var result structs.DivisionStruct
	if err := json.Unmarshal(req.Get(DivisionIdByBookId+BookId, 0), &result); err != nil {
		fmt.Println("json unmarshal error:", err)
		return nil
	}
	return result.Data.DivisionList
}

func GetCatalogueByDivisionId(DivisionId string) []structs.ChapterList {
	var result structs.ChapterStruct
	if err := json.Unmarshal(req.Get(CatalogueDetailedByDivisionId+DivisionId, 0), &result); err != nil {
		fmt.Println("json unmarshal error:", err)
		return nil
	}
	return result.Data.ChapterList
}

func Login(account, password string) {
	var result structs.LoginStruct
	response := req.Get(fmt.Sprintf(LoginByAccount, account, password), 0)
	if json.Unmarshal(response, &result) == nil {
		cfg.Vars.Cat.Token = result.Data.LoginToken
		cfg.Vars.Cat.Account = result.Data.ReaderInfo.Account
		cfg.SaveJson()
	} else {
		fmt.Println("Login failed!")
	}
}
func GetBookDetailById(bid string) structs.BookInfo {
	var result structs.DetailStruct
	res := req.Get(fmt.Sprintf(BookDetailedById, bid), 0)
	if err := json.Unmarshal(res, &result); err != nil {
		fmt.Println("json unmarshal error:", err)
		return structs.BookInfo{}
	}
	return result.Data.BookInfo
}

func Search(bookName string, page int) []structs.BookList {
	res := req.Get(fmt.Sprintf(SearchDetailedByKeyword, page, bookName), 0)
	var result structs.SearchStruct
	if err := json.Unmarshal(res, &result); err != nil {
		fmt.Println("json unmarshal error:", err)
	}
	return result.Data.BookList
}

func main() []structs.ChapterList {
	Login("", "")
	GetBookDetailById("")
	Search("", 0)
	var chapterList []structs.ChapterList
	for index, division := range GetDivisionIdByBookId("") {
		fmt.Println("index:", index, "\t\tdivision:", division.DivisionName)
		for _, chapter := range GetCatalogueByDivisionId(division.DivisionID) {
			if chapter.IsValid == "1" {
				chapterList = append(chapterList, chapter)
			}
		}
	}
	return chapterList

}

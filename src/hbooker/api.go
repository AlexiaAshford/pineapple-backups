package HbookerAPI

import (
	"encoding/json"
	"fmt"
	"sf/cfg"
	req "sf/src/hbooker/request"
	"sf/src/hbooker/util"
	structs "sf/structural/hbooker_structs"
)

func GetDivisionIdByBookId(BookId string) []structs.DivisionList {
	var result structs.DivisionStruct
	response := req.Get(DivisionIdByBookId+BookId, 0)
	if err := json.Unmarshal(response, &result); err != nil {
		fmt.Println("json unmarshal error:", err)
	}
	return result.Data.DivisionList
}

func GetCatalogueByDivisionId(DivisionId string) []structs.ChapterList {
	var result structs.ChapterStruct
	response := req.Get(CatalogueDetailedByDivisionId+DivisionId, 0)
	if err := json.Unmarshal(response, &result); err != nil {
		fmt.Println("json unmarshal error:", err)
	}
	return result.Data.ChapterList
}

func Login(account, password string) {
	var result structs.LoginStruct
	response := req.Get(fmt.Sprintf(LoginByAccount, account, password), 0)
	if json.Unmarshal(response, &result) == nil {
		cfg.Vars.Cat.CommonParams.LoginToken = result.Data.LoginToken
		cfg.Vars.Cat.CommonParams.Account = result.Data.ReaderInfo.Account
		cfg.SaveJson()
	} else {
		fmt.Println("Login failed!")
	}
}
func GetBookDetailById(bid string) structs.BookInfo {
	var result structs.DetailStruct
	response := req.Get(fmt.Sprintf(BookDetailedById, bid), 0)
	if err := json.Unmarshal(response, &result); err != nil {
		fmt.Println("json unmarshal error:", err)
		return structs.BookInfo{}
	}
	return result.Data.BookInfo
}

func Search(bookName string, page int) []structs.BookList {
	var result structs.SearchStruct
	response := req.Get(fmt.Sprintf(SearchDetailedByKeyword, page, bookName), 0)
	if err := json.Unmarshal(response, &result); err != nil {
		fmt.Println("json unmarshal error:", err)
	}
	return result.Data.BookList
}

func GetKeyByCid(chapterId string) string {
	var result structs.KeyStruct
	response := req.Get(ChapterKeyByCid+chapterId, 0)
	if err := json.Unmarshal(response, &result); err != nil {
		fmt.Println("json unmarshal error:", err)
	}
	return result.Data.Command
}

func GetContent(chapterId string) structs.ChapterInfo {
	var result structs.ContentStruct
	chapterKey := GetKeyByCid(chapterId)
	response := req.Get(fmt.Sprintf(ContentDetailedByCid, chapterId, chapterKey), 0)
	if err := json.Unmarshal(response, &result); err != nil {
		fmt.Println("json unmarshal error:", err)
	}
	bytes := util.Decode(result.Data.ChapterInfo.TxtContent, chapterKey)
	result.Data.ChapterInfo.TxtContent = bytes
	return result.Data.ChapterInfo
}

func main() []structs.ChapterList {
	Login("", "")
	GetBookDetailById("")
	GetContent("")
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

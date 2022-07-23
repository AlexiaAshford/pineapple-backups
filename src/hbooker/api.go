package HbookerAPI

import (
	"encoding/json"
	"fmt"
	"os"
	req "sf/src/https"
	structs "sf/structural/hbooker_structs"
)

func GetDivisionIdByBookId(BookId string) []structs.DivisionList {
	var result structs.DivisionStruct
	response, _ := req.Request("POST", QueryParams(DivisionIdByBookId+BookId), "")
	if err := json.Unmarshal(Decode(string(response), ""), &result); err != nil {
		fmt.Println("json unmarshal error:", err)
	}
	return result.Data.DivisionList
}

func GetCatalogueByDivisionId(DivisionId string) []structs.ChapterList {
	var result structs.ChapterStruct
	response, _ := req.Request("POST", QueryParams(CatalogueDetailedByDivisionId+DivisionId), "")
	if err := json.Unmarshal(Decode(string(response), ""), &result); err != nil {
		fmt.Println("json unmarshal error:", err)
	}
	return result.Data.ChapterList
}

func GetBookDetailById(bid string) structs.DetailStruct {
	var result structs.DetailStruct
	response, _ := req.Request("POST", QueryParams(fmt.Sprintf(BookDetailedById, bid)), "")
	if err := json.Unmarshal(Decode(string(response), ""), &result); err == nil {
		return result
	} else {
		fmt.Println("BookDetailById json unmarshal error:", err)
		os.Exit(1)
	}
	return structs.DetailStruct{}
}

func Search(bookName string, page int) structs.SearchStruct {
	var result structs.SearchStruct
	response, _ := req.Request(
		"POST", QueryParams(fmt.Sprintf(SearchDetailedByKeyword, page, bookName)), "")
	if err := json.Unmarshal(Decode(string(response), ""), &result); err != nil {
		fmt.Println("json unmarshal error:", err)
	}
	return result
}

//func Login(account, password string) {
//	var result structs.LoginStruct
//	response, _ := req.Request("POST", fmt.Sprintf(WebSite+LoginByAccount, account, password), "")
//	if json.Unmarshal(Decode(string(response), ""), &result) == nil {
//		cfg.Vars.Cat.Params.LoginToken = result.Data.LoginToken
//		cfg.Vars.Cat.Params.Account = result.Data.ReaderInfo.Account
//		cfg.SaveJson()
//	} else {
//		fmt.Println("Login failed!")
//	}
//}

func GetKeyByCid(chapterId string) string {
	var result structs.KeyStruct
	response, _ := req.Request("POST", QueryParams(ChapterKeyByCid+chapterId), "")
	if err := json.Unmarshal(Decode(string(response), ""), &result); err != nil {
		fmt.Println("json unmarshal error:", err)
	}
	return result.Data.Command
}

func GetContent(chapterId string) structs.ChapterInfo {
	var result structs.ContentStruct
	chapterKey := GetKeyByCid(chapterId)
	response, _ := req.Request("POST", QueryParams(fmt.Sprintf(ContentDetailedByCid, chapterId, chapterKey)), "")
	if err := json.Unmarshal(Decode(string(response), ""), &result); err != nil {
		fmt.Println("json unmarshal error:", err)
	}
	result.Data.ChapterInfo.TxtContent = string(Decode(result.Data.ChapterInfo.TxtContent, chapterKey))
	return result.Data.ChapterInfo
}

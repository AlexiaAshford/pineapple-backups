package HbookerAPI

import (
	"encoding/json"
	"fmt"
	"sf/cfg"
	req "sf/src/https"
	structs "sf/structural/hbooker_structs"
)

func GetDivisionIdByBookId(BookId string) []structs.DivisionList {
	var result structs.DivisionStruct
	response := req.Get(DivisionIdByBookId+BookId, 0)
	if err := json.Unmarshal([]byte(Decode(string(response), "")), &result); err != nil {
		fmt.Println("json unmarshal error:", err)
	}
	return result.Data.DivisionList
}

func GetCatalogueByDivisionId(DivisionId string) []structs.ChapterList {
	var result structs.ChapterStruct
	response := req.Get(CatalogueDetailedByDivisionId+DivisionId, 0)
	if err := json.Unmarshal([]byte(Decode(string(response), "")), &result); err != nil {
		fmt.Println("json unmarshal error:", err)
	}
	return result.Data.ChapterList
}

func Login(account, password string) {
	var result structs.LoginStruct
	response := req.Get(fmt.Sprintf(LoginByAccount, account, password), 0)
	if json.Unmarshal([]byte(Decode(string(response), "")), &result) == nil {
		cfg.Vars.Cat.CommonParams.LoginToken = result.Data.LoginToken
		cfg.Vars.Cat.CommonParams.Account = result.Data.ReaderInfo.Account
		cfg.SaveJson()
	} else {
		fmt.Println("Login failed!")
	}
}
func GetBookDetailById(bid string) structs.BookInfo {
	var result structs.DetailStruct
	response := req.Get(WebSite+fmt.Sprintf(BookDetailedById, bid)+QueryParams(), 0)
	//fmt.Println("GetBookDetailById", string(response))
	if err := json.Unmarshal([]byte(Decode(string(response), "")), &result); err != nil {
		fmt.Println("json unmarshal error:", err)
		return structs.BookInfo{}
	}
	return result.Data.BookInfo
}

func Search(bookName string, page int) []structs.BookList {
	var result structs.SearchStruct
	response := req.Get(WebSite+fmt.Sprintf(SearchDetailedByKeyword, page, bookName), 0)
	if err := json.Unmarshal([]byte(Decode(string(response), "")), &result); err != nil {
		fmt.Println("json unmarshal error:", err)
	}
	return result.Data.BookList
}

// QueryParams Struct to map
func QueryParams() string {
	var Params string
	QueryMap := map[string]interface{}{
		"login_token":  cfg.Vars.Cat.CommonParams.LoginToken,
		"account":      cfg.Vars.Cat.CommonParams.Account,
		"app_version":  cfg.Vars.Cat.CommonParams.AppVersion,
		"device_token": cfg.Vars.Cat.CommonParams.DeviceToken,
	}
	for k, v := range QueryMap {
		Params += fmt.Sprintf("&%s=%s", k, v)
	}
	fmt.Println(Params)
	return Params
}

func GetKeyByCid(chapterId string) string {
	var result structs.KeyStruct
	response := req.Get(WebSite+ChapterKeyByCid+chapterId+QueryParams(), 0)
	if err := json.Unmarshal([]byte(Decode(string(response), "")), &result); err != nil {
		fmt.Println("json unmarshal error:", err)
	}
	return result.Data.Command
}

func GetContent(chapterId string) structs.ChapterInfo {
	var result structs.ContentStruct
	chapterKey := GetKeyByCid(chapterId)
	response := req.Get(WebSite+fmt.Sprintf(ContentDetailedByCid, chapterId, chapterKey)+QueryParams(), 0)
	if err := json.Unmarshal([]byte(Decode(string(response), "")), &result); err != nil {
		fmt.Println("json unmarshal error:", err)
	}
	bytes := Decode(result.Data.ChapterInfo.TxtContent, chapterKey)
	result.Data.ChapterInfo.TxtContent = bytes
	return result.Data.ChapterInfo
}

func main() {
	Login("", "")
	GetContent("")
	Search("", 0)

}

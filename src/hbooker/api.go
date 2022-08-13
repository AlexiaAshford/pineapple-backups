package hbooker

import (
	"encoding/json"
	"fmt"
	"github.com/gookit/color"
	"os"
	"sf/cfg"
	"sf/src/hbooker/Encrypt"
	req "sf/src/https"
	structs "sf/structural/hbooker_structs"
	"strconv"
	"time"
)

func GetDivisionIdByBookId(BookId string) []structs.DivisionList {
	var result structs.DivisionStruct
	params := map[string]string{"book_id": BookId}
	response, _ := req.Request("POST", QueryParams(DivisionIdByBookId, params), "")
	if err := json.Unmarshal(Encrypt.Decode(string(response), ""), &result); err != nil {
		fmt.Println("json unmarshal error:", err)
	}
	return result.Data.DivisionList
}

func GetCatalogueByDivisionId(DivisionId string) []structs.ChapterList {
	var result structs.ChapterStruct
	params := map[string]string{"division_id": DivisionId}
	response, _ := req.Request("POST", QueryParams(CatalogueDetailedByDivisionId, params), "")
	if err := json.Unmarshal(Encrypt.Decode(string(response), ""), &result); err != nil {
		fmt.Println("json unmarshal error:", err)
	}
	return result.Data.ChapterList
}

func GetBookDetailById(bid string) structs.DetailStruct {
	var result structs.DetailStruct
	params := map[string]string{"book_id": bid}
	response, _ := req.Request("POST", QueryParams(BookDetailedById, params), "")
	if err := json.Unmarshal(Encrypt.Decode(string(response), ""), &result); err == nil {
		return result
	} else {
		fmt.Println("BookDetailById json unmarshal error:", err)
		os.Exit(1)
	}
	return structs.DetailStruct{}
}

func Search(bookName string, page int) structs.SearchStruct {
	var SearchStruct structs.SearchStruct
	params := map[string]string{"count": "10", "page": strconv.Itoa(page), "category_index": "0", "key": bookName}
	response, _ := req.Request("POST", QueryParams(SearchDetailedByKeyword, params), "")
	if err := json.Unmarshal(Encrypt.Decode(string(response), ""), &SearchStruct); err != nil {
		fmt.Println("json unmarshal error:", err)
	}
	return SearchStruct
}

//func Login(account, password string) {
//	var result structs.LoginStruct
//	response, _ := req.Request("POST", fmt.Sprintf(WebSite+LoginByAccount, account, password), "")
//	if json.Unmarshal(Decode(string(response), ""), &result) == nil {
//		cfg.Apps.Cat.Params.LoginToken = result.Data.LoginToken
//		cfg.Apps.Cat.Params.Account = result.Data.ReaderInfo.Account
//		cfg.SaveJson()
//	} else {
//		fmt.Println("Login failed!")
//	}
//}

func UseGeetest() {
	var GeetestStruct structs.GeetestStruct
	response, _ := req.Request("POST", WebSite+UseGeetestSignup, "")
	if err := json.Unmarshal(Encrypt.Decode(string(response), ""), &GeetestStruct); err != nil {
		fmt.Println("json unmarshal error:", err)
	}
	fmt.Println(GeetestStruct.Tip)
}

func GeetestRegister(userID string) (string, string) {
	var GeetestChallenge structs.GeetestChallenge
	params := map[string]string{"t": strconv.FormatInt(time.Now().UnixNano()/1e6, 10), "user_id": userID}
	response, _ := req.Request("POST", QueryParams(GeetestFirstRegister, params), "")
	if err := json.Unmarshal(response, &GeetestChallenge); err != nil {
		fmt.Println("json unmarshal error:", err)
	}
	return GeetestChallenge.Challenge, GeetestChallenge.Gt
}
func TestGeetest(userID string) {
	UseGeetest()
	challenge, gt := GeetestRegister(userID)
	status, CaptchaType, errorDetail := GetFullBG(&Geetest{GT: gt, Challenge: challenge})
	fmt.Println(status, CaptchaType, errorDetail)
	if status == "success" {
		color.Infoln("验证码类型：", CaptchaType, "")
	} else {
		color.Errorln("获取图片失败 Error: ", status, " ErrorDetail:", errorDetail)
		TestGeetest(userID)
	}
}

func GetRecommend() structs.RecommendStruct {
	RecommendStruct := structs.RecommendStruct{}
	params := map[string]string{"theme_type": "NORMAL", "tab_type": "200"}
	response, _ := req.Request("POST", QueryParams(Recommend, params), "")
	if err := json.Unmarshal(Encrypt.Decode(string(response), ""), &RecommendStruct); err != nil {
		fmt.Println("json unmarshal error:", err)
	}
	return RecommendStruct
}

func JsonUnmarshal(response []byte, Struct interface{}) any {
	err := json.Unmarshal(response, Struct)
	if err != nil {
		fmt.Println("json unmarshal error:", err)
	}
	return Struct
}

func GetChangeRecommend() []structs.ChangeBookList {
	bookIdList := "100250589,100283902,100186621,100287528,100309123,100325245"
	params := map[string]string{"book_id": bookIdList, "from_module_name": "长篇好书"}
	response, _ := req.Request("POST", QueryParams(ChangeRecommend, params), "")
	result := JsonUnmarshal(Encrypt.Decode(string(response), ""), &structs.ChangeRecommendStruct{})
	return result.(*structs.ChangeRecommendStruct).Data.BookList
}

func GetKeyByCid(chapterId string) string {
	var result structs.KeyStruct
	params := map[string]string{"chapter_id": chapterId}
	response, _ := req.Request("POST", QueryParams(ChapterKeyByCid, params), "")
	if err := json.Unmarshal(Encrypt.Decode(string(response), ""), &result); err != nil {
		fmt.Println("json unmarshal error:", err)
	}
	return result.Data.Command
}

func GetContent(chapterId string) (structs.ContentStruct, bool) {
	var result structs.ContentStruct
	key := GetKeyByCid(chapterId)
	params := map[string]string{"chapter_id": chapterId, "chapter_command": key}
	response, _ := req.Request("POST", QueryParams(ContentDetailedByCid, params), "")
	if err := json.Unmarshal(Encrypt.Decode(string(response), ""), &result); err == nil {
		for i := 0; i < cfg.Vars.MaxRetry; i++ {
			if result.Code == "100000" {
				result.Data.ChapterInfo.TxtContent = string(Encrypt.Decode(result.Data.ChapterInfo.TxtContent, key))
				return result, true
			}
		}
	} else {
		fmt.Println("json unmarshal error:", err)
	}
	return structs.ContentStruct{}, false
}

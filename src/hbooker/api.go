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
	var result structs.SearchStruct
	params := map[string]string{"count": "10", "page": strconv.Itoa(page), "category_index": "0", "key": bookName}
	response, _ := req.Request("POST", QueryParams(SearchDetailedByKeyword, params), "")
	if err := json.Unmarshal(Encrypt.Decode(string(response), ""), &result); err != nil {
		fmt.Println("json unmarshal error:", err)
	}
	return result
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
	var result structs.GeetestStruct
	response, _ := req.Request("POST", WebSite+UseGeetestSignup, "")
	if err := json.Unmarshal(Encrypt.Decode(string(response), ""), &result); err != nil {
		fmt.Println("json unmarshal error:", err)
	}
	fmt.Println(result.Tip)
}

func GeetestRegister(userID string) (string, string) {
	var result structs.GeetestChallenge
	params := map[string]string{"t": strconv.FormatInt(time.Now().UnixNano()/1e6, 10), "user_id": userID}
	response, _ := req.Request("POST", QueryParams(GeetestFirstRegister, params), "")
	if err := json.Unmarshal(response, &result); err != nil {
		fmt.Println("json unmarshal error:", err)
	}
	return result.Challenge, result.Gt
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
	var result structs.RecommendStruct
	params := map[string]string{"theme_type": "NORMAL", "tab_type": "200"}
	response, _ := req.Request("POST", QueryParams(Recommend, params), "")
	fmt.Println(string(Encrypt.Decode(string(response), "")))
	if err := json.Unmarshal(Encrypt.Decode(string(response), ""), &result); err != nil {
		fmt.Println("json unmarshal error:", err)
	}
	return result
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

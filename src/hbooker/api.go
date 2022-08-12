package hbooker

import (
	"encoding/json"
	"fmt"
	"github.com/gookit/color"
	"os"
	"sf/cfg"
	"sf/src/hbooker/Geetest"
	req "sf/src/https"
	structs "sf/structural/hbooker_structs"
	"strconv"
	"time"
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
	if err := json.Unmarshal(Decode(string(response), ""), &result); err != nil {
		fmt.Println("json unmarshal error:", err)
	}
	fmt.Println(result.Tip)
}

func GeetestRegister(userID string) (string, string) {
	var result structs.GeetestChallenge
	RegisterApi := fmt.Sprintf(GeetestFirstRegister, strconv.FormatInt(time.Now().UnixNano()/1e6, 10), userID)
	response, _ := req.Request("POST", QueryParams(RegisterApi), "")
	if err := json.Unmarshal(response, &result); err != nil {
		fmt.Println("json unmarshal error:", err)
	}
	return result.Challenge, result.Gt
}
func TestGeetest(userID string) {
	UseGeetest()
	challenge, gt := GeetestRegister(userID)
	status, CaptchaType, errorDetail := Geetest.GetFullBG(&Geetest.Geetest{GT: gt, Challenge: challenge})
	fmt.Println(status, CaptchaType, errorDetail)
	if status == "success" {
		color.Infoln("验证码类型：", CaptchaType, "")
	} else {
		color.Errorln("获取图片失败 Error: ", status, " ErrorDetail:", errorDetail)
		TestGeetest(userID)
	}
}

func GetKeyByCid(chapterId string) string {
	var result structs.KeyStruct
	response, _ := req.Request("POST", QueryParams(ChapterKeyByCid+chapterId), "")
	if err := json.Unmarshal(Decode(string(response), ""), &result); err != nil {
		fmt.Println("json unmarshal error:", err)
	}
	return result.Data.Command
}

func GetContent(chapterId string) (structs.ContentStruct, bool) {
	var result structs.ContentStruct
	chapterKey := GetKeyByCid(chapterId)
	response, _ := req.Request("POST", QueryParams(fmt.Sprintf(ContentDetailedByCid, chapterId, chapterKey)), "")
	if err := json.Unmarshal(Decode(string(response), ""), &result); err == nil {
		for i := 0; i < cfg.Vars.MaxRetry; i++ {
			if result.Code == "100000" {
				result.Data.ChapterInfo.TxtContent = string(Decode(result.Data.ChapterInfo.TxtContent, chapterKey))
				return result, true
			}
		}
	} else {
		fmt.Println("json unmarshal error:", err)
	}
	return structs.ContentStruct{}, false
}

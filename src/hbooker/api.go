package hbooker

import (
	"fmt"
	"github.com/gookit/color"
	"sf/cfg"
	"sf/src/hbooker/Encrypt"
	req "sf/src/https"
	structs "sf/structural/hbooker_structs"
	"strconv"
	"time"
)

func GetDivisionIdByBookId(BookId string) []structs.DivisionList {
	url := QueryParams(DivisionIdByBookId, map[string]string{"book_id": BookId})
	return req.JsonUnmarshal(req.Get("POST", url), &structs.DivisionStruct{}).(*structs.DivisionStruct).Data.DivisionList
}

func GetCatalogueByDivisionId(DivisionId string) []structs.ChapterList {
	url := QueryParams(CatalogueDetailedByDivisionId, map[string]string{"division_id": DivisionId})
	return req.JsonUnmarshal(req.Get("POST", url), &structs.ChapterStruct{}).(*structs.ChapterStruct).Data.ChapterList
}

func GetBookDetailById(bid string) *structs.DetailStruct {
	url := QueryParams(BookDetailedById, map[string]string{"book_id": bid})
	return req.JsonUnmarshal(req.Get("POST", url), &structs.DetailStruct{}).(*structs.DetailStruct)

}

func Search(KeyWord string, page int) *structs.SearchStruct {
	params := map[string]string{"count": "10", "page": strconv.Itoa(page), "category_index": "0", "key": KeyWord}
	return req.JsonUnmarshal(req.Get("POST", QueryParams(SearchDetailed, params)), &structs.SearchStruct{}).(*structs.SearchStruct)
}

func Login(account, password string) {
	url := QueryParams(LoginByAccount, map[string]string{"login_name": account, "password": password})
	result := req.JsonUnmarshal(req.Get("POST", url), &structs.LoginStruct{}).(*structs.LoginStruct)
	if result.Code == "100000" {
		cfg.Apps.Cat.Params.LoginToken = result.Data.LoginToken
		cfg.Apps.Cat.Params.Account = result.Data.ReaderInfo.Account
		cfg.SaveJson()
	} else {
		fmt.Println("Login failed!")
	}
}

func UseGeetest() *structs.GeetestStruct {
	return req.JsonUnmarshal(req.Get("POST", WebSite+UseGeetestSignup), &structs.GeetestStruct{}).(*structs.GeetestStruct)
}

func GeetestRegister(userID string) (string, string) {
	params := map[string]string{"t": strconv.FormatInt(time.Now().UnixNano()/1e6, 10), "user_id": userID}
	response, _ := req.Request("POST", QueryParams(GeetestFirstRegister, params))
	result := req.JsonUnmarshal(response, &structs.GeetestChallenge{}).(*structs.GeetestChallenge)
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

func GetRecommend() *structs.RecommendStruct {
	url := QueryParams(Recommend, map[string]string{"theme_type": "NORMAL", "tab_type": "200"})
	return req.JsonUnmarshal(req.Get("POST", url), &structs.RecommendStruct{}).(*structs.RecommendStruct)
}
func GetChangeRecommend() []structs.ChangeBookList {
	bookIdList := "100250589,100283902,100186621,100287528,100309123,100325245"
	url := QueryParams(ChangeRecommend, map[string]string{"book_id": bookIdList, "from_module_name": "长篇好书"})
	return req.JsonUnmarshal(req.Get("POST", url), &structs.ChangeRecommendStruct{}).(*structs.ChangeRecommendStruct).Data.BookList
}

func GetKeyByCid(chapterId string) string {
	url := QueryParams(ChapterKeyByCid, map[string]string{"chapter_id": chapterId})
	return req.JsonUnmarshal(req.Get("POST", url), &structs.KeyStruct{}).(*structs.KeyStruct).Data.Command
}

func GetContent(chapterId string) (*structs.ContentStruct, bool) {
	key := GetKeyByCid(chapterId)
	url := QueryParams(ContentDetailedByCid, map[string]string{"chapter_id": chapterId, "chapter_command": key})
	result := req.JsonUnmarshal(req.Get("POST", url), &structs.ContentStruct{}).(*structs.ContentStruct)
	for retry := 0; retry < cfg.Vars.MaxRetry; retry++ {
		if result.Code == "100000" {
			TxtContent := Encrypt.Decode(result.Data.ChapterInfo.TxtContent, key)
			result.Data.ChapterInfo.TxtContent = string(TxtContent)
			return result, true
		}
	}
	return &structs.ContentStruct{}, false
}

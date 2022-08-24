package hbooker

import (
	"fmt"
	"github.com/gookit/color"
	"sf/cfg"
	"sf/src/hbooker/Encrypt"
	req "sf/src/https"
	structs "sf/struct/hbooker_structs"
	"strconv"
	"time"
)

func GetDivisionIdByBookId(BookId string) []structs.DivisionList {
	url := req.QueryParams(req.CatDivisionIdByBookId, map[string]string{"book_id": BookId})
	return req.Get(url, &structs.DivisionStruct{}).(*structs.DivisionStruct).Data.DivisionList
}

func GetCatalogueByDivisionId(DivisionId string) []structs.ChapterList {
	url := req.QueryParams(req.CatalogueDetailedByDivisionId, map[string]string{"division_id": DivisionId})
	return req.Get(url, &structs.ChapterStruct{}).(*structs.ChapterStruct).Data.ChapterList
}

func GetBookDetailById(bid string) *structs.DetailStruct {
	url := req.QueryParams(req.CatBookDetailedById, map[string]string{"book_id": bid})
	return req.Get(url, &structs.DetailStruct{}).(*structs.DetailStruct)

}

func Search(KeyWord string, page int) *structs.SearchStruct {
	params := map[string]string{"count": "10", "page": strconv.Itoa(page), "category_index": "0", "key": KeyWord}
	return req.Get(req.QueryParams(req.CatSearchDetailed, params), &structs.SearchStruct{}).(*structs.SearchStruct)
}

func Login(account, password string) {
	url := req.QueryParams(req.CatLoginByAccount, map[string]string{"login_name": account, "password": password})
	result := req.Get(url, &structs.LoginStruct{}).(*structs.LoginStruct)
	if result.Code == "100000" {
		cfg.Apps.Cat.Params.LoginToken = result.Data.LoginToken
		cfg.Apps.Cat.Params.Account = result.Data.ReaderInfo.Account
		cfg.SaveJson()
	} else {
		fmt.Println("Login failed!")
	}
}

func UseGeetest() *structs.GeetestStruct {
	return req.Get(req.CatWebSite+req.CatUseGeetestSignup, &structs.GeetestStruct{}).(*structs.GeetestStruct)
}

func GeetestRegister(userID string) (string, string) {
	params := map[string]string{"t": strconv.FormatInt(time.Now().UnixNano()/1e6, 10), "user_id": userID}
	response, _ := req.Request("POST", req.QueryParams(req.CatGeetestFirstRegister, params))
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
	url := req.QueryParams(req.CatRecommend, map[string]string{"theme_type": "NORMAL", "tab_type": "200"})
	return req.Get(url, &structs.RecommendStruct{}).(*structs.RecommendStruct)
}
func GetChangeRecommend() []structs.ChangeBookList {
	bookIdList := "100250589,100283902,100186621,100287528,100309123,100325245"
	url := req.QueryParams(req.CatChangeRecommend, map[string]string{"book_id": bookIdList, "from_module_name": "长篇好书"})
	return req.Get(url, &structs.ChangeRecommendStruct{}).(*structs.ChangeRecommendStruct).Data.BookList
}

func GetKeyByCid(chapterId string) string {
	url := req.QueryParams(req.CatChapterKeyByCid, map[string]string{"chapter_id": chapterId})
	return req.Get(url, &structs.KeyStruct{}).(*structs.KeyStruct).Data.Command
}

func GetContent(chapterId string) (*structs.ContentStruct, bool) {
	key := GetKeyByCid(chapterId)
	url := req.QueryParams(req.CatContentDetailedByCid, map[string]string{"chapter_id": chapterId, "chapter_command": key})
	result := req.Get(url, &structs.ContentStruct{}).(*structs.ContentStruct)
	for retry := 0; retry < cfg.Vars.MaxRetry; retry++ {
		if result.Code == "100000" {
			TxtContent := Encrypt.Decode(result.Data.ChapterInfo.TxtContent, key)
			result.Data.ChapterInfo.TxtContent = string(TxtContent)
			return result, true
		}
	}
	return &structs.ContentStruct{}, false
}

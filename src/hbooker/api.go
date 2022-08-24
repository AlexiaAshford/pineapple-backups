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
	response := req.Get(req.CatDivisionIdByBookId, &structs.DivisionStruct{}, map[string]string{"book_id": BookId})
	return response.(*structs.DivisionStruct).Data.DivisionList
}

func GetCatalogueByDivisionId(DivisionId string) []structs.ChapterList {
	response := req.Get(req.CatalogueAPI(""), &structs.ChapterStruct{}, map[string]string{"division_id": DivisionId})
	return response.(*structs.ChapterStruct).Data.ChapterList
}

func GetBookDetailById(bid string) *structs.DetailStruct {
	return req.Get(req.BookInfoAPI(""), &structs.DetailStruct{}, map[string]string{"book_id": bid}).(*structs.DetailStruct)

}

func Search(KeyWord string, page int) *structs.SearchStruct {
	params := map[string]string{"count": "10", "page": strconv.Itoa(page), "category_index": "0", "key": KeyWord}
	return req.Get(req.SearchAPI(), &structs.SearchStruct{}, params).(*structs.SearchStruct)
}

func Login(account, password string) {
	params := map[string]string{"login_name": account, "password": password}
	response := req.Get(req.CatLoginByAccount, &structs.LoginStruct{}, params)
	result := response.(*structs.LoginStruct)
	if result.Code == "100000" {
		cfg.Apps.Cat.Params.LoginToken = result.Data.LoginToken
		cfg.Apps.Cat.Params.Account = result.Data.ReaderInfo.Account
		cfg.SaveJson()
	} else {
		fmt.Println("Login failed!")
	}
}

func UseGeetest() *structs.GeetestStruct {
	return req.Get(req.CatUseGeetestSignup, &structs.GeetestStruct{}, nil).(*structs.GeetestStruct)
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
	response := req.Get(req.CatRecommend, &structs.RecommendStruct{}, map[string]string{"theme_type": "NORMAL", "tab_type": "200"})
	return response.(*structs.RecommendStruct)
}
func GetChangeRecommend() []structs.ChangeBookList {
	bookIdList := "100250589,100283902,100186621,100287528,100309123,100325245"
	response := req.Get(req.CatChangeRecommend, &structs.ChangeRecommendStruct{}, map[string]string{"book_id": bookIdList, "from_module_name": "长篇好书"})
	return response.(*structs.ChangeRecommendStruct).Data.BookList
}

func GetKeyByCid(chapterId string) string {
	return req.Get(req.CatChapterKeyByCid, &structs.KeyStruct{}, map[string]string{"chapter_id": chapterId}).(*structs.KeyStruct).Data.Command
}

func GetContent(chapterId string, ChapterKey string) (*structs.ContentStruct, bool) {
	params := map[string]string{"chapter_id": chapterId, "chapter_command": ChapterKey}
	result := req.Get(req.CatContentDetailedByCid, &structs.ContentStruct{}, params).(*structs.ContentStruct)
	for retry := 0; retry < cfg.Vars.MaxRetry; retry++ {
		if result.Code == "100000" {
			TxtContent := Encrypt.Decode(result.Data.ChapterInfo.TxtContent, ChapterKey)
			result.Data.ChapterInfo.TxtContent = string(TxtContent)
			return result, true
		}
	}
	return &structs.ContentStruct{}, false
}

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

func GET_DIVISION(BookId string) []structs.DivisionList {
	response := req.Get("book/get_division_list", &structs.DivisionStruct{}, map[string]string{"book_id": BookId})
	return response.(*structs.DivisionStruct).Data.DivisionList
}

func GET_CATALOGUE(DivisionId string) []structs.ChapterList {
	params := map[string]string{"division_id": DivisionId}
	return req.Get("chapter/get_updated_chapter_by_division_id", &structs.ChapterStruct{}, params).(*structs.ChapterStruct).Data.ChapterList
}

func GET_BOOK_INFORMATION(bid string) *structs.DetailStruct {
	return req.Get("book/get_info_by_id", &structs.DetailStruct{}, map[string]string{"book_id": bid}).(*structs.DetailStruct)

}

func GET_SEARCH(KeyWord string, page int) *structs.SearchStruct {
	params := map[string]string{"count": "10", "page": strconv.Itoa(page), "category_index": "0", "key": KeyWord}
	return req.Get("bookcity/get_filter_search_book_list", &structs.SearchStruct{}, params).(*structs.SearchStruct)
}

func Login(account, password string) {
	params := map[string]string{"login_name": account, "password": password}
	response := req.Get("signup/login", &structs.LoginStruct{}, params)
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
	return req.Get("signup/use_geetest", &structs.GeetestStruct{}, nil).(*structs.GeetestStruct)
}

func GeetestRegister(userID string) (string, string) {
	params := map[string]string{"t": strconv.FormatInt(time.Now().UnixNano()/1e6, 10), "user_id": userID}
	response, _ := req.Request("POST", req.QueryParams("signup/geetest_first_register", params))
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
	params := map[string]string{"theme_type": "NORMAL", "tab_type": "200"}
	return req.Get("bookcity/get_index_list", &structs.RecommendStruct{}, params).(*structs.RecommendStruct)
}
func GetChangeRecommend() []structs.ChangeBookList {
	params := map[string]string{"book_id": "100250589,100283902,100186621,100287528,100309123,100325245", "from_module_name": "长篇好书"}
	return req.Get("bookcity/change_recommend_exposure_books", &structs.ChangeRecommendStruct{}, params).(*structs.ChangeRecommendStruct).Data.BookList
}

func GetKeyByCid(chapterId string) string {
	params := map[string]string{"chapter_id": chapterId}
	return req.Get("chapter/get_chapter_cmd", &structs.KeyStruct{}, params).(*structs.KeyStruct).Data.Command
}

func GetContent(chapterId string, ChapterKey string) (*structs.ContentStruct, bool) {
	params := map[string]string{"chapter_id": chapterId, "chapter_command": ChapterKey}
	result := req.Get("chapter/get_cpt_ifm", &structs.ContentStruct{}, params).(*structs.ContentStruct)
	for retry := 0; retry < cfg.Vars.MaxRetry; retry++ {
		if result.Code == "100000" {
			TxtContent := Encrypt.Decode(result.Data.ChapterInfo.TxtContent, ChapterKey)
			result.Data.ChapterInfo.TxtContent = string(TxtContent)
			return result, true
		}
	}
	return &structs.ContentStruct{}, false
}

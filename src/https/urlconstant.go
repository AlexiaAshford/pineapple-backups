package https

import (
	"fmt"
	"sf/cfg"
)

const (
	CatalogueDetailedByDivisionId = "chapter/get_updated_chapter_by_division_id"
	CatDivisionIdByBookId         = "book/get_division_list"
	CatBookDetailedById           = "book/get_info_by_id"
	CatSearchDetailed             = "bookcity/get_filter_search_book_list"
	CatChapterKeyByCid            = "chapter/get_chapter_cmd"
	CatContentDetailedByCid       = "chapter/get_cpt_ifm"
	CatRecommend                  = "bookcity/get_index_list"
	CatChangeRecommend            = "bookcity/change_recommend_exposure_books"
	CatWebSite                    = "https://app.hbooker.com/"
	CatUseGeetestSignup           = "signup/use_geetest"
	CatGeetestFirstRegister       = "signup/geetest_first_register"
	CatLoginByAccount             = "signup/login"
)

// QueryParams Struct to map
func QueryParams(url string, ParamsData map[string]string) string {
	var Params string
	queryRequisite := map[string]interface{}{
		"login_token":  cfg.Apps.Cat.Params.LoginToken,
		"account":      cfg.Apps.Cat.Params.Account,
		"app_version":  cfg.Apps.Cat.Params.AppVersion,
		"device_token": cfg.Apps.Cat.Params.DeviceToken,
	}
	for k, v := range ParamsData {
		Params += fmt.Sprintf("&%s=%s", k, v)
	}
	for k, v := range queryRequisite {
		Params += fmt.Sprintf("&%s=%s", k, v)
	}
	return CatWebSite + url + "?" + Params
}

const (
	SFBookDetailedById        = "novels/%v?expand="
	SFWebSite                 = "https://minipapi.sfacg.com/pas/mpapi/"
	SFAccountDetailedByApi    = "user"
	SFCatalogueDetailedById   = "novels/%v/dirs?expand=originNeedFireMoney"
	SFContentDetailedByCid    = "Chaps/%v?expand=content&autoOrder=true"
	SFSearchDetailedByKeyword = "search/novels/result?q=%v&size=20&page=%v&expand="
	SFLoginByAccount          = "sessions"
	//SFBookShelfDetailed       = "novels/%v/shelf?expand="
)

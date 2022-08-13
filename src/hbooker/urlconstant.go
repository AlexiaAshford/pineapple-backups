package hbooker

import (
	"fmt"
	"sf/cfg"
)

const (
	CatalogueDetailedByDivisionId = "chapter/get_updated_chapter_by_division_id"
	DivisionIdByBookId            = "book/get_division_list"
	BookDetailedById              = "book/get_info_by_id"
	SearchDetailedByKeyword       = "bookcity/get_filter_search_book_list"
	ChapterKeyByCid               = "chapter/get_chapter_cmd"
	ContentDetailedByCid          = "chapter/get_cpt_ifm"
	Recommend                     = "bookcity/get_index_list"
	ChangeRecommend               = "bookcity/change_recommend_exposure_books"
	WebSite                       = "https://app.hbooker.com/"
	UseGeetestSignup              = "signup/use_geetest"
	GeetestFirstRegister          = "signup/geetest_first_register"
	//LoginByAccount                = "signup/login?login_name=%v&password=%v"
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
	return WebSite + url + "?" + Params
}

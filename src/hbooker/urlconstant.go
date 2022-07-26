package hbooker

import (
	"fmt"
	"sf/cfg"
)

const (
	CatalogueDetailedByDivisionId = "chapter/get_updated_chapter_by_division_id?division_id="
	DivisionIdByBookId            = "book/get_division_list?book_id="
	BookDetailedById              = "book/get_info_by_id?book_id=%v"
	SearchDetailedByKeyword       = "bookcity/get_filter_search_book_list?count=10&page=%v&category_index=0&key=%v"
	ChapterKeyByCid               = "chapter/get_chapter_cmd?chapter_id="
	ContentDetailedByCid          = "chapter/get_cpt_ifm?chapter_id=%v&chapter_command=%v"
	WebSite                       = "https://app.hbooker.com/"
	//LoginByAccount                = "signup/login?login_name=%v&password=%v"
)

// QueryParams Struct to map
func QueryParams(url string) string {
	var Params string
	QueryMap := map[string]interface{}{
		"login_token":  cfg.Vars.Cat.Params.LoginToken,
		"account":      cfg.Vars.Cat.Params.Account,
		"app_version":  cfg.Vars.Cat.Params.AppVersion,
		"device_token": cfg.Vars.Cat.Params.DeviceToken,
	}
	for k, v := range QueryMap {
		Params += fmt.Sprintf("&%s=%s", k, v)
	}
	return WebSite + url + Params
}

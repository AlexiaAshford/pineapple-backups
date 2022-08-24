package https

import (
	"fmt"
	url_ "net/url"
	"sf/cfg"
)

const (
	CatDivisionIdByBookId   = "book/get_division_list"
	CatChapterKeyByCid      = "chapter/get_chapter_cmd"
	CatContentDetailedByCid = "chapter/get_cpt_ifm"
	CatRecommend            = "bookcity/get_index_list"
	CatChangeRecommend      = "bookcity/change_recommend_exposure_books"
	CatWebSite              = "https://app.hbooker.com/"
	CatUseGeetestSignup     = "signup/use_geetest"
	CatGeetestFirstRegister = "signup/geetest_first_register"
	CatLoginByAccount       = "signup/login"
)

const (
	SFWebSite              = "https://minipapi.sfacg.com/pas/mpapi/"
	SFAccountDetailedByApi = "user"
	SFContentDetailedByCid = "Chaps/%v?expand=content&autoOrder=true"
	SFLogin                = "sessions"
	//SFBookShelfDetailed       = "novels/%v/shelf?expand="
)

func SearchAPI(keyWord string, page int) string {
	//params := map[string]string{"q": url_.QueryEscape(keyWord), "size": "20", "page": strconv.Itoa(page)}
	switch cfg.Vars.AppType {
	case "cat":
		return "bookcity/get_filter_search_book_list"
	case "sfacg":
		return fmt.Sprintf("search/novels/result?q=%v&size=20&page=%v&expand=", url_.QueryEscape(keyWord), page)
	}
	return ""
}

func BookInfoAPI(NovelId string) string {
	switch cfg.Vars.AppType {
	case "cat":
		return "book/get_info_by_id"
	case "sfacg":
		return fmt.Sprintf("novels/%v?expand=", NovelId)
	}
	return ""
}

func CatalogueAPI(NovelId string) string {
	switch cfg.Vars.AppType {
	case "cat":
		return "chapter/get_updated_chapter_by_division_id"
	case "sfacg":
		return fmt.Sprintf("novels/%v/dirs?expand=originNeedFireMoney", NovelId)
	}
	return ""
}

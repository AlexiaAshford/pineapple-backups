package boluobao

import (
	"fmt"
	url_ "net/url"
	"sf/cfg"
	req "sf/src/https"
	"sf/structural/sfacg_structs"
)

func GetBookDetailedById(bookId string) *sfacg_structs.BookInfo {
	url := WebSite + fmt.Sprintf(BookDetailedById, bookId)
	return req.JsonUnmarshal(req.Get("GET", url), &sfacg_structs.BookInfo{}).(*sfacg_structs.BookInfo)
}

func GetAccountDetailedByApi() *sfacg_structs.Account {
	url := WebSite + AccountDetailedByApi
	return req.JsonUnmarshal(req.Get("GET", url), &sfacg_structs.Account{}).(*sfacg_structs.Account)
}

func GetCatalogue(NovelID string) *sfacg_structs.Catalogue {
	url := fmt.Sprintf(WebSite+CatalogueDetailedById, NovelID)
	return req.JsonUnmarshal(req.Get("GET", url), &sfacg_structs.Catalogue{}).(*sfacg_structs.Catalogue)
}

func GetContentDetailedByCid(cid string) (*sfacg_structs.Content, bool) {
	url := fmt.Sprintf(WebSite+ContentDetailedByCid, cid)
	result := req.JsonUnmarshal(req.Get("GET", url), &sfacg_structs.Content{}).(*sfacg_structs.Content)
	for retry := 0; retry < cfg.Vars.MaxRetry; retry++ {
		if result.Status.HTTPCode == 200 {
			return result, true
		} else {
			fmt.Println("Error:", result.Status.Msg)
		}
	}
	return result, false
}

func GetSearchDetailedByKeyword(keyword string, page int) *sfacg_structs.Search {
	url := WebSite + fmt.Sprintf(SearchDetailedByKeyword, url_.QueryEscape(keyword), page)
	return req.JsonUnmarshal(req.Get("GET", url), &sfacg_structs.Search{}).(*sfacg_structs.Search)

}

func PostLoginByAccount(username, password string) *sfacg_structs.Login {
	data := fmt.Sprintf(`{"username":"%s", "password": "%s"}`, username, password)
	response, Cookie := req.LoginSession(WebSite+LoginByAccount, []byte(data))
	result := req.JsonUnmarshal(response, &sfacg_structs.Login{}).(*sfacg_structs.Login)
	for _, cookie := range Cookie {
		result.Cookie += cookie.Name + "=" + cookie.Value + ";"
	}
	return result
}

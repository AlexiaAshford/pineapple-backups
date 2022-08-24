package boluobao

import (
	"fmt"
	url_ "net/url"
	"sf/cfg"
	req "sf/src/https"
	"sf/struct/sfacg_structs"
)

func GetBookDetailedById(bookId string) *sfacg_structs.BookInfo {
	return req.Get(req.SET_URL(fmt.Sprintf(req.SFBookDetailedById, bookId)), &sfacg_structs.BookInfo{}).(*sfacg_structs.BookInfo)
}

func GetAccountDetailedByApi() *sfacg_structs.Account {
	url := req.SET_URL(req.SFAccountDetailedByApi)
	return req.Get(url, &sfacg_structs.Account{}).(*sfacg_structs.Account)
}

func GetCatalogue(NovelID string) *sfacg_structs.Catalogue {
	url := req.SET_URL(fmt.Sprintf(req.SFCatalogueDetailedById, NovelID))
	return req.Get(url, &sfacg_structs.Catalogue{}).(*sfacg_structs.Catalogue)
}

func GetContentDetailedByCid(cid string) (*sfacg_structs.Content, bool) {
	url := req.SET_URL(fmt.Sprintf(req.SFContentDetailedByCid, cid))
	result := req.Get(url, &sfacg_structs.Content{}).(*sfacg_structs.Content)
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
	url := req.SET_URL(fmt.Sprintf(req.SFSearchDetailedByKeyword, url_.QueryEscape(keyword), page))
	return req.Get(url, &sfacg_structs.Search{}).(*sfacg_structs.Search)

}

func PostLoginByAccount(username, password string) *sfacg_structs.Login {
	data := fmt.Sprintf(`{"username":"%s", "password": "%s"}`, username, password)
	response, Cookie := req.LoginSession(req.SET_URL(req.SFLoginByAccount), []byte(data))
	result := req.JsonUnmarshal(response, &sfacg_structs.Login{}).(*sfacg_structs.Login)
	for _, cookie := range Cookie {
		result.Cookie += cookie.Name + "=" + cookie.Value + ";"
	}
	return result
}

package boluobao

import (
	"fmt"
	url_ "net/url"
	"sf/cfg"
	req "sf/src/https"
	"sf/struct/sfacg_structs"
)

func GetBookDetailedById(bookId string) *sfacg_structs.BookInfo {
	return req.Get(fmt.Sprintf(req.SFBookDetailedById, bookId), &sfacg_structs.BookInfo{}, nil).(*sfacg_structs.BookInfo)
}

func GetAccountDetailedByApi() *sfacg_structs.Account {
	return req.Get(req.SFAccountDetailedByApi, &sfacg_structs.Account{}, nil).(*sfacg_structs.Account)
}

func GetCatalogue(NovelID string) *sfacg_structs.Catalogue {
	return req.Get(fmt.Sprintf(req.SFCatalogueDetailedById, NovelID), &sfacg_structs.Catalogue{}, nil).(*sfacg_structs.Catalogue)
}

func GetContentDetailedByCid(cid string) (*sfacg_structs.Content, bool) {
	result := req.Get(fmt.Sprintf(req.SFContentDetailedByCid, cid), &sfacg_structs.Content{}, nil).(*sfacg_structs.Content)
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
	return req.Get(fmt.Sprintf(req.SFSearch, url_.QueryEscape(keyword), page), &sfacg_structs.Search{}, nil).(*sfacg_structs.Search)

}

func PostLoginByAccount(username, password string) *sfacg_structs.Login {
	data := []byte(fmt.Sprintf(`{"username":"%s", "password": "%s"}`, username, password))
	response, Cookie := req.LoginSession(req.SET_URL(req.SFLogin, nil), data)
	result := req.JsonUnmarshal(response, &sfacg_structs.Login{}).(*sfacg_structs.Login)
	for _, cookie := range Cookie {
		result.Cookie += cookie.Name + "=" + cookie.Value + ";"
	}
	return result
}

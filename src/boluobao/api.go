package boluobao

import (
	"fmt"
	url_ "net/url"
	req "sf/src/https"
	"sf/struct/sfacg_structs"
	"strconv"
)

func GET_BOOK_INFORMATION(NovelId string) *sfacg_structs.BookInfo {
	params := map[string]string{"expand": "intro,tags,sysTags,totalNeedFireMoney,originTotalNeedFireMoney"}
	return req.Get("novels/"+NovelId, &sfacg_structs.BookInfo{}, params).(*sfacg_structs.BookInfo)
}

func GET_ACCOUNT_INFORMATION() *sfacg_structs.Account {
	return req.Get("user", &sfacg_structs.Account{}, nil).(*sfacg_structs.Account)
}

func GET_CATALOGUE(NovelID string) *sfacg_structs.Catalogue {
	params := map[string]string{"expand": "originNeedFireMoney"}
	return req.Get(fmt.Sprintf("novels/%v/dirs", NovelID), &sfacg_structs.Catalogue{}, params).(*sfacg_structs.Catalogue)
}

func GET_CONTENT(cid string) *sfacg_structs.Content {
	params := map[string]string{"expand": "content&autoOrder=true"}
	return req.Get("Chaps/"+cid, &sfacg_structs.Content{}, params).(*sfacg_structs.Content)
}

func GET_SEARCH(keyword string, page int) *sfacg_structs.Search {
	params := map[string]string{"q": url_.QueryEscape(keyword), "size": "20", "page": strconv.Itoa(page)}
	return req.Get("search/novels/result", &sfacg_structs.Search{}, params).(*sfacg_structs.Search)

}

func PostLoginByAccount(username, password string) *sfacg_structs.Login {
	data := fmt.Sprintf(`{"username":"%s", "password": "%s"}`, username, password)
	response, Cookie := req.LoginSession(req.SET_URL("sessions", nil), []byte(data))
	result := req.JsonUnmarshal(response, &sfacg_structs.Login{}).(*sfacg_structs.Login)
	for _, cookie := range Cookie {
		result.Cookie += cookie.Name + "=" + cookie.Value + ";"
	}
	return result
}

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
	params := map[string]string{"expand": "content"}
	if result := req.Get("Chaps/"+cid, &sfacg_structs.Content{}, params); result != nil {
		return result.(*sfacg_structs.Content)
	} else {
		return GET_CONTENT(cid) // retry once if failed to get content
	}
}

func GET_SEARCH(keyword string, page int) *sfacg_structs.Search {
	params := map[string]string{"q": url_.QueryEscape(keyword), "size": "20", "page": strconv.Itoa(page)}
	return req.Get("search/novels/result", &sfacg_structs.Search{}, params).(*sfacg_structs.Search)

}

func LOGIN_ACCOUNT(username, password string) *sfacg_structs.Login {
	params := fmt.Sprintf(`{"username":"%s", "password": "%s"}`, username, password)
	response, Cookie := req.Login(req.SET_URL("sessions", nil), []byte(params))
	for _, cookie := range Cookie {
		response.Cookie += cookie.Name + "=" + cookie.Value + ";"
	}
	return response
}

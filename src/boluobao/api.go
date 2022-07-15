package boluobao

import (
	"encoding/json"
	"fmt"
	"net/url"
	req "sf/src/boluobao/request"
	"sf/structural/sfacg_structs"
)

func GetBookDetailedById(bookId string) sfacg_structs.BookInfo {
	var BookData sfacg_structs.BookInfo
	response := req.Get(fmt.Sprintf(BookDetailedById, bookId), 0)
	if err := json.Unmarshal(response, &BookData); err == nil {
		return BookData
	} else {
		fmt.Println("Error:", err)
		return sfacg_structs.BookInfo{}
	}
}
func GetAccountDetailedByApi() sfacg_structs.Account {
	var AccountData sfacg_structs.Account
	if err := json.Unmarshal(req.Get(AccountDetailedByApi, 0), &AccountData); err == nil {
		return AccountData
	} else {
		fmt.Println("Error:", err)
		return sfacg_structs.Account{}
	}
}

func GetCatalogueDetailedById(NovelID string) sfacg_structs.Catalogue {
	var CatalogueData sfacg_structs.Catalogue
	response := req.Get(fmt.Sprintf(CatalogueDetailedById, NovelID), 0)
	if err := json.Unmarshal(response, &CatalogueData); err == nil {
		return CatalogueData
	} else {
		fmt.Println("Error:", err)
		return sfacg_structs.Catalogue{}
	}
}

func GetContentDetailedByCid(cid string) sfacg_structs.Content {
	var ContentData sfacg_structs.Content
	response := req.Get(fmt.Sprintf(ContentDetailedByCid, cid), 0)
	if err := json.Unmarshal(response, &ContentData); err == nil {
		return ContentData
	} else {
		fmt.Println("Error:", err)
		return sfacg_structs.Content{}
	}
}

func GetSearchDetailedByKeyword(keyword string) sfacg_structs.Search {
	var SearchData sfacg_structs.Search
	response := req.Get(fmt.Sprintf(SearchDetailedByKeyword, url.QueryEscape(keyword)), 0)
	if err := json.Unmarshal(response, &SearchData); err == nil {
		return SearchData // return result of search
	} else {
		fmt.Println("Error:", err)
		return sfacg_structs.Search{} // return empty struct if error
	}

}

func PostLoginByAccount(username, password string) sfacg_structs.Login {
	var LoginData sfacg_structs.Login
	result, Cookie := req.POST(LoginByAccount,
		fmt.Sprintf(`{"username":"%s", "password": "%s"}`, username, password),
	)
	if err := json.Unmarshal(result, &LoginData); err != nil {
		fmt.Println(err)
		return sfacg_structs.Login{}
	}
	LoginData.Cookie = Cookie
	return LoginData
}

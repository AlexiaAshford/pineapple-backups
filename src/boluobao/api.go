package boluobao

import (
	"encoding/json"
	"fmt"
	"net/url"
	req "sf/src/boluobao/request"
	structs2 "sf/src/sfacg_structs"
)

func GetBookDetailedById(bookId string) structs2.BookInfo {
	var BookData structs2.BookInfo
	response := req.Get(fmt.Sprintf(BookDetailedById, bookId))
	if err := json.Unmarshal(response, &BookData); err == nil {
		return BookData
	} else {
		fmt.Println("Error:", err)
		return structs2.BookInfo{}
	}
}
func GetAccountDetailedByApi() structs2.Account {
	var AccountData structs2.Account
	if err := json.Unmarshal(req.Get(AccountDetailedByApi), &AccountData); err == nil {
		return AccountData
	} else {
		fmt.Println("Error:", err)
		return structs2.Account{}
	}
}

func GetCatalogueDetailedById(NovelID string) structs2.Catalogue {
	var CatalogueData structs2.Catalogue
	response := req.Get(fmt.Sprintf(CatalogueDetailedById, NovelID))
	if err := json.Unmarshal(response, &CatalogueData); err == nil {
		return CatalogueData
	} else {
		fmt.Println("Error:", err)
		return structs2.Catalogue{}
	}
}

func GetContentDetailedByCid(cid string) structs2.Content {
	var ContentData structs2.Content
	response := req.Get(fmt.Sprintf(ContentDetailedByCid, cid))
	if err := json.Unmarshal(response, &ContentData); err == nil {
		return ContentData
	} else {
		fmt.Println("Error:", err)
		return structs2.Content{}
	}
}

func GetSearchDetailedByKeyword(keyword string) structs2.Search {
	var SearchData structs2.Search
	response := req.Get(fmt.Sprintf(SearchDetailedByKeyword, url.QueryEscape(keyword)))
	if err := json.Unmarshal(response, &SearchData); err == nil {
		return SearchData // return result of search
	} else {
		fmt.Println("Error:", err)
		return structs2.Search{} // return empty struct if error
	}

}

func PostLoginByAccount(username, password string) structs2.Login {
	var LoginData structs2.Login
	result, Cookie := req.POST(LoginByAccount,
		fmt.Sprintf(`{"username":"%s", "password": "%s"}`, username, password),
	)
	if err := json.Unmarshal(result, &LoginData); err != nil {
		fmt.Println(err)
		return structs2.Login{}
	}
	LoginData.Cookie = Cookie
	return LoginData
}

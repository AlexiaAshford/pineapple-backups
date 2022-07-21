package boluobao

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sf/src/https"
	"sf/structural/sfacg_structs"
)

func GetBookDetailedById(bookId string) sfacg_structs.BookInfo {
	var BookData sfacg_structs.BookInfo
	response := https.Get(WebSite+fmt.Sprintf(BookDetailedById, bookId), 0)
	if err := json.Unmarshal(response, &BookData); err == nil {
		return BookData
	} else {
		fmt.Println("Error:", err)
		return sfacg_structs.BookInfo{}
	}
}
func GetAccountDetailedByApi() sfacg_structs.Account {
	var AccountData sfacg_structs.Account
	response := https.Get(WebSite+AccountDetailedByApi, 0)
	if err := json.Unmarshal(response, &AccountData); err == nil {
		return AccountData
	} else {
		fmt.Println("Error:", err)
		return sfacg_structs.Account{}
	}
}

func GetCatalogueDetailedById(NovelID string) sfacg_structs.Catalogue {
	var CatalogueData sfacg_structs.Catalogue
	response := https.Get(fmt.Sprintf(WebSite+CatalogueDetailedById, NovelID), 0)
	if err := json.Unmarshal(response, &CatalogueData); err == nil {
		return CatalogueData
	} else {
		fmt.Println("Error:", err)
		return sfacg_structs.Catalogue{}
	}
}

func GetContentDetailedByCid(cid string) sfacg_structs.Content {
	var ContentData sfacg_structs.Content
	response := https.Get(fmt.Sprintf(WebSite+ContentDetailedByCid, cid), 0)
	if err := json.Unmarshal(response, &ContentData); err == nil {
		return ContentData
	} else {
		fmt.Println("Error:", err)
		return sfacg_structs.Content{}
	}
}

func GetSearchDetailedByKeyword(keyword string, page int) sfacg_structs.Search {
	var SearchData sfacg_structs.Search
	searchAPI := fmt.Sprintf(WebSite+SearchDetailedByKeyword, url.QueryEscape(keyword), page)
	if err := json.Unmarshal(https.Get(searchAPI, 0), &SearchData); err == nil {
		return SearchData // return result of search
	} else {
		fmt.Println("Error:", err)
		return sfacg_structs.Search{} // return empty struct if error
	}

}

func PostLoginByAccount(username, password string) sfacg_structs.Login {
	var LoginData sfacg_structs.Login
	result, Cookie := https.POST(WebSite+LoginByAccount,
		fmt.Sprintf(`{"username":"%s", "password": "%s"}`, username, password),
	)
	if err := json.Unmarshal(result, &LoginData); err != nil {
		fmt.Println(err)
		return sfacg_structs.Login{}
	}
	LoginData.Cookie = Cookie
	return LoginData
}

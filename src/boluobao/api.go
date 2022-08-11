package boluobao

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sf/cfg"
	"sf/src/https"
	"sf/structural/sfacg_structs"
)

func GetBookDetailedById(bookId string) sfacg_structs.BookInfo {
	var BookData sfacg_structs.BookInfo
	response, _ := https.Request("GET", WebSite+fmt.Sprintf(BookDetailedById, bookId), "")
	if err := json.Unmarshal(response, &BookData); err == nil {
		return BookData
	} else {
		fmt.Println("Error:", err)
		return sfacg_structs.BookInfo{}
	}
}
func GetAccountDetailedByApi() sfacg_structs.Account {
	var AccountData sfacg_structs.Account
	response, _ := https.Request("GET", WebSite+AccountDetailedByApi, "")
	if err := json.Unmarshal(response, &AccountData); err == nil {
		return AccountData
	} else {
		fmt.Println("Error:", err)
		return sfacg_structs.Account{}
	}
}

func GetCatalogue(NovelID string) sfacg_structs.Catalogue {
	var CatalogueData sfacg_structs.Catalogue
	response, _ := https.Request("GET", fmt.Sprintf(WebSite+CatalogueDetailedById, NovelID), "")
	if err := json.Unmarshal(response, &CatalogueData); err == nil {
		return CatalogueData
	} else {
		fmt.Println("Error:", err)
		return sfacg_structs.Catalogue{}
	}
}

func GetContentDetailedByCid(cid string) (sfacg_structs.Content, bool) {
	var ContentData sfacg_structs.Content
	response, _ := https.Request("GET", fmt.Sprintf(WebSite+ContentDetailedByCid, cid), "")
	if err := json.Unmarshal(response, &ContentData); err == nil {
		for i := 0; i < cfg.Vars.MaxRetry; i++ {
			if ContentData.Status.HTTPCode == 200 {
				return ContentData, true
			} else {
				fmt.Println("Error:", ContentData.Status.Msg)
			}
		}
	} else {
		fmt.Println("ContentData Json Error:", err)
	}
	return sfacg_structs.Content{}, false
}

func GetSearchDetailedByKeyword(keyword string, page int) sfacg_structs.Search {
	var SearchData sfacg_structs.Search
	response, _ := https.Request("GET",
		WebSite+fmt.Sprintf(SearchDetailedByKeyword, url.QueryEscape(keyword), page), "")
	if err := json.Unmarshal(response, &SearchData); err == nil {
		return SearchData // return result of search
	} else {
		fmt.Println("Error:", err)
		return sfacg_structs.Search{} // return empty struct if error
	}

}

func PostLoginByAccount(username, password string) sfacg_structs.Login {
	var LoginData sfacg_structs.Login
	result, Cookie := https.Request("POST", WebSite+LoginByAccount,
		fmt.Sprintf(`{"username":"%s", "password": "%s"}`, username, password),
	)
	if err := json.Unmarshal(result, &LoginData); err != nil {
		fmt.Println("LoginData err:", err)
		return sfacg_structs.Login{}
	}
	LoginData.Cookie = Cookie
	return LoginData
}

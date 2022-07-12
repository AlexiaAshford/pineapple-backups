package boluobao

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	req "sf/src/boluobao/request"
	structs2 "sf/src/sfacg_structs"
)

func Get_book_detailed_by_id(bookId string) structs2.BookInfo {
	var BookData structs2.BookInfo
	response := req.Get(fmt.Sprintf("novels/%v?expand=", bookId))
	if err := json.Unmarshal(response, &BookData); err == nil {
		return BookData
	} else {
		fmt.Println("Error:", err)
		return structs2.BookInfo{}
	}
}
func Get_account_detailed_by_api() structs2.Account {
	var AccountData structs2.Account
	if err := json.Unmarshal(req.Get("user"), &AccountData); err == nil {
		return AccountData
	} else {
		fmt.Println("Error:", err)
		return structs2.Account{}
	}
}

func Get_catalogue_detailed_by_id(NovelID string) structs2.Catalogue {
	var CatalogueData structs2.Catalogue
	response := req.Get("novels/" + NovelID + "/dirs?expand=originNeedFireMoney")
	if err := json.Unmarshal(response, &CatalogueData); err == nil {
		return CatalogueData
	} else {
		fmt.Println("Error:", err)
		return structs2.Catalogue{}
	}
}

func Get_content_detailed_by_cid(cid string) structs2.Content {
	var ContentData structs2.Content
	response := req.Get(fmt.Sprintf("Chaps/%v?expand=content&autoOrder=true", cid))
	if err := json.Unmarshal(response, &ContentData); err == nil {
		return ContentData
	} else {
		fmt.Println("Error:", err)
		return structs2.Content{}
	}
}

func Get_search_detailed_by_keyword(keyword string) structs2.Search {
	var SearchData structs2.Search
	keywordParams := url.QueryEscape(keyword) + "&size=20&page=0&expand="
	response := req.Get(fmt.Sprintf("search/novels/result?q=%v", keywordParams))
	if err := json.Unmarshal(response, &SearchData); err == nil {
		return SearchData // return result of search
	} else {
		fmt.Println("Error:", err)
		return structs2.Search{} // return empty struct if error
	}

}

func Post_login_by_account(username, password string) ([]*http.Cookie, structs2.Login) {
	var LoginData structs2.Login
	result, Cookie := req.POST("sessions",
		fmt.Sprintf(`{"username":"%s", "password": "%s"}`, username, password),
	)
	if err := json.Unmarshal(result, &LoginData); err != nil {
		panic(err)
	} else {
		return Cookie, LoginData
	}
}

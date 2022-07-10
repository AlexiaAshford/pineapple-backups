package src

import (
	"encoding/json"
	"fmt"
	"net/http"
	req "sf/src/sfacg/request"
	"sf/structs"
)

func Get_book_detailed_by_id(bookId string) structs.BookInfo {
	var BookData structs.BookInfo
	response := req.Get(fmt.Sprintf("novels/%v?expand=", bookId))
	if err := json.Unmarshal(response, &BookData); err == nil {
		return BookData
	} else {
		fmt.Println("Error:", err)
		return structs.BookInfo{}
	}
}
func Get_account_detailed_by_api() structs.Account {
	var AccountData structs.Account
	if err := json.Unmarshal(req.Get("user"), &AccountData); err == nil {
		return AccountData
	} else {
		fmt.Println("Error:", err)
		return structs.Account{}
	}
}

func Get_catalogue_detailed_by_id(NovelID string) structs.Catalogue {
	var CatalogueData structs.Catalogue
	response := req.Get("novels/" + NovelID + "/dirs?expand=originNeedFireMoney")
	if err := json.Unmarshal(response, &CatalogueData); err == nil {
		return CatalogueData
	} else {
		fmt.Println("Error:", err)
		return structs.Catalogue{}
	}
}

func Get_content_detailed_by_cid(cid string) structs.Content {
	var ContentData structs.Content
	response := req.Get(fmt.Sprintf("Chaps/%v?expand=content&autoOrder=true", cid))
	if err := json.Unmarshal(response, &ContentData); err == nil {
		return ContentData
	} else {
		fmt.Println("Error:", err)
		return structs.Content{}
	}
}

func Post_login_by_account(username, password string) ([]*http.Cookie, structs.Login) {
	var LoginData structs.Login
	result, CookieArray := req.POST("sessions",
		fmt.Sprintf(`{"username":"%s", "password": "%s"}`, username, password),
	)
	if err := json.Unmarshal(result, &LoginData); err != nil {
		panic(err)
	} else {
		return CookieArray, LoginData
	}
}

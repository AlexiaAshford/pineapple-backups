package src

import (
	"fmt"
	"sf/src/boluobao"
	cfg "sf/src/config"
)

func AccountDetailed() string {
	response := boluobao.Get_account_detailed_by_api()
	if response.Status.HTTPCode == 200 {
		return fmt.Sprintf("AccountName:%v", response.Data.NickName)
	} else {
		fmt.Println(response.Status.Msg)
		return response.Status.Msg
	}

}

func LoginAccount(username string, password string) {
	CookieArray, status := boluobao.Post_login_by_account(username, password)
	if status.Status.HTTPCode == 200 {
		cfg.Var.Cookie = ""
		for _, cookie := range CookieArray {
			cfg.Var.Cookie += cookie.Name + "=" + cookie.Value + ";"
		}
		cfg.Var.UserName, cfg.Var.Password = username, password
		cfg.SaveJson()
		fmt.Println(AccountDetailed(), "\tLogin successful!")
	} else {
		fmt.Println("Login failed:", status.Status.Msg)
	}
}

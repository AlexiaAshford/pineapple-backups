package src

import (
	"fmt"
	"sf/src/boluobao"
	cfg "sf/src/config"
)

func AccountDetailed() string {
	response := boluobao.GetAccountDetailedByApi()
	if response.Status.HTTPCode == 200 {
		return fmt.Sprintf("AccountName:%v", response.Data.NickName)
	} else {
		return response.Status.Msg
	}

}

func LoginAccount(username string, password string) {
	LoginData := boluobao.PostLoginByAccount(username, password)
	if LoginData.Status.HTTPCode == 200 {
		cfg.Load()
		cfg.Var.Cookie = ""
		for _, cookie := range LoginData.Cookie {
			cfg.Var.Cookie += cookie.Name + "=" + cookie.Value + ";"
		}
		cfg.Var.UserName, cfg.Var.Password = username, password
		cfg.SaveJson()
		if AccountDetailed() == "需要登录才能访问该资源" {
			fmt.Println("Login failed, login again")
			LoginAccount(username, password)
		} else {
			fmt.Println("Login successful!\t", AccountDetailed())
		}
	} else {
		fmt.Println("Login failed:", LoginData.Status.Msg)
	}
}

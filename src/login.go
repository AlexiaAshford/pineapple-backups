package src

import (
	"fmt"
	"os"
	cfg "sf/configuration"
	"sf/src/boluobao"
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
	cfg.Var.Sfacg.Cookie = ""
	if LoginData.Status.HTTPCode == 200 {
		cfg.Load()
		cfg.Var.Sfacg.Cookie = ""
		for _, cookie := range LoginData.Cookie {
			cfg.Var.Sfacg.Cookie += cookie.Name + "=" + cookie.Value + ";"
		}
		cfg.Var.Sfacg.UserName, cfg.Var.Sfacg.Password = username, password
		cfg.SaveJson()
		if AccountDetailed() == "需要登录才能访问该资源" {
			fmt.Println("Login failed, login again")
			LoginAccount(username, password)
		} else {
			fmt.Println("Login successful!\t", AccountDetailed())
		}
	} else {
		fmt.Println("Login failed:", LoginData.Status.Msg)
		os.Exit(1)
	}
}

package src

import (
	"fmt"
	"os"
	cfg "sf/cfg"
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
	cfg.Vars.Sfacg.Cookie = ""
	if LoginData.Status.HTTPCode == 200 {
		cfg.Load()
		cfg.Vars.Sfacg.Cookie = ""
		for _, cookie := range LoginData.Cookie {
			cfg.Vars.Sfacg.Cookie += cookie.Name + "=" + cookie.Value + ";"
		}
		cfg.Vars.Sfacg.UserName, cfg.Vars.Sfacg.Password = username, password
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

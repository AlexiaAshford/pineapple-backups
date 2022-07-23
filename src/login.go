package src

import (
	"fmt"
	"os"
	"sf/cfg"
	"sf/src/boluobao"
)

func AccountDetailed() string {
	response := boluobao.GetAccountDetailedByApi()
	if response.Status.HTTPCode == 200 {
		return fmt.Sprintf("account name:%v\taccount id:%v",
			response.Data.NickName, response.Data.AccountID,
		)
	} else {
		return response.Status.Msg
	}

}

func LoginAccount(username string, password string, retry int) {
	LoginData := boluobao.PostLoginByAccount(username, password)
	cfg.Vars.Sfacg.Cookie = ""
	if LoginData.Status.HTTPCode == 200 {
		cfg.Vars.Sfacg.Cookie = ""
		for _, cookie := range LoginData.Cookie {
			cfg.Vars.Sfacg.Cookie += cookie.Name + "=" + cookie.Value + ";"
		}
		cfg.Vars.Sfacg.UserName, cfg.Vars.Sfacg.Password = username, password
		cfg.SaveJson()
		if AccountDetailed() == "需要登录才能访问该资源" {
			fmt.Println("Your login attempt was not successful, try again retry:", retry+1)
			LoginAccount(username, password, retry+1)
		} else {
			if retry == 0 {
				fmt.Println("Login successful!\n" + AccountDetailed())
			} else {
				fmt.Println("Login again successful!\n" + AccountDetailed())
			}
		}
	} else {
		fmt.Println("Login failed:", LoginData.Status.Msg)
		os.Exit(1)
	}
}

package src

import (
	"fmt"
	"github.com/VeronicaAlexia/pineapple-backups/config"
	"github.com/VeronicaAlexia/pineapple-backups/pkg/tools"
)

func LoginAccount(username string, password string) {
	cookie, err := config.APP.SFacg.Client.API.Login(username, password)
	if err != nil {
		fmt.Println("login failed!" + err.Error())
		return
	}
	config.APP.SFacg.Client.API.HttpClient.Cookie = cookie
	config.Apps.Sfacg.Cookie = cookie
	config.SaveJson()
}

func InputAccountToken() bool {
	for i := 0; i < config.Vars.MaxRetry; i++ {
		LoginToken := tools.InputStr("you must input 32 characters login token:")
		if len(LoginToken) != 32 {
			fmt.Println("Login token is 32 characters, please input again:")
		} else {
			config.Apps.Hbooker.LoginToken = LoginToken
			config.Apps.Hbooker.Account = tools.InputStr("you must input account:")
			config.SaveJson()
			return true
		}
	}
	return false
}

package src

import (
	"fmt"
	"github.com/VeronicaAlexia/BoluobaoAPI/boluobao/account"
	"github.com/VeronicaAlexia/pineapple-backups/config"
	"github.com/VeronicaAlexia/pineapple-backups/pkg/tools"
	"os"
)

func AccountDetailed() string {
	response := account.GET_ACCOUNT_INFORMATION()
	if response.Status.HTTPCode == 200 {
		return fmt.Sprintf("account name:%v\taccount id:%v",
			response.Data.NickName, response.Data.AccountID,
		)
	} else {
		return response.Status.Msg.(string)
	}

}

func LoginAccount(username string, password string, retry int) {
	if retry > 3 {
		fmt.Println("login max retry, login failed!")
		os.Exit(1)
	}
	config.Apps.Sfacg.Cookie = account.LOGIN_ACCOUNT(username, password)
	config.Apps.Sfacg.UserName = username
	config.Apps.Sfacg.Password = password
	config.SaveJson()
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
}

func TestCatAccount() bool {
	if config.Apps.Cat.Params.Account != "" && config.Apps.Cat.Params.LoginToken != "" {
		return true
	} else {
		if ok := InputAccountToken(); !ok {
			return false
		}
		return true
	}
}

func AutoAccount() bool {
	if config.Apps.Sfacg.UserName != "" && config.Apps.Sfacg.Password != "" {
		if AccountDetailed() == "需要登录才能访问该资源" {
			fmt.Printf("cookie is Invalid,attempt to auto login!\naccount:%v\npassword:%v\n",
				config.Apps.Sfacg.UserName, config.Apps.Sfacg.Password)
			// auto login and get cookie
			LoginAccount(config.Apps.Sfacg.UserName, config.Apps.Sfacg.Password, 0)
		}
		return true
	}
	return false
}

//func LoginAccountToken() bool {
//	hbooker.GET_LOGIN_TOKEN()
//}

func InputAccountToken() bool {
	for i := 0; i < config.Vars.MaxRetry; i++ {
		LoginToken := tools.InputStr("you must input 32 characters login token:")
		if len(LoginToken) != 32 {
			fmt.Println("Login token is 32 characters, please input again:")
		} else {
			config.Apps.Cat.Params.LoginToken = LoginToken
			config.Apps.Cat.Params.Account = tools.InputStr("you must input account:")
			config.SaveJson()
			return true
		}
	}
	return false
}

func TestAppTypeAndAccount() {
	// test AppType and Account is valid
	switch config.Vars.AppType {
	case "cat":
		if !TestCatAccount() {
			fmt.Println("please input account and login token, please input again:")
			os.Exit(1)
		}
	case "sfacg":
		if !AutoAccount() {
			fmt.Println("please input account and password, please input again:")
			os.Exit(1)
		}
	default:
		panic("app type %v is invalid, please input again:" + config.Vars.AppType)
	}

}

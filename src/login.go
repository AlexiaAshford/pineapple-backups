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
	cfg.Apps.Sfacg.Cookie = ""
	if LoginData.Status.HTTPCode == 200 {
		cfg.Apps.Sfacg.Cookie = ""
		for _, cookie := range LoginData.Cookie {
			cfg.Apps.Sfacg.Cookie += cookie.Name + "=" + cookie.Value + ";"
		}
		cfg.Apps.Sfacg.UserName, cfg.Apps.Sfacg.Password = username, password
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

func TestCatAccount() bool {
	if cfg.Apps.Cat.Params.Account != "" && cfg.Apps.Cat.Params.LoginToken != "" {
		return true
	} else {
		if ok := InputAccountToken(); !ok {
			return false
		}
		return true
	}
}

func AutoAccount() bool {
	if cfg.Apps.Sfacg.UserName != "" && cfg.Apps.Sfacg.Password != "" {
		if AccountDetailed() == "需要登录才能访问该资源" {
			fmt.Printf("cookie is Invalid,attempt to auto login!\naccount:%v\npassword:%v\n",
				cfg.Apps.Sfacg.UserName, cfg.Apps.Sfacg.Password)
			// auto login and get cookie
			LoginAccount(cfg.Apps.Sfacg.UserName, cfg.Apps.Sfacg.Password, 0)
		}
		return true
	}
	return false
}

func InputAccountToken() bool {
	for i := 0; i < cfg.Vars.MaxRetry; i++ {
		LoginToken := cfg.InputStr("you must input 32 characters login token:")
		if len(LoginToken) != 32 {
			fmt.Println("Login token is 32 characters, please input again:")
		} else {
			cfg.Apps.Cat.Params.LoginToken = LoginToken
			cfg.Apps.Cat.Params.Account = cfg.InputStr("you must input account:")
			cfg.SaveJson()
			return true
		}
	}
	return false
}

func TestAppTypeAndAccount() {
	if !cfg.TestList([]string{"sfacg", "cat"}, cfg.Vars.AppType) { // check app type and cheng edit config
		panic("app type %v is invalid, please input again:" + cfg.Vars.AppType)
	}

	if cfg.Vars.AppType == "cat" {
		if !TestCatAccount() {
			fmt.Println("please input account and login token, please input again:")
			os.Exit(1)
		}
		cfg.Vars.AppType = "cat"
	} else if cfg.Vars.AppType == "sfacg" {
		if !AutoAccount() {
			fmt.Println("please input account and password, please input again:")
			os.Exit(1)
		}
	}
}

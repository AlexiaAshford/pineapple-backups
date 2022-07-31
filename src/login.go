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

func TestCatAccount() bool {
	if cfg.Vars.Cat.Params.Account != "" && cfg.Vars.Cat.Params.LoginToken != "" {
		return true
	} else {
		if ok := InputAccountToken(); !ok {
			return false
		}
		return true
	}
}

func AutoAccount() bool {
	if cfg.Vars.Sfacg.UserName != "" && cfg.Vars.Sfacg.Password != "" {
		if cfg.Vars.Sfacg.Cookie != "" {
			if AccountDetailed() == "需要登录才能访问该资源" {
				fmt.Printf("cookie is Invalid,attempt to auto login!\naccount:%v\npassword:%v\n",
					cfg.Vars.Sfacg.UserName, cfg.Vars.Sfacg.Password)
				// auto login and get cookie
				LoginAccount(cfg.Vars.Sfacg.UserName, cfg.Vars.Sfacg.Password, 0)
			}
		} else {
			fmt.Printf("cookie is empty,attempt to auto login!\naccount:%v\npassword:%v\n",
				cfg.Vars.Sfacg.UserName, cfg.Vars.Sfacg.Password)
			// auto login and get cookie
			LoginAccount(cfg.Vars.Sfacg.UserName, cfg.Vars.Sfacg.Password, 0)
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
			cfg.Vars.Cat.Params.LoginToken = LoginToken
			cfg.Vars.Cat.Params.Account = cfg.InputStr("you must input account:")
			cfg.SaveJson()
			return true
		}
	}
	return false
}

func TestAppTypeAndAccount(appType string) {
	if cfg.TestList([]string{"sfacg", "cat"}, appType) { // check app type and cheng edit config
		cfg.Vars.AppType = appType
		if cfg.Vars.AppType == "cat" {
			if !TestCatAccount() {
				fmt.Println("input account and login token, please input again:")
				os.Exit(1)
			}
			cfg.Vars.AppType = "cat"
		} else if cfg.Vars.AppType == "sfacg" {
			if !AutoAccount() {
				fmt.Println("input account and password, please input again:")
				os.Exit(1)
			}
		}
	} else {
		panic("app type %v is invalid, please input again:" + appType)
	}
}

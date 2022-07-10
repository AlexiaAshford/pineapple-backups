package sfacg

import (
	"encoding/json"
	"fmt"
	cfg "sf/config"
	"sf/src/request"
)

func AccountDetailed() string {
	var AccountInfor UserStruct
	if err := json.Unmarshal(request.Get("user"), &AccountInfor); err == nil {
		if AccountInfor.Status.HTTPCode == 200 {
			return fmt.Sprintf("AccountName:%v", AccountInfor.Data.NickName)
		} else {
			fmt.Println(AccountInfor.Status.Msg)
			return AccountInfor.Status.Msg
		}
	}
	return "AccountDetailed error"
}

func LoginAccount(username string, password string) {
	var status LoginStatus
	result, CookieArray := request.POST("sessions",
		fmt.Sprintf(`{"username":"%s", "password": "%s"}`, username, password),
	)
	if err := json.Unmarshal(result, &status); err != nil {
		panic(err)
	}
	if status.Status.HTTPCode == 200 {
		cfg.Load()
		CookieMap := make(map[string]string)
		for _, cookie := range CookieArray {
			CookieMap[cookie.Name] = cookie.Value
		}
		cfg.Var.Cookie, cfg.Var.UserName, cfg.Var.Password = CookieMap, username, password
		cfg.SaveJson()
		fmt.Println(AccountDetailed(), "\tLogin successful!")
	} else {
		fmt.Println("Login failed:", status.Status.Msg)
	}
}

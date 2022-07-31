package https

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"sf/cfg"
)

var client = &http.Client{}

func Base64Bytes() string {
	var encoded bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &encoded)
	authentication := []byte(cfg.Vars.Sfacg.UserName + "&" + cfg.Vars.Sfacg.Password)
	if _, err := encoder.Write(authentication); err == nil {
		if err = encoder.Close(); err == nil {
			return string(encoded.Bytes())
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Println(err)
	}
	return ""

}

func CatAppHeaders() map[string]string {
	HeaderCollection := make(map[string]string)
	HeaderCollection["User-Agent"] = cfg.Vars.Cat.UserAgent
	HeaderCollection["Content-Type"] = "application/json"
	HeaderCollection["Cookie"] = "Account:" + cfg.Vars.Cat.Params.Account +
		"; LoginToken:" + cfg.Vars.Cat.Params.LoginToken
	return HeaderCollection
}

func SfWeChatHeaders(TestCookie bool) map[string]string {
	HeaderCollection := make(map[string]string)
	if cfg.Vars.Sfacg.Cookie == "" && TestCookie == true {
		fmt.Println("Cookie is empty, please login first!")
		os.Exit(1)
	}
	HeaderCollection["sf-minip-info"] = cfg.Vars.Sfacg.UserAgent
	HeaderCollection["Content-Type"] = "application/json"
	HeaderCollection["test-sfacg-cookie"] = "cookie:" + cfg.Vars.Sfacg.Cookie
	HeaderCollection["Authorization"] = Base64Bytes()
	HeaderCollection["Cookie"] = cfg.Vars.Sfacg.Cookie
	HeaderCollection["account-sfacg"] = cfg.Vars.Sfacg.UserName + "&" + cfg.Vars.Sfacg.Password
	return HeaderCollection
}

func SetHeaders(req *http.Request, TestCookie bool) {
	var HeaderCollection map[string]string
	if cfg.Vars.AppType == "sfacg" {
		HeaderCollection = SfWeChatHeaders(TestCookie)
	} else if cfg.Vars.AppType == "cat" {
		HeaderCollection = CatAppHeaders()
	} else {
		fmt.Println(cfg.Vars.AppType, "AppType is invalid, please check config file")
		os.Exit(1)
	}

	for k, v := range HeaderCollection {
		req.Header.Set(k, v)
	}
}

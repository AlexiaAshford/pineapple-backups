package https

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"sf/cfg"
)

func Base64Bytes() string {
	var encoded bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &encoded)
	authentication := []byte(cfg.Apps.Sfacg.UserName + "&" + cfg.Apps.Sfacg.Password)
	if _, err := encoder.Write(authentication); err == nil {
		defer func(encoder io.WriteCloser) {
			_ = encoder.Close()
		}(encoder)
		return string(encoded.Bytes())
	} else {
		fmt.Println("encoder.Write:", err)
	}
	return ""

}

func CatAppHeaders() map[string]string {
	HeaderCollection := make(map[string]string)
	HeaderCollection["User-Agent"] = cfg.Apps.Cat.UserAgent
	HeaderCollection["Content-Type"] = "application/json"
	HeaderCollection["Cookie"] = "Account:" + cfg.Apps.Cat.Params.Account +
		"; LoginToken:" + cfg.Apps.Cat.Params.LoginToken
	return HeaderCollection
}

func SfWeChatHeaders(TestCookie bool) map[string]string {
	HeaderCollection := make(map[string]string)
	if cfg.Apps.Sfacg.Cookie == "" && TestCookie == true {
		fmt.Println("Cookie is empty, please login first!")
		os.Exit(1)
	}
	HeaderCollection["sf-minip-info"] = cfg.Apps.Sfacg.UserAgent
	HeaderCollection["Content-Type"] = "application/json"
	HeaderCollection["test-sfacg-cookie"] = "cookie:" + cfg.Apps.Sfacg.Cookie
	HeaderCollection["Authorization"] = Base64Bytes()
	HeaderCollection["Cookie"] = cfg.Apps.Sfacg.Cookie
	HeaderCollection["account-sfacg"] = cfg.Apps.Sfacg.UserName + "&" + cfg.Apps.Sfacg.Password
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

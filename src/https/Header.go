package https

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"sf/cfg"
)

func Base64Bytes(UserName, Password string) string {
	var encoded bytes.Buffer
	authentication := []byte(UserName + "&" + Password)
	if _, err := base64.NewEncoder(base64.StdEncoding, &encoded).Write(authentication); err == nil {
		return string(encoded.Bytes())
	} else {
		fmt.Println("encoder.Write:", err)
	}
	return ""

}

func SET_THE_HEADERS(req *http.Request, TestCookie bool) {
	HeaderCollection := make(map[string]string)
	HeaderCollection["Content-Type"] = "application/json"
	switch cfg.Vars.AppType {

	case "sfacg":
		if cfg.Apps.Sfacg.Cookie == "" && TestCookie == true {
			fmt.Println(req.URL.String())
			fmt.Println("Cookie is empty, please login first!")
			os.Exit(1)
		}
		HeaderCollection["sf-minip-info"] = cfg.Apps.Sfacg.UserAgent
		//HeaderCollection["Authorization"] = Base64Bytes(cfg.Apps.Sfacg.UserName, cfg.Apps.Sfacg.Password)
		HeaderCollection["Cookie"] = cfg.Apps.Sfacg.Cookie
		//HeaderCollection["account-sfacg"] = cfg.Apps.Sfacg.UserName + "&" + cfg.Apps.Sfacg.Password
	case "cat":
		HeaderCollection["User-Agent"] = cfg.Apps.Cat.UserAgent
		HeaderCollection["Cookie"] = "Account:" + cfg.Apps.Cat.Params.Account + ";" + cfg.Apps.Cat.Params.LoginToken
		HeaderCollection["Authorization"] = Base64Bytes(cfg.Apps.Cat.Params.Account, cfg.Apps.Cat.Params.LoginToken)

	default:
		fmt.Println(cfg.Vars.AppType, "AppType is invalid, please check config file")
		os.Exit(1)
	}
	for HeaderKey, HeaderValue := range HeaderCollection {
		req.Header.Set(HeaderKey, HeaderValue)

	}
}

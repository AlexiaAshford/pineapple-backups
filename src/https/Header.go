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

func SetHeaders(req *http.Request, TestCookie bool) {
	cfg.Load()
	var encoded bytes.Buffer
	HeaderCollection := make(map[string]string)
	if cfg.Vars.Sfacg.Cookie == "" && TestCookie == true {
		fmt.Println("Cookie is empty, please login first!")
		os.Exit(1)
	}
	HeaderCollection["Cookie"] = cfg.Vars.Sfacg.Cookie
	HeaderCollection["sf-minip-info"] = "minip_novel/1.0.70(android;11)/wxmp"
	HeaderCollection["Content-Type"] = "application/json"
	HeaderCollection["test-sfacg-cookie"] = "cookie:" + cfg.Vars.Sfacg.Cookie
	encoder := base64.NewEncoder(base64.StdEncoding, &encoded)
	authentication := []byte(cfg.Vars.Sfacg.UserName + "&" + cfg.Vars.Sfacg.Password)
	if _, err := encoder.Write(authentication); err == nil {
		if err = encoder.Close(); err == nil {
			HeaderCollection["Authorization"] = string(encoded.Bytes())
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Println(err)
	}
	HeaderCollection["account-sfacg"] = string(authentication)
	for k, v := range HeaderCollection {
		req.Header.Set(k, v)
	}
}

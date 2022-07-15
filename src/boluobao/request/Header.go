package request

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"sf/configuration"
)

var client = &http.Client{}

func SetHeaders(req *http.Request, TestCookie bool) {
	configuration.Load()
	var encoded bytes.Buffer
	HeaderCollection := make(map[string]string)
	if configuration.Var.Sfacg.Cookie == "" && TestCookie == true {
		fmt.Println("Cookie is empty, please login first!")
		os.Exit(1)
	}
	HeaderCollection["Cookie"] = configuration.Var.Sfacg.Cookie
	HeaderCollection["sf-minip-info"] = "minip_novel/1.0.70(android;11)/wxmp"
	HeaderCollection["Content-Type"] = "application/json"
	HeaderCollection["test-sfacg-cookie"] = "cookie:" + configuration.Var.Sfacg.Cookie
	encoder := base64.NewEncoder(base64.StdEncoding, &encoded)
	authentication := []byte(configuration.Var.Sfacg.UserName + "&" + configuration.Var.Sfacg.Password)
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

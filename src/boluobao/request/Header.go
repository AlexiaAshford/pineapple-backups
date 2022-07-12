package request

import (
	"fmt"
	"net/http"
	"os"
	"sf/config"
)

var client = &http.Client{}

func SetHeaders(req *http.Request) {
	config.Load()
	Header := make(map[string]string)
	if config.Var.Cookie == "" {
		fmt.Println("Cookie is empty, please login first!")
		os.Exit(1)
	} else {
		Header["Cookie"] = config.Var.Cookie
	}
	Header["sf-minip-info"] = "minip_novel/1.0.70(android;11)/wxmp"
	Header["Content-Type"] = "application/json"
	if config.Var.UserName != "" || config.Var.Password != "" {
		Header["test-go"] = "cookie:" + config.Var.Cookie
		Header["Authorization"] = config.Var.Authorization
		Header["account-go"] = config.Var.UserName + "&" + config.Var.Password
	}
	fmt.Println(Header)
	for k, v := range Header {
		req.Header.Set(k, v)
	}
}

package request

import (
	"fmt"
	"net/http"
	"sf/config"
)

var client = &http.Client{}

func SetHeaders(req *http.Request) {
	config.Load()
	Header := make(map[string]string)
	if config.Var.Cookie == "" {
		fmt.Println("Cookie is empty, please login first!")
	} else {
		Header["Cookie"] = config.Var.Cookie
	}
	Header["sf-minip-info"] = "minip_novel/1.0.70(android;11)/wxmp"
	Header["Content-Type"] = "application/json"
	//Header["test-crawling"] = "golang/http"
	for k, v := range Header {
		req.Header.Set(k, v)
	}
}

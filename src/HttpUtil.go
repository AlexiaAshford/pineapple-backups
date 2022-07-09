package src

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"sf/setting"
)

var (
	client  = &http.Client{}
	WebSite = "https://api.sfacg.com/%v"
)

func SetHeaders(req *http.Request) {
	setting.Load()
	Header := make(map[string]string)
	Header["Host"] = "api.sfacg.com"
	Header["Connection"] = "keep-alive"
	Header["accept"] = "application/vnd.sfacg.api+json;version=1"
	Header["Cookie"] = setting.FormatCookie(setting.Var.Cookie[".SFCommunity"], setting.Var.Cookie["session_APP"])
	Header["User-Agent"] = setting.Var.UserAgent
	Header["authorization"] = setting.Var.Authorization
	Header["Content-type"] = "application/json; charset=UTF-8"
	Header["Accept-Encoding"] = "gzip,compress,br,deflate"

	for k, v := range Header {
		req.Header.Set(k, v)
	}
}

func Get(url string) []byte {
	if req, err := http.NewRequest("GET", fmt.Sprintf(WebSite, url), nil); err == nil {
		SetHeaders(req)
		if resp, err := client.Do(req); err == nil {
			if bodyText, ok := ioutil.ReadAll(resp.Body); ok == nil {
				//fmt.Println(string(bodyText))
				return bodyText
			}
		}
	}
	return nil
}

func POST(url, dataJson string) ([]byte, []*http.Cookie) {
	if request, err := http.NewRequest(
		"POST", fmt.Sprintf(WebSite, url), bytes.NewBuffer([]byte(dataJson)),
	); err != nil {
		fmt.Println(err)
	} else {
		SetHeaders(request)
		if response, ok := client.Do(request); ok == nil {
			if body, bodyError := ioutil.ReadAll(response.Body); bodyError == nil {
				return body, response.Cookies()
			}
		} else {
			fmt.Println(err)

		}
	}
	return nil, nil
}

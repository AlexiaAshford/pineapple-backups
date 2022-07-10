package request

import (
	"net/http"
	"sf/config"
)

var client = &http.Client{}

func SetHeaders(req *http.Request) {
	config.Load()
	Header := make(map[string]string)
	Header["Host"] = "api.sfacg.com"
	Header["Connection"] = "keep-alive"
	Header["accept"] = "application/vnd.sfacg.api+json;version=1"
	Header["Cookie"] = config.FormatCookie(config.Var.Cookie[".SFCommunity"], config.Var.Cookie["session_APP"])
	Header["User-Agent"] = config.Var.UserAgent
	Header["authorization"] = config.Var.Authorization
	Header["Content-type"] = "application/json; charset=UTF-8"
	Header["Accept-Encoding"] = "gzip,compress,br,deflate"
	for k, v := range Header {
		req.Header.Set(k, v)
	}
}

package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/VeronicaAlexia/pineapple-backups/config"
	"github.com/VeronicaAlexia/pineapple-backups/pkg/encryption"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type HttpUtils struct {
	url         string
	method      string
	cookie      []*http.Cookie
	response    *http.Request
	app_type    string
	query_data  *url.Values
	result_body []byte
}

func MustNewRequest(method, url string, data io.Reader) *http.Request {
	if request, err := http.NewRequest(method, url, data); err != nil {
		fmt.Println(config.Error("MustNewRequest", err, 93))
		return nil
	} else {
		return request
	}
}
func (is *HttpUtils) GetEncodeParams() *bytes.Reader {
	return bytes.NewReader([]byte(is.query_data.Encode()))
}
func (is *HttpUtils) GetResultBody() string {
	return string(is.result_body)
}

func (is *HttpUtils) GetCookie() []*http.Cookie {
	return is.cookie
}
func (is *HttpUtils) GetValue(key string) string {
	return is.query_data.Get(key)
}

func (is *HttpUtils) GetUrl() string {
	return is.url
}

func (is *HttpUtils) Add(key string, value string) *HttpUtils {
	is.query_data.Add(key, value)
	return is
}

func NewHttpUtils(api_url, method string) *HttpUtils {
	req := &HttpUtils{method: method, query_data: &url.Values{}}
	req.app_type = config.Vars.AppType
	if req.app_type == "cat" {
		req.url = CatWebSite + strings.ReplaceAll(api_url, CatWebSite, "")
		req.Add("login_token", config.Apps.Cat.Params.LoginToken).
			Add("account", config.Apps.Cat.Params.Account).
			Add("app_version", config.Apps.Cat.Params.AppVersion).
			Add("device_token", config.Apps.Cat.Params.DeviceToken)
	} else if req.app_type == "sfacg" {
		req.url = SFWebSite + strings.ReplaceAll(api_url, SFWebSite, "")
	} else {
		req.url = api_url
	}

	return req
}

func (is *HttpUtils) NEW_SET_THE_HEADERS() {
	HeaderCollection := make(map[string]string)
	HeaderCollection["Content-Type"] = "application/x-www-form-urlencoded"
	switch config.Vars.AppType {
	case "sfacg":
		HeaderCollection["sf-minip-info"] = "minip_novel/1.0.70(android;11)/wxmp"
		HeaderCollection["Authorization"] = Base64Bytes(config.Apps.Sfacg.UserName, config.Apps.Sfacg.Password)
		HeaderCollection["Cookie"] = config.Apps.Sfacg.Cookie
		HeaderCollection["account-sfacg"] = config.Apps.Sfacg.UserName + "&" + config.Apps.Sfacg.Password
	case "cat":
		HeaderCollection["User-Agent"] = "Android com.kuangxiangciweimao.novel 2.9.291"
		HeaderCollection["Cookie"] = "Account:" + config.Apps.Cat.Params.Account + ";" + config.Apps.Cat.Params.LoginToken
		HeaderCollection["Authorization"] = Base64Bytes(config.Apps.Cat.Params.Account, config.Apps.Cat.Params.LoginToken)

	default:
		fmt.Println(config.Vars.AppType, "AppType is invalid, please check config file")
		os.Exit(1)
	}
	for HeaderKey, HeaderValue := range HeaderCollection {
		is.response.Header.Set(HeaderKey, HeaderValue)
	}
}

func (is *HttpUtils) NewRequests() *HttpUtils {
	is.result_body = nil
	is.response = MustNewRequest(is.method, is.url, is.GetEncodeParams())
	is.NEW_SET_THE_HEADERS()
	if response, ok := http.DefaultClient.Do(is.response); ok == nil {
		is.cookie = response.Cookies()
		result_body, _ := io.ReadAll(response.Body)
		if config.Vars.AppType == "cat" && !strings.Contains(is.url, "jpg") {
			is.result_body = encryption.Decode(string(result_body), "")
		} else {
			is.result_body = result_body
		}
	} else {
		fmt.Println(config.Error(is.method+":"+is.url, ok, 67))
	}
	return is
}

func (is *HttpUtils) Unmarshal(s any) *HttpUtils {
	err := json.Unmarshal(is.result_body, s)
	if err != nil {
		fmt.Println(config.Error("json unmarshal", err, 18))
	}
	return is
}

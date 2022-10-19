package https

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/VeronicaAlexia/pineapple-backups/config"
	"github.com/VeronicaAlexia/pineapple-backups/src/encryption"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type HttpUtils struct {
	url        string
	method     string
	response   *http.Request
	AppType    string
	query_data *url.Values
	ResultBody []byte
}

func (is *HttpUtils) params() *bytes.Reader {
	return bytes.NewReader([]byte(is.query_data.Encode()))
}
func (is *HttpUtils) Add(key string, value string) *HttpUtils {
	if key != "" {
		is.query_data.Add(key, value)

	}
	return is
}
func (is *HttpUtils) Test() {
	fmt.Println(is.query_data.Encode())
	fmt.Println(string(is.ResultBody))
}

func NewHttpUtils(api_url, method string) *HttpUtils {
	req := &HttpUtils{method: method, query_data: &url.Values{}}
	req.AppType = config.Vars.AppType
	if req.AppType == "cat" {
		req.url = CatWebSite + strings.ReplaceAll(api_url, CatWebSite, "")
		req.Add("login_token", config.Apps.Cat.Params.LoginToken).
			Add("account", config.Apps.Cat.Params.Account).
			Add("app_version", config.Apps.Cat.Params.AppVersion).
			Add("device_token", config.Apps.Cat.Params.DeviceToken)
	} else if req.AppType == "sfacg" {
		req.url = SFWebSite + strings.ReplaceAll(api_url, SFWebSite, "")
	} else {
		req.url = api_url
	}

	return req
}

func MustNewRequest(method, url string, data io.Reader) *http.Request {
	if request, err := http.NewRequest(method, url, data); err != nil {
		fmt.Println(config.Error("MustNewRequest", err, 93))
		return nil
	} else {
		return request
	}
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
	is.ResultBody = nil
	is.response = MustNewRequest(is.method, is.url, is.params())
	is.NEW_SET_THE_HEADERS()
	if response, ok := http.DefaultClient.Do(is.response); ok == nil {
		result_body, _ := io.ReadAll(response.Body)
		if config.Vars.AppType == "cat" && !strings.Contains(is.url, "jpg") {
			is.ResultBody = encryption.Decode(string(result_body), "")
		} else {
			is.ResultBody = result_body
		}
	} else {
		fmt.Println(config.Error(is.method+":"+is.url, ok, 67))
	}
	return is
}

func (is *HttpUtils) Unmarshal(s any) *HttpUtils {
	err := json.Unmarshal(is.ResultBody, s)
	if err != nil {
		fmt.Println(config.Error("json unmarshal", err, 18))
	}
	return is
}

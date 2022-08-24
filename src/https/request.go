package https

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sf/cfg"
	"sf/src/hbooker/Encrypt"
	"strings"
)

var client = &http.Client{}

func JsonUnmarshal(response []byte, Struct any) any {
	err := json.Unmarshal(response, Struct)
	if err != nil {
		fmt.Println("json unmarshal error:", err)
	}
	return Struct
}

func QueryParams(url string, ParamsData map[string]string) string {
	var Params string
	queryRequisite := map[string]interface{}{
		"login_token":  cfg.Apps.Cat.Params.LoginToken,
		"account":      cfg.Apps.Cat.Params.Account,
		"app_version":  cfg.Apps.Cat.Params.AppVersion,
		"device_token": cfg.Apps.Cat.Params.DeviceToken,
	}
	for k, v := range ParamsData {
		Params += fmt.Sprintf("&%s=%s", k, v)
	}
	for k, v := range queryRequisite {
		Params += fmt.Sprintf("&%s=%s", k, v)
	}
	return url + "?" + Params
}

func SET_URL(url string, params map[string]string) string {
	switch cfg.Vars.AppType {
	case "cat":
		return CatWebSite + strings.Replace(QueryParams(url, params), CatWebSite, "", -1)
	case "sfacg":
		return SFWebSite + strings.Replace(url, SFWebSite, "", -1)
	case "happybooker":
		return CatWebSite + strings.Replace(QueryParams(url, params), CatWebSite, "", -1)
	default:
		return url
	}
}
func LoginSession(url string, dataJson []byte) ([]byte, []*http.Cookie) {
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(dataJson))
	if err != nil {
		fmt.Printf("LoginSession error:%v\n", err)
	} else {
		SetHeaders(request, false)
		if response, ok := client.Do(request); ok == nil {
			if body, err := io.ReadAll(response.Body); err == nil {
				return body, response.Cookies()
			}
		} else {
			fmt.Println("client.Do error:", err)
		}
	}

	return nil, nil
}

func Request(method string, url string) ([]byte, error) {
	if method != "GET" && method != "POST" && method != "PUT" {
		panic("Error: method must be GET or POST or PUT, but now is " + method)
	}
	if request, err := http.NewRequest(method, url, nil); err != nil {
		fmt.Printf("NewRequest %v error:%v\n", method, err)
	} else {
		SetHeaders(request, true)
		if response, ok := client.Do(request); ok == nil {
			// delete ioutil.ReadAll and use io.ReadAll instead
			if body, bodyError := io.ReadAll(response.Body); bodyError == nil {
				return body, nil
			}
		}
	}
	return nil, errors.New("request error:" + method + "url:" + url)
}

func Get(url string, structural any, params map[string]string) any {
	if cfg.Vars.AppType == "cat" {
		if result, ok := Request("POST", SET_URL(url, params)); ok == nil {
			return JsonUnmarshal(Encrypt.Decode(string(result), ""), structural)
		} else {
			fmt.Println(ok)
		}
	} else if cfg.Vars.AppType == "sfacg" {
		if result, ok := Request("GET", SET_URL(url, params)); ok == nil {
			return JsonUnmarshal(result, structural)
		} else {
			fmt.Println(ok)
		}
	}
	return nil
}
func GetCover(imgUrl string) []byte {
	if res, err := Request("GET", imgUrl); err == nil {
		return res
	} else {
		fmt.Println(err)
	}
	return nil
}

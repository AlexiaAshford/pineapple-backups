package https

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/VeronicaAlexia/pineapple-backups/config"
	"github.com/VeronicaAlexia/pineapple-backups/src/app/hbooker/Encrypt"
	"github.com/VeronicaAlexia/pineapple-backups/struct/sfacg_structs"
	"io"
	"net/http"
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

func SET_URL(url string) string {
	switch config.Vars.AppType {
	case "cat":
		return CatWebSite + strings.ReplaceAll(url, CatWebSite, "")
	case "sfacg":
		return SFWebSite + strings.ReplaceAll(url, SFWebSite, "")
	case "happybooker":
		return HappyWebSite + strings.ReplaceAll(url, HappyWebSite, "")
	default:
		return url
	}
}
func Login(url string, dataJson []byte) (*sfacg_structs.Login, []*http.Cookie) {
	request, err := http.NewRequest("POST", SET_URL(url), bytes.NewBuffer(dataJson))
	if err != nil {
		fmt.Printf("Login session error:%v\n", err)
		return nil, nil
	}
	SET_THE_HEADERS(request, false)
	response, ok := client.Do(request)
	if ok != nil {
		return nil, nil
	}
	body, _ := io.ReadAll(response.Body)
	return JsonUnmarshal(body, &sfacg_structs.Login{}).(*sfacg_structs.Login), response.Cookies()
}

func Request(method string, url string) ([]byte, error) {
	if method != "GET" && method != "POST" && method != "PUT" {
		panic("Error: method must be GET or POST or PUT, but now is " + method)
	}
	if request, err := http.NewRequest(method, url, nil); err != nil {
		fmt.Printf("NewRequest %v error:%v\n", method, err)
	} else {
		SET_THE_HEADERS(request, true)
		if response, ok := client.Do(request); ok == nil {
			// delete ioutil.ReadAll and use io.ReadAll instead
			if body, bodyError := io.ReadAll(response.Body); bodyError == nil {
				return body, nil
			}
		}
	}
	return nil, errors.New("request error:" + method + "url:" + url)
}

func Get(url string, structural any) any {
	if config.Vars.AppType == "cat" {
		if result, ok := Request("POST", SET_URL(url)); ok == nil {
			//fmt.Println(string(Encrypt.Decode(string(result), "")))
			return JsonUnmarshal(Encrypt.Decode(string(result), ""), structural)
		} else {
			fmt.Println(ok)
		}
	} else if config.Vars.AppType == "sfacg" {
		if result, ok := Request("GET", SET_URL(url)); ok == nil {
			return JsonUnmarshal(result, structural)
		} else {
			fmt.Println(ok)
		}
	} else {
		fmt.Println("not support, please use cat or sfacg, now is " + config.Vars.AppType)
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

//
//func QueryParams(url string, ParamsData map[string]string) string {
//	var Params string
//	if config.Vars.AppType == "cat" {
//		queryRequisite := map[string]interface{}{
//			"login_token":  config.Apps.Cat.Params.LoginToken,
//			"account":      config.Apps.Cat.Params.Account,
//			"app_version":  config.Apps.Cat.Params.AppVersion,
//			"device_token": config.Apps.Cat.Params.DeviceToken,
//		}
//		for k, v := range queryRequisite {
//			Params += fmt.Sprintf("&%s=%s", k, v)
//		}
//	}
//	for k, v := range ParamsData {
//		Params += fmt.Sprintf("&%s=%s", k, v)
//	}
//	return url + "?" + Params
//}

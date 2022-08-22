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
)

var client = &http.Client{}

func JsonUnmarshal(response []byte, Struct any) any {
	err := json.Unmarshal(response, Struct)
	if err != nil {
		fmt.Println("json unmarshal error:", err)
	}
	return Struct
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

func Get(method string, url string) []byte {
	if cfg.Vars.AppType == "cat" {
		if result, ok := Request(method, url); ok == nil {
			return Encrypt.Decode(string(result), "")
		} else {
			fmt.Println(ok)
		}
	} else if cfg.Vars.AppType == "sfacg" {
		if result, ok := Request(method, url); ok == nil {
			return result
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

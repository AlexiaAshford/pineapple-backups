package https

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

func Request(method string, url string, dataJson string) ([]byte, []*http.Cookie) {

	if method != "GET" && method != "POST" && method != "PUT" {
		panic("Error: method must be GET or POST or PUT, but now is " + method)
	}
	if request, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(dataJson))); err != nil {
		fmt.Printf("NewRequest %v error:%v\n", method, err)
	} else {
		if url == "https://minipapi.sfacg.com/pas/mpapi/sessions" {
			// login request no need to test cookie
			SetHeaders(request, false)
		} else {
			SetHeaders(request, true)
		}
		if response, ok := client.Do(request); ok == nil {
			// delete ioutil.ReadAll and use io.ReadAll instead
			if body, bodyError := io.ReadAll(response.Body); bodyError == nil {
				return body, response.Cookies()
			}
		} else {
			fmt.Println(err)
		}
	}
	return nil, nil
}

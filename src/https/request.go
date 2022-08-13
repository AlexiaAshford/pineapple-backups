package https

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func Request(method string, url string, dataJson string) ([]byte, []*http.Cookie) {
	var client = &http.Client{}

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

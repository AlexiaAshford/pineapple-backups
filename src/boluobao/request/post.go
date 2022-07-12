package request

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func POST(url, dataJson string) ([]byte, []*http.Cookie) {
	api := fmt.Sprintf("https://minipapi.sfacg.com/pas/mpapi/%v", url)
	if request, err := http.NewRequest("POST", api, bytes.NewBuffer([]byte(dataJson))); err != nil {
		fmt.Println("NewRequest post error:", err)
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

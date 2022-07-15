package request

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func Get(url string) []byte {
	api := fmt.Sprintf("%v", url)
	if req, err := http.NewRequest("GET", api, nil); err == nil {
		SetHeaders(req, true)
		if resp, ok := client.Do(req); ok == nil {
			if bodyText, ReadBody := ioutil.ReadAll(resp.Body); ReadBody == nil {
				return bodyText
			} else {
				fmt.Println("Read Body Error:", ReadBody)
			}
		} else {
			fmt.Println("Get error:", err)
		}
	}
	return nil
}

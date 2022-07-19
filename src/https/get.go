package https

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func Get(Url string, Retry int) []byte {
	if req, err := http.NewRequest("GET", Url, nil); err == nil {
		SetHeaders(req, true)
		if resp, ok := client.Do(req); ok == nil {
			if bodyText, ReadBody := ioutil.ReadAll(resp.Body); ReadBody == nil {
				return bodyText
			} else {
				fmt.Println("Read Body Error:", ReadBody)
			}
		} else {
			if Retry >= 5 {
				fmt.Println("Get error:", err)
				os.Exit(1)
			} else {
				return Get(Url, Retry+1)
			}
		}
	}
	return nil
}

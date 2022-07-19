package https

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func Get(url string, Retry int) []byte {
	if req, err := http.NewRequest("GET", url, nil); err == nil {
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
				return Get(url, Retry+1)
			}
		}
	}
	return nil
}

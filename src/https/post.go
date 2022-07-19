package https

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func POST(url, dataJson string) ([]byte, []*http.Cookie) {
	if request, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(dataJson))); err != nil {
		fmt.Println("NewRequest post error:", err)
	} else {
		SetHeaders(request, false)
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

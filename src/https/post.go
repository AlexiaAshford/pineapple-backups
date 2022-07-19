package https

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func POST(URL string, dataJson string) ([]byte, []*http.Cookie) {
	params := bytes.NewBuffer([]byte(dataJson))
	if request, err := http.NewRequest("POST", URL, params); err != nil {
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

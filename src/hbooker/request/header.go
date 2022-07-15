package request

import (
	"net/http"
)

var client = &http.Client{}

func SetHeaders(req *http.Request, isGet bool) {
	if isGet {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		req.Header.Set("Content-Type", "application/json")
	}
}

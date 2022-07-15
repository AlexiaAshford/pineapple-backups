package sfacg_structs

import "net/http"

type Login struct {
	Status struct {
		HTTPCode  int         `json:"httpCode"`
		ErrorCode int         `json:"errorCode"`
		MsgType   int         `json:"msgType"`
		Msg       interface{} `json:"msg"`
	} `json:"status"`
	Cookie []*http.Cookie
	Data   interface{} `json:"data"`
}

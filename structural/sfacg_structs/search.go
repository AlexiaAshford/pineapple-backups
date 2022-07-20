package sfacg_structs

type Search struct {
	Status struct {
		HTTPCode  int         `json:"httpCode"`
		ErrorCode int         `json:"errorCode"`
		MsgType   int         `json:"msgType"`
		Msg       interface{} `json:"msg"`
	} `json:"status"`
	Data struct {
		Novels []BookInfoData `json:"novels"`
	} `json:"data"`
}

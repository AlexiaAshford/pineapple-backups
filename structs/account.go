package structs

type Account struct {
	Status struct {
		HTTPCode  int    `json:"httpCode"`
		ErrorCode int    `json:"errorCode"`
		MsgType   int    `json:"msgType"`
		Msg       string `json:"msg"`
	} `json:"status"`
	Data struct {
		UserName     string `json:"userName"`
		NickName     string `json:"nickName"`
		Email        string `json:"email"`
		AccountID    int    `json:"accountId"`
		RoleName     string `json:"roleName"`
		FireCoin     int    `json:"fireCoin"`
		Avatar       string `json:"avatar"`
		IsAuthor     bool   `json:"isAuthor"`
		PhoneNum     string `json:"phoneNum"`
		RegisterDate string `json:"registerDate"`
	} `json:"data"`
}

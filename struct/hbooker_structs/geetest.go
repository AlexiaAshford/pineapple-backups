package hbooker_structs

type GeetestStruct struct {
	Code int         `json:"code"`
	Tip  string      `json:"tip"`
	Data GeetestData `json:"data"`
}
type GeetestData struct {
	NeedUseGeetest int `json:"need_use_geetest"`
	CodeLen        int `json:"code_len"`
}
type GeetestChallenge struct {
	Success    int    `json:"success"`
	Gt         string `json:"gt"`
	Challenge  string `json:"challenge"`
	NewCaptcha bool   `json:"new_captcha"`
}

package hbooker

import (
	"github.com/go-resty/resty/v2"
	"github.com/gookit/color"
	"github.com/tidwall/gjson"
	"sf/src/hbooker/Geetest"
	"strconv"
	"time"
)

// GetGtChallenge 从Demo获取gt和challenge
func GetGtChallenge() {
	url := "https://www.geetest.com/demo/gt/register-slide?t=" + strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	client := resty.New()
	res, err := client.R().
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36 Edg/91.0.864.64").
		Get(url)
	if err != nil {
		color.Errorln("GetGt Error: ", err)
		panic(err)
	}
	response := res.String()
	gt := gjson.Get(response, "gt").String()
	challenge := gjson.Get(response, "challenge").String()
	g := Geetest.Geetest{
		GT:        gt,
		Challenge: challenge,
	}
	status, CaptchaType, errorDetail := Geetest.GetFullBG(&g)
	if status == "success" {
		color.Infoln("验证码类型：", CaptchaType, "")
	} else {
		color.Errorln("获取图片失败   Err: ", status, " ErrorDetail:", errorDetail)
	}
}

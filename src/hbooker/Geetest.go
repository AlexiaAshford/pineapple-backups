package hbooker

import (
	"github.com/go-resty/resty/v2"
	"github.com/gookit/color"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"sf/src/hbooker/Encrypt"
	"strconv"
	"time"
)

// Geetest 未来会精简这个结构体
type Geetest struct {
	GT        string
	Challenge string
	S         string
	C         []int64
	FullBG    string
	BG        string
	SecKey    string
}

var client = resty.New()

// GetSC 流程的第一步，获取s,c，二代极验获取全图
func GetSC(g *Geetest) int64 {
	temp := `{"gt":"","challenge":"","offline":false,"new_captcha":true,"product":"custom","width":"300px","next_width":"278px","area":"#area","bg_color":"gray","https":true,"protocol":"https://","click":"/static/js/click.2.8.1.js","slide":"/static/js/slide.7.6.0.js","pencil":"/static/js/pencil.1.0.3.js","fullpage":"/static/js/fullpage.8.7.9.js","beeline":"/static/js/beeline.1.0.1.js","aspect_radio":{"beeline":50,"click":128,"slide":103,"voice":128,"pencil":128},"static_servers":["static.geetest.com/","dn-staticdown.qbox.me/"],"type":"fullpage","maze":"/static/js/maze.1.0.1.js","geetest":"/static/js/geetest.6.0.9.js","voice":"/static/js/voice.1.2.0.js","cc":8,"ww":true,"i":"6322!!7608!!CSS1Compat!!1!!-1!!-1!!-1!!-1!!-1!!-1!!-1!!-1!!-1!!2!!3!!-1!!-1!!-1!!-1!!-1!!-1!!-1!!-1!!-1!!-1!!1!!-1!!-1!!-1!!10!!44!!0!!0!!737!!784!!1687!!888!!zh-CN!!zh-CN,zh!!-1!!1.5!!24!!Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Safari/537.36!!1!!1!!1706!!960!!1707!!960!!1!!1!!1!!-1!!Linux x86_64!!0!!-8!!0a93728bbc5be4241e024b729f6c5e5d!!3f4cf59bdc7d0da206835d0e495e258a!!internal-pdf-viewer,mhjfbmdgcfjbbpaeojofohoefgiehjai!!0!!-1!!0!!8!!Arial,BitstreamVeraSansMono,Courier,CourierNew,Helvetica,Monaco,Times,TimesNewRoman,Wingdings,Wingdings2,Wingdings3!!1563788965804!!-1,-1,2,1,2,0,15,0,49,1,4,4,10,70,71,73,74,74,74,-1!!-1!!-1!!12!!-1!!-1!!-1!!5!!false!!false"}`
	value, _ := sjson.Set(temp, "gt", g.GT)
	text, _ := sjson.Set(value, "challenge", g.Challenge)
	url := "https://api.geetest.com/get.php"
	w := g.CalW(text, false)
	res, err := client.R().
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36 Edg/91.0.864.64").
		SetQueryParam("gt", g.GT).
		SetQueryParam("challenge", g.Challenge).
		SetQueryParam("w", w).
		Get(url)
	if err != nil {
		color.Errorln("StepOne Error: ", err)
		return 0
	} else {
		response := res.String()
		//color.Green.Println(response) // 输出结果，可注释

		if gjson.Get(response, "c").Exists() {
			a := gjson.Get(response, "c").Array()
			var b []int64
			for _, v := range a {
				b = append(b, v.Int())
			}
			g.C = b
			g.S = gjson.Get(response, "s").String()
			g.FullBG = "https://static.geetest.com/" + gjson.Get(response, "fullbg").String()
			g.BG = "https://static.geetest.com/" + gjson.Get(response, "bg").String()
			return 1
		} else {
			a := gjson.Get(response, "data.c").Array()
			var b []int64
			for _, v := range a {
				b = append(b, v.Int())
			}
			g.C = b
			g.S = gjson.Get(response, "data.s").String()
			return 2
		}
	}
}

// GetCaptchaType 流程的第二步，获取验证码类型
func GetCaptchaType(g *Geetest) (string, string, string) {
	//w := g.CalW(Function.CalA(g.C, g.S, g.GT, g.Challenge), true)
	w := g.CalW("{}", true)
	url := "https://api.geetest.com/ajax.php"
	res, err := client.R().
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36 Edg/91.0.864.64").
		SetQueryParam("gt", g.GT).
		SetQueryParam("challenge", g.Challenge).
		SetQueryParam("w", w).
		Get(url)
	if err != nil {
		color.Errorln("StepOne Error: ", err)
		return "error", "", "GetCaptchaType Request Failed!"
	} else {
		response := res.String()
		//color.Green.Println(response) // 输出结果，可注释

		status := gjson.Get(response, "status").String()
		CaptchaType := gjson.Get(response, "data.result").String()
		errorDetail := gjson.Get(response, "error").String()
		return status, CaptchaType, errorDetail
	}
}

// GetFullBG 流程的第三步，极验三代获取背景图片，返回信息为 函数状态，验证码类型，错误详情
func GetFullBG(g *Geetest) (string, string, string) {
	success := GetSC(g)
	if success == 0 { // 失败
		return "failed", "", "GetSC Request Failed!"
	} else if success == 1 { // 极验二代，不需要再获取图片
		return "success", "Geetest2", ""
	} else if success == 2 {
		status, CaptchaType, errorDetail := GetCaptchaType(g)
		if status == "error" {
			return status, CaptchaType, errorDetail
		} else {
			url := "https://api.geetest.com/get.php"
			res, err := client.R().
				SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36 Edg/91.0.864.64").
				SetQueryParam("gt", g.GT).
				SetQueryParam("challenge", g.Challenge).
				SetQueryParam("callback", "geetest_"+strconv.FormatInt(time.Now().UnixNano()/1e6, 10)).
				Get(url)
			if err != nil {
				color.Errorln("StepOne Error: ", err)
				return "error", "", "GetFullBG Request Failed!"
			} else {
				response := res.String()
				//color.Green.Println(response) // 输出结果，可注释
				a := gjson.Get(response, "c").Array()
				var b []int64
				for _, v := range a {
					b = append(b, v.Int())
				}
				g.C = b
				g.S = gjson.Get(response, "s").String()
				g.FullBG = "https://static.geetest.com/" + gjson.Get(response, "fullbg").String()
				g.BG = "https://static.geetest.com/" + gjson.Get(response, "bg").String()
				g.Challenge = gjson.Get(response, "challenge").String()
				return status, CaptchaType, errorDetail
			}
		}
	}
	return "unknownError", "", "GetFullBG直接跳过了计算"
}

func Slide(g *Geetest) {
	// 下载fullbg bg
	// 计算距离（简化处理）
	var distance int64 = 123
	// 计算轨迹（简化处理）
	trace := [][]int64{{distance, distance, distance}, {}, {}, {}}

	lastTrace := trace[len(trace)-1]
	c := lastTrace[0]
	passTime := lastTrace[len(lastTrace)-1]
	aa := Encrypt.CalAA(trace, g.C, g.S)
	w := Encrypt.GetRequestW(g.GT, g.Challenge, aa, strconv.FormatInt(passTime, 10), c)
	res, err := client.R().
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36 Edg/91.0.864.64").
		SetQueryParam("gt", g.GT).
		SetQueryParam("challenge", g.Challenge).
		SetQueryParam("w", w).
		SetQueryParam("callback", "geetest_"+strconv.FormatInt(time.Now().UnixNano()/1e6, 10)).
		Get("https://api.geetest.com/ajax.php")
	if err != nil {
		color.Errorln("Slide Request Error: ", err)
	} else {
		response := res.String()
		//color.Green.Println(response) // 输出结果，可注释

		success := gjson.Get(response, "success").String()
		validate := gjson.Get(response, "validate").String()
		score := gjson.Get(response, "score").String()
		color.Infoln("成功：", success, "\tValidate: ", validate, "\tScore: ", score)
	}
}

// CalW 计算validate前需要提交的W参数
func (g Geetest) CalW(text string, flag bool) string {
	EncSecKey := ""
	if !flag {
		secKey := Encrypt.CreateSecretKey()
		if len(secKey) == 16 {
			EncSecKey = Encrypt.RSAEncrypt(secKey)
			encTextByte := Encrypt.AESEncrypt([]byte(text), secKey)
			encText := Encrypt.BytesToString(encTextByte)
			return encText + EncSecKey
		} else {
			return g.CalW(text, flag)
		}
	}
	encTextByte := Encrypt.AESEncrypt([]byte(text), Encrypt.CreateSecretKey())
	encText := Encrypt.BytesToString(encTextByte)
	return encText + EncSecKey
}

package hbooker_structs

var Login = struct {
	Code string      `json:"code"`
	Data LoginData   `json:"data"`
	Tip  interface{} `json:"tip"`
}{}

type LoginData struct {
	LoginToken string     `json:"login_token"`
	UserCode   string     `json:"user_code"`
	ReaderInfo ReaderInfo `json:"reader_info"`
	PropInfo   PropInfo   `json:"prop_info"`
	IsSetYoung string     `json:"is_set_young"`
}

type PropInfo struct {
	RESTGiftHlb     string `json:"rest_gift_hlb"`
	RESTHlb         string `json:"rest_hlb"`
	RESTYp          string `json:"rest_yp"`
	RESTRecommend   string `json:"rest_recommend"`
	RESTTotalBlade  string `json:"rest_total_blade"`
	RESTMonthBlade  string `json:"rest_month_blade"`
	RESTTotal100    string `json:"rest_total_100"`
	RESTTotal588    string `json:"rest_total_588"`
	RESTTotal1688   string `json:"rest_total_1688"`
	RESTTotal5000   string `json:"rest_total_5000"`
	RESTTotal10000  string `json:"rest_total_10000"`
	RESTTotal100000 string `json:"rest_total_100000"`
	RESTTotal50000  string `json:"rest_total_50000"`
	RESTTotal160000 string `json:"rest_total_160000"`
}

type ReaderInfo struct {
	ReaderID       string        `json:"reader_id"`
	Account        string        `json:"account"`
	IsBind         string        `json:"is_bind"`
	IsBindQq       string        `json:"is_bind_qq"`
	IsBindWeixin   string        `json:"is_bind_weixin"`
	IsBindHuawei   string        `json:"is_bind_huawei"`
	IsBindApple    string        `json:"is_bind_apple"`
	PhoneNum       string        `json:"phone_num"`
	MobileVal      string        `json:"mobileVal"`
	Email          string        `json:"email"`
	License        string        `json:"license"`
	ReaderName     string        `json:"reader_name"`
	AvatarURL      string        `json:"avatar_url"`
	AvatarThumbURL string        `json:"avatar_thumb_url"`
	BaseStatus     string        `json:"base_status"`
	ExpLV          string        `json:"exp_lv"`
	ExpValue       string        `json:"exp_value"`
	Gender         string        `json:"gender"`
	VipLV          string        `json:"vip_lv"`
	VipValue       string        `json:"vip_value"`
	IsAuthor       string        `json:"is_author"`
	IsUploader     string        `json:"is_uploader"`
	BookAge        string        `json:"book_age"`
	CategoryPrefer []interface{} `json:"category_prefer"`
	UsedDecoration []interface{} `json:"used_decoration"`
	Rank           string        `json:"rank"`
	Ctime          string        `json:"ctime"`
}

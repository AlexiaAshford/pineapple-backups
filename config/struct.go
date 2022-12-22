package config

type AppConfig struct {
	Hbooker HbookerCommonParams `json:"common_params"`
	Sfacg   BoluobaoConfig      `json:"sfacg_config"`
	Config  AccountConfig       `json:"account_config"`
}

type AccountConfig struct {
	ConfigName string
	OutputName string
	CoverFile  string
	AppKey     string
	DeviceId   string
	AppType    string
	ThreadNum  int
	MaxRetry   int
	Epub       bool
}

type BoluobaoConfig struct {
	UserName     string `json:"account"`
	Password     string `json:"password"`
	Cookie       string `json:"cookie"`
	WeChatCookie string `json:"we_chat_cookie"`
}

type HbookerCommonParams struct {
	LoginToken  string `json:"login_token"`
	Account     string `json:"account"`
	AppVersion  string `json:"app_version"`
	DeviceToken string `json:"device_token"`
}

type Books struct {
	NovelName  string
	NovelID    string
	IsFinish   bool
	MarkCount  string
	NovelCover string
	AuthorName string
	CharCount  string
	SignStatus string
}
type MyBookInfoJsonPro struct {
	Book       Books
	NewBooks   map[string]string
	OutputPath string
	ConfigPath string
	CoverPath  string
	//BackupsPath  string
	BookInfoList []Books
	DownloadList []string
}

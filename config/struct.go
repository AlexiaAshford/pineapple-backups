package config

var Apps = AppConfig{}

type AppConfig struct {
	Hbooker HbookerCommonParams `json:"common_params"`
	Sfacg   BoluobaoConfig      `json:"sfacg_config"`
	Config  ScriptConfig        `json:"script_config"`
}

type ScriptConfig struct {
	ConfigName string `json:"config_name"`
	OutputName string `json:"output_name"`
	CoverFile  string `json:"cover_file"`
	DeviceId   string `json:"device_id"`
	AppType    string `json:"app_type"`
	ThreadNum  int    `json:"thread_num"`
	MaxRetry   int    `json:"max_retry"`
	Epub       bool   `json:"epub"`
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
type BookInfo struct {
	Book       Books
	NewBooks   map[string]string
	OutputPath string
	ConfigPath string
	CoverPath  string
	//BackupsPath  string
	BookInfoList []Books
	DownloadList []string
}

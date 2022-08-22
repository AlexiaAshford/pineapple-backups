package _struct

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
type MyConfigPro struct {
	ConfigFile      string
	SaveFile        string
	AppType         string
	MaxThreadNumber int
	MaxRetry        int
}

type MyAppPro struct {
	Sfacg SfacgApp `json:"sfacg"`
	Cat   CatApp   `json:"cat"`
}

type SfacgApp struct {
	UserName  string `json:"account"`
	Password  string `json:"password"`
	Cookie    string `json:"Cookie"`
	UserAgent string `json:"user-agent"`
}

type CatApp struct {
	Params struct {
		LoginToken  string `json:"login_token"`
		Account     string `json:"account"`
		AppVersion  string `json:"app_version"`
		DeviceToken string `json:"device_token"`
	} `json:"common_params"`
	UserAgent string `json:"user-agent"`
}
type MyBookInfoJsonPro struct {
	BookInfo     Books
	BookInfoList []Books
}

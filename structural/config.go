package structural

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
type MyJsonPro struct {
	ConfigFile      string `json:"config_file"`
	SaveFile        string `json:"save_file"`
	AppType         string `json:"app_type"`
	MaxThreadNumber int    `json:"max_thread_number"`
	MaxRetry        int    `json:"max_retry"`
	Sfacg           struct {
		UserName  string `json:"account"`
		Password  string `json:"password"`
		Cookie    string `json:"Cookie"`
		UserAgent string `json:"user-agent"`
	} `json:"sfacg"`
	Cat struct {
		Params struct {
			LoginToken  string `json:"login_token"`
			Account     string `json:"account"`
			AppVersion  string `json:"app_version"`
			DeviceToken string `json:"device_token"`
		} `json:"common_params"`
		UserAgent string `json:"user-agent"`
	} `json:"cat"`
}

type MyBookInfoJsonPro struct {
	BookInfo     Books
	BookInfoList []Books
}

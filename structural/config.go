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
	Sfacg           struct {
		UserName string `json:"UserName"`
		Password string `json:"Password"`
		Cookie   string `json:"Cookie"`
	} `json:"sfacg"`
	Cat struct {
		Params struct {
			LoginToken  string `json:"login_token"`
			Account     string `json:"account"`
			AppVersion  string `json:"app_version"`
			DeviceToken string `json:"device_token"`
		} `json:"common_params"`
		UserAgent string `json:"UserAgent"`
	} `json:"cat"`
}

type MyBookInfoJsonPro struct {
	BookInfo     Books
	BookInfoList []Books
}

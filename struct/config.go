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

type MyAppPro struct {
	Sfacg  SfacgApp    `json:"sfacg"`
	Cat    CatApp      `json:"cat"`
	Config MyConfigPro `json:"config"`
}

type SfacgApp struct {
	UserName string `json:"account"`
	Password string `json:"password"`
	Cookie   string `json:"Cookie"`
}

type CatApp struct {
	Params struct {
		LoginToken  string `json:"login_token"`
		Account     string `json:"account"`
		AppVersion  string `json:"app_version"`
		DeviceToken string `json:"device_token"`
	} `json:"common_params"`
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

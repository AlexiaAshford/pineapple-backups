package sfacg_structs

type BookInfoData struct {
	AuthorID       int     `json:"authorId"`
	LastUpdateTime string  `json:"lastUpdateTime"`
	MarkCount      int     `json:"markCount"`
	NovelCover     string  `json:"novelCover"`
	BgBanner       string  `json:"bgBanner"`
	NovelID        int     `json:"novelId"`
	NovelName      string  `json:"novelName"`
	Point          float64 `json:"point"`
	IsFinish       bool    `json:"isFinish"`
	AuthorName     string  `json:"authorName"`
	CharCount      int     `json:"charCount"`
	ViewTimes      int     `json:"viewTimes"`
	TypeID         int     `json:"typeId"`
	AllowDown      bool    `json:"allowDown"`
	SignStatus     string  `json:"signStatus"`
	CategoryID     int     `json:"categoryId"`
}

var BookInfo = struct {
	Status struct {
		HTTPCode  int         `json:"httpCode"`
		ErrorCode int         `json:"errorCode"`
		MsgType   int         `json:"msgType"`
		Msg       interface{} `json:"msg"`
	} `json:"status"`
	Data BookInfoData `json:"data"`
}{}

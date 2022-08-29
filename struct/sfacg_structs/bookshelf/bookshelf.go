package bookshelf

type InfoData struct {
	Status struct {
		HTTPCode  int         `json:"httpCode"`
		ErrorCode int         `json:"errorCode"`
		MsgType   int         `json:"msgType"`
		Msg       interface{} `json:"msg"`
	} `json:"status"`
	Data []struct {
		AccountID  int    `json:"accountId"`
		PocketID   int    `json:"pocketId"`
		Name       string `json:"name"`
		TypeID     int    `json:"typeId"`
		CreateTime string `json:"createTime"`
		IsFull     bool   `json:"isFull"`
		CanModify  bool   `json:"canModify"`
		Expand     struct {
			Novels []NovelsList `json:"novels"`
		} `json:"expand"`
	} `json:"data"`
}
type NovelsList struct {
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
	IsSensitive    bool    `json:"isSensitive"`
	Expand         struct {
		Discount int `json:"discount"`
	} `json:"expand"`
	IsSticky       bool        `json:"isSticky"`
	StickyDateTime interface{} `json:"stickyDateTime"`
	MarkDateTime   string      `json:"markDateTime"`
}

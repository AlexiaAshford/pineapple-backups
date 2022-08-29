package bookshelf

type InfoData struct {
	Status Status `json:"status"`
	Data   []Data `json:"data"`
}
type Status struct {
	HTTPCode  int         `json:"httpCode"`
	ErrorCode int         `json:"errorCode"`
	MsgType   int         `json:"msgType"`
	Msg       interface{} `json:"msg"`
}
type Expand struct {
	Novels []interface{} `json:"novels"`
}
type Albums struct {
	AlbumID         int         `json:"albumId"`
	NovelID         int         `json:"novelId"`
	AuthorID        int         `json:"authorId"`
	LatestChapterID int         `json:"latestChapterId"`
	VisitTimes      int         `json:"visitTimes"`
	Name            string      `json:"name"`
	LastUpdateTime  string      `json:"lastUpdateTime"`
	CoverBig        string      `json:"coverBig"`
	CoverSmall      string      `json:"coverSmall"`
	CoverMedium     string      `json:"coverMedium"`
	Expand          Expand      `json:"expand"`
	IsSticky        bool        `json:"isSticky"`
	StickyDateTime  interface{} `json:"stickyDateTime"`
	MarkDateTime    string      `json:"markDateTime"`
}
type Data struct {
	AccountID  int    `json:"accountId"`
	PocketID   int    `json:"pocketId"`
	Name       string `json:"name"`
	TypeID     int    `json:"typeId"`
	CreateTime string `json:"createTime"`
	IsFull     bool   `json:"isFull"`
	CanModify  bool   `json:"canModify"`
}

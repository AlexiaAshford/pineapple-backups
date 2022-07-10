package structs

type Catalogue struct {
	Status Status `json:"status"`
	Data   Data   `json:"data"`
}
type Status struct {
	HTTPCode  int         `json:"httpCode"`
	ErrorCode int         `json:"errorCode"`
	MsgType   int         `json:"msgType"`
	Msg       interface{} `json:"msg"`
}
type ChapterList struct {
	ChapID              int         `json:"chapId"`
	NovelID             int         `json:"novelId"`
	VolumeID            int         `json:"volumeId"`
	NeedFireMoney       int         `json:"needFireMoney"`
	OriginNeedFireMoney int         `json:"originNeedFireMoney"`
	CharCount           int         `json:"charCount"`
	RowNum              int         `json:"rowNum"`
	ChapOrder           int         `json:"chapOrder"`
	Title               string      `json:"title"`
	Content             interface{} `json:"content"`
	Sno                 float64     `json:"sno"`
	IsVip               bool        `json:"isVip"`
	AddTime             string      `json:"AddTime"`
	UpdateTime          interface{} `json:"updateTime"`
	CanUnlockWithAd     bool        `json:"canUnlockWithAd"`
	IsRubbish           bool        `json:"isRubbish"`
	AuditStatus         int         `json:"auditStatus"`
}
type VolumeList struct {
	VolumeID    int           `json:"volumeId"`
	Title       string        `json:"title"`
	Sno         float64       `json:"sno"`
	ChapterList []ChapterList `json:"chapterList"`
}
type Data struct {
	NovelID        int          `json:"novelId"`
	LastUpdateTime string       `json:"lastUpdateTime"`
	VolumeList     []VolumeList `json:"volumeList"`
}

package structs

type Content struct {
	Status struct {
		HTTPCode  int         `json:"httpCode"`
		ErrorCode int         `json:"errorCode"`
		MsgType   int         `json:"msgType"`
		Msg       interface{} `json:"msg"`
	} `json:"status"`
	Data struct {
		ChapID     int    `json:"chapId"`
		NovelID    int    `json:"novelId"`
		VolumeID   int    `json:"volumeId"`
		CharCount  int    `json:"charCount"`
		RowNum     int    `json:"rowNum"`
		ChapOrder  int    `json:"chapOrder"`
		Title      string `json:"title"`
		AddTime    string `json:"addTime"`
		UpdateTime string `json:"updateTime"`
		Sno        int    `json:"sno"`
		IsVip      bool   `json:"isVip"`
		Expand     struct {
			NeedFireMoney       int         `json:"needFireMoney"`
			OriginNeedFireMoney int         `json:"originNeedFireMoney"`
			Content             string      `json:"content"`
			Tsukkomi            interface{} `json:"tsukkomi"`
			ChatLines           interface{} `json:"chatLines"`
		} `json:"expand"`
		Ntitle      string `json:"ntitle"`
		IsRubbish   bool   `json:"isRubbish"`
		AuditStatus int    `json:"auditStatus"`
	} `json:"data"`
}

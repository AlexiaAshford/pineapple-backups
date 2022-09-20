package task

type Task struct {
	Status struct {
		HttpCode  int         `json:"httpCode"`
		ErrorCode int         `json:"errorCode"`
		MsgType   int         `json:"msgType"`
		Msg       interface{} `json:"msg"`
	} `json:"status"`
	Data []struct {
		RecordId    int    `json:"recordId"`
		TaskId      int    `json:"taskId"`
		RequireNum  int    `json:"requireNum"`
		CompleteNum int    `json:"completeNum"`
		Status      int    `json:"status"`
		Name        string `json:"name"`
		Desc        string `json:"desc"`
		Link        string `json:"link"`
		Type        string `json:"type"`
		AddDate     string `json:"addDate"`
		Tips1       string `json:"tips1"`
		Tips2       string `json:"tips2"`
		BonusInfo   []struct {
			BonusType int `json:"bonusType"`
			Bonus     int `json:"bonus"`
		} `json:"bonusInfo"`
		TaskType  int           `json:"taskType"`
		Category  int           `json:"category"`
		Comic     interface{}   `json:"comic"`
		ComicList []interface{} `json:"comicList"`
		Exp       int           `json:"exp"`
	} `json:"data"`
}

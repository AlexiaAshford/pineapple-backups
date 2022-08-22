package book_info

type BoluobaoBookInfo struct {
	Uptime     string  `json:"lastUpdateTime"`
	MarkCount  int     `json:"markCount"`
	NovelCover string  `json:"novelCover"`
	BgBanner   string  `json:"bgBanner"`
	NovelID    int     `json:"novelId"`
	NovelName  string  `json:"novelName"`
	Point      float64 `json:"point"`
	IsFinish   bool    `json:"isFinish"`
	AuthorName string  `json:"authorName"`
	CharCount  int     `json:"charCount"`
	SignStatus string  `json:"signStatus"`
}

type BoluobaoBookSrc struct {
	Status struct {
		HTTPCode  int `json:"httpCode"`
		ErrorCode int `json:"errorCode"`
		MsgType   int `json:"msgType"`
		Msg       any `json:"msg"`
	} `json:"status"`
	Data BoluobaoBookInfo `json:"data"`
}

type CatBookSrc struct {
	Code string `json:"code"`
	Data struct {
		BookInfo CatBookInfo `json:"book_info"`
	} `json:"data"`
	Tip any `json:"tip"`
}

type CatBookInfo struct {
	NovelID         string `json:"book_id"`
	NovelName       string `json:"book_name"`
	Description     string `json:"description"`
	Tag             string `json:"tag"`
	TotalWordCount  string `json:"total_word_count"`
	UpStatus        string `json:"up_status"`
	UpdateStatus    string `json:"update_status"`
	IsPaid          string `json:"is_paid"`
	NovelCover      string `json:"cover"`
	AuthorName      string `json:"author_name"`
	Uptime          string `json:"uptime"`
	LastChapterInfo struct {
		ChapterID         string `json:"chapter_id"`
		BookID            string `json:"book_id"`
		ChapterIndex      string `json:"chapter_index"`
		ChapterTitle      string `json:"chapter_title"`
		Uptime            string `json:"uptime"`
		Mtime             string `json:"mtime"`
		RecommendBookInfo string `json:"recommend_book_info"`
	} `json:"last_chapter_info"`
	TagList struct {
		TagID   string `json:"tag_id"`
		TagType string `json:"tag_type"`
		TagName string `json:"tag_name"`
	} `json:"tag_list"`
	BookType        string `json:"book_type"`
	TransverseCover string `json:"transverse_cover"`
}

package hbooker_structs

type DetailStruct struct {
	Code         string        `json:"code"`
	Tip          any           `json:"tip"`
	Data         DetailData    `json:"data"`
	ScrollChests []ScrollChest `json:"scroll_chests"`
}

type DetailData struct {
	BookInfo    BookInfo   `json:"book_info"`
	RelatedList []BookInfo `json:"related_list"`
}

type BookInfo struct {
	BookID          string          `json:"book_id"`
	BookName        string          `json:"book_name"`
	Description     string          `json:"description"`
	BookSrc         string          `json:"book_src"`
	CategoryIndex   string          `json:"category_index"`
	Tag             string          `json:"tag"`
	TotalWordCount  string          `json:"total_word_count"`
	UpStatus        string          `json:"up_status"`
	UpdateStatus    string          `json:"update_status"`
	IsPaid          string          `json:"is_paid"`
	Cover           string          `json:"cover"`
	AuthorName      string          `json:"author_name"`
	Uptime          string          `json:"uptime"`
	Newtime         string          `json:"newtime"`
	ReviewAmount    string          `json:"review_amount"`
	RewardAmount    string          `json:"reward_amount"`
	ChapterAmount   string          `json:"chapter_amount"`
	LastChapterInfo LastChapterInfo `json:"last_chapter_info"`
	TagList         []TagList       `json:"tag_list"`
	BookType        string          `json:"book_type"`
	TransverseCover string          `json:"transverse_cover"`
}

type LastChapterInfo struct {
	ChapterID         string `json:"chapter_id"`
	BookID            string `json:"book_id"`
	ChapterIndex      string `json:"chapter_index"`
	ChapterTitle      string `json:"chapter_title"`
	Uptime            string `json:"uptime"`
	Mtime             string `json:"mtime"`
	RecommendBookInfo string `json:"recommend_book_info"`
}

type TagList struct {
	TagID   string `json:"tag_id"`
	TagType string `json:"tag_type"`
	TagName string `json:"tag_name"`
}

type ScrollChest struct {
	ChestID     string `json:"chest_id"`
	ReaderName  string `json:"reader_name"`
	Gender      string `json:"gender"`
	AvatarURL   string `json:"avatar_url"`
	BookName    string `json:"book_name"`
	Cost        int64  `json:"cost"`
	ChestImgURL string `json:"chest_img_url"`
	PropID      int64  `json:"prop_id"`
	Content     string `json:"content"`
}

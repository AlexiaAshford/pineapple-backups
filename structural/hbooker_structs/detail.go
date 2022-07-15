package hbooker_structs

type DetailStruct struct {
	Code         string        `json:"code"`
	Data         DetailData    `json:"data"`
	ScrollChests []ScrollChest `json:"scroll_chests"`
}

type DetailData struct {
	BookInfo                 BookInfo                   `json:"book_info"`
	IsInshelf                string                     `json:"is_inshelf"`
	IsBuy                    string                     `json:"is_buy"`
	UpReaderInfo             UpReaderInfo               `json:"up_reader_info"`
	RelatedList              []BookInfo                 `json:"related_list"`
	BookShortageReommendList []BookShortageReommendList `json:"book_shortage_reommend_list"`
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
	Discount        string          `json:"discount"`
	DiscountEndTime string          `json:"discount_end_time"`
	Cover           string          `json:"cover"`
	AuthorName      string          `json:"author_name"`
	Uptime          string          `json:"uptime"`
	Newtime         string          `json:"newtime"`
	ReviewAmount    string          `json:"review_amount"`
	RewardAmount    string          `json:"reward_amount"`
	ChapterAmount   string          `json:"chapter_amount"`
	IsOriginal      string          `json:"is_original"`
	TotalClick      string          `json:"total_click"`
	MonthClick      string          `json:"month_click"`
	WeekClick       string          `json:"week_click"`
	MonthNoVipClick string          `json:"month_no_vip_click"`
	WeekNoVipClick  string          `json:"week_no_vip_click"`
	TotalRecommend  string          `json:"total_recommend"`
	MonthRecommend  string          `json:"month_recommend"`
	WeekRecommend   string          `json:"week_recommend"`
	TotalFavor      string          `json:"total_favor"`
	MonthFavor      string          `json:"month_favor"`
	WeekFavor       string          `json:"week_favor"`
	CurrentYp       string          `json:"current_yp"`
	TotalYp         string          `json:"total_yp"`
	CurrentBlade    string          `json:"current_blade"`
	TotalBlade      string          `json:"total_blade"`
	WeekFansValue   string          `json:"week_fans_value"`
	MonthFansValue  string          `json:"month_fans_value"`
	TotalFansValue  string          `json:"total_fans_value"`
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

type BookShortageReommendList struct {
	BookID          string `json:"book_id"`
	BookName        string `json:"book_name"`
	Cover           string `json:"cover"`
	Discount        string `json:"discount"`
	DiscountEndTime string `json:"discount_end_time"`
	Introduce       string `json:"introduce"`
	Tag             string `json:"tag"`
	TotalFavor      string `json:"total_favor"`
}

type UpReaderInfo struct {
	ReaderID       string        `json:"reader_id"`
	Account        string        `json:"account"`
	ReaderName     string        `json:"reader_name"`
	AvatarURL      string        `json:"avatar_url"`
	AvatarThumbURL string        `json:"avatar_thumb_url"`
	BaseStatus     string        `json:"base_status"`
	ExpLV          string        `json:"exp_lv"`
	ExpValue       string        `json:"exp_value"`
	Gender         string        `json:"gender"`
	VipLV          string        `json:"vip_lv"`
	VipValue       string        `json:"vip_value"`
	IsAuthor       string        `json:"is_author"`
	IsUploader     string        `json:"is_uploader"`
	IsFollowing    string        `json:"is_following"`
	UsedDecoration []interface{} `json:"used_decoration"`
	IsInBlacklist  string        `json:"is_in_blacklist"`
	Ctime          string        `json:"ctime"`
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

package bookshelf

var GetShelfList = struct {
	Code string `json:"code"`
	Tip  any    `json:"tip"`
	Data struct {
		ShelfList []struct {
			ShelfID    string `json:"shelf_id"`
			ReaderID   string `json:"reader_id"`
			ShelfName  string `json:"shelf_name"`
			ShelfIndex string `json:"shelf_index"`
			BookLimit  string `json:"book_limit"`
		} `json:"shelf_list"`
	} `json:"data"`
	ScrollChests []interface{} `json:"scroll_chests"`
}{}

type GetShelfBookList struct {
	Code string `json:"code"`
	Tip  any    `json:"tip"`
	Data struct {
		BookList []struct {
			IsBuy    string `json:"is_buy"`
			BookInfo struct {
				BookID          string `json:"book_id"`
				BookName        string `json:"book_name"`
				IsOriginal      string `json:"is_original"`
				CategoryIndex   string `json:"category_index"`
				TotalWordCount  string `json:"total_word_count"`
				ReviewAmount    string `json:"review_amount"`
				Cover           string `json:"cover"`
				Discount        string `json:"discount"`
				DiscountEndTime string `json:"discount_end_time"`
				GloryTag        struct {
					TagName    string `json:"tag_name"`
					CornerName string `json:"corner_name"`
					LabelIcon  string `json:"label_icon"`
					LinkURL    string `json:"link_url"`
				} `json:"glory_tag"`
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
				UpdateStatus string `json:"update_status"`
				BookType     string `json:"book_type"`
			} `json:"book_info"`
			ModTime                   string `json:"mod_time"`
			LastReadChapterID         string `json:"last_read_chapter_id"`
			LastReadChapterUpdateTime string `json:"last_read_chapter_update_time"`
			IsNotify                  string `json:"is_notify"`
		} `json:"book_list"`
	} `json:"data"`
	ScrollChests []interface{} `json:"scroll_chests"`
}

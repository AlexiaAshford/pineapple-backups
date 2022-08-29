package bookshelf

type GetShelfList struct {
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
}

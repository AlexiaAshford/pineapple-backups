package hbooker_structs

type RecommendStruct struct {
	Code string `json:"code"`
	Data Data   `json:"data"`
}
type ChangeRecommendStruct struct {
	Code int        `json:"code"`
	Tip  string     `json:"tip"`
	Data ChangeData `json:"data"`
}

type ChangeBookList struct {
	BookID          string `json:"book_id"`
	BookName        string `json:"book_name"`
	CategoryIndex   string `json:"category_index"`
	Description     string `json:"description"`
	AuthorName      string `json:"author_name"`
	Cover           string `json:"cover"`
	DiscountEndTime string `json:"discount_end_time"`
	UpStatus        string `json:"up_status"`
	TotalWordCount  string `json:"total_word_count"`
	IsOriginal      int    `json:"is_original"`
	Introduce       string `json:"introduce"`
}

type ChangeData struct {
	BookList []ChangeBookList `json:"book_list"`
}
type CarouselList struct {
	PageURL       string `json:"page_url,omitempty"`
	CarouselCover string `json:"carousel_cover"`
	Rank          string `json:"rank"`
	CarouselType  string `json:"carousel_type"`
	BookID        string `json:"book_id,omitempty"`
}
type Header struct {
	Title         string `json:"title"`
	Image         string `json:"image"`
	CiweicatImage string `json:"ciweicat_image,omitempty"`
	InnerApp      string `json:"inner_app"`
	Version       string `json:"version"`
	CategoryIndex int    `json:"category_index,omitempty"`
}
type RankList struct {
	Image         string   `json:"image"`
	CiweicatImage string   `json:"ciweicat_image"`
	InnerApp      string   `json:"inner_app"`
	Order         string   `json:"order"`
	Title         string   `json:"title"`
	SubTitle      string   `json:"sub_title"`
	BookInfo      BookInfo `json:"book_info"`
}
type DesBookList struct {
	BookID          string `json:"book_id"`
	BookName        string `json:"book_name"`
	CategoryIndex   string `json:"category_index"`
	Description     string `json:"description"`
	AuthorName      string `json:"author_name"`
	Cover           string `json:"cover"`
	Discount        int    `json:"discount"`
	DiscountEndTime string `json:"discount_end_time"`
	TotalClick      int    `json:"total_click"`
	IsOriginal      int    `json:"is_original"`
	Introduce       string `json:"introduce"`
}
type BossModule struct {
	ModuleImage         string        `json:"module_image"`
	CiweicatModuleImage string        `json:"ciweicat_module_image"`
	ListID              string        `json:"list_id"`
	ListURL             string        `json:"list_url"`
	DesBookList         []DesBookList `json:"des_book_list"`
	RightButtonType     string        `json:"right_button_type"`
}
type EditorModule struct {
	ModuleImage         string        `json:"module_image"`
	CiweicatModuleImage string        `json:"ciweicat_module_image"`
	ModuleTitle         string        `json:"module_title"`
	ListID              string        `json:"list_id"`
	ListURL             string        `json:"list_url"`
	DesBookList         []DesBookList `json:"des_book_list"`
	RightButtonType     string        `json:"right_button_type"`
}
type SingleBooklist struct {
	ListID    string `json:"list_id"`
	ListURL   string `json:"list_url"`
	ListCover string `json:"list_cover"`
}
type PicBookList struct {
	BookID          string `json:"book_id"`
	BookName        string `json:"book_name"`
	Cover           string `json:"cover"`
	Discount        int    `json:"discount"`
	DiscountEndTime string `json:"discount_end_time"`
	IsOriginal      int    `json:"is_original"`
	CategoryIndex   string `json:"category_index"`
	Description     string `json:"description"`
	TotalFavor      string `json:"total_favor"`
	TotalClick      int    `json:"total_click"`
	UpReaderID      string `json:"up_reader_id"`
	AuthorName      string `json:"author_name"`
}
type MoreBooklist struct {
	ListID    string `json:"list_id"`
	ListCover string `json:"list_cover"`
}
type ModuleList struct {
	ModuleType          int            `json:"module_type"`
	Header              []Header       `json:"header,omitempty"`
	ModuleID            int            `json:"module_id"`
	RankList            []RankList     `json:"rank_list,omitempty"`
	BossModule          BossModule     `json:"boss_module,omitempty"`
	ListID              string         `json:"list_id,omitempty"`
	ModuleTitle         string         `json:"module_title,omitempty"`
	ModuleIntroduce     string         `json:"module_introduce,omitempty"`
	ModuleImage         string         `json:"module_image,omitempty"`
	CiweicatModuleImage string         `json:"ciweicat_module_image,omitempty"`
	EditorModule        EditorModule   `json:"editor_module,omitempty"`
	SingleBooklist      SingleBooklist `json:"single_booklist,omitempty"`
	RightButtonType     string         `json:"right_button_type,omitempty"`
	InternalModuleTitle string         `json:"internal_module_title,omitempty"`
	PicBookList         []PicBookList  `json:"pic_book_list,omitempty"`
	ListURL             string         `json:"list_url,omitempty"`
	MoreBooklist        []MoreBooklist `json:"more_booklist,omitempty"`
}
type NewBookList struct {
	BookID          string `json:"book_id"`
	Cover           string `json:"cover"`
	Discount        int    `json:"discount"`
	DiscountEndTime string `json:"discount_end_time"`
	BookName        string `json:"book_name"`
	IsOriginal      int    `json:"is_original"`
}
type UpBookList struct {
	BookID          string `json:"book_id"`
	Cover           string `json:"cover"`
	Discount        int    `json:"discount"`
	DiscountEndTime string `json:"discount_end_time"`
	BookName        string `json:"book_name"`
	IsOriginal      int    `json:"is_original"`
}

type CiweicatNewBookList struct {
	BookList            []BookList `json:"book_list"`
	CiweicatModuleImage string     `json:"ciweicat_module_image"`
	RightButtonType     string     `json:"right_button_type"`
	ModuleIntroduce     string     `json:"module_introduce"`
	ModuleTitle         string     `json:"module_title"`
}
type CiweicatUpBookList struct {
	BookList            []BookList `json:"book_list"`
	CiweicatModuleImage string     `json:"ciweicat_module_image"`
	RightButtonType     string     `json:"right_button_type"`
	ModuleIntroduce     string     `json:"module_introduce"`
	ModuleTitle         string     `json:"module_title"`
}
type Data struct {
	CarouselList        []CarouselList      `json:"carousel_list"`
	ModuleList          []ModuleList        `json:"module_list"`
	NewBookList         []NewBookList       `json:"new_book_list"`
	UpBookList          []UpBookList        `json:"up_book_list"`
	CiweicatNewBookList CiweicatNewBookList `json:"ciweicat_new_book_list"`
	CiweicatUpBookList  CiweicatUpBookList  `json:"ciweicat_up_book_list"`
	NewRecommend        int                 `json:"new_recommend"`
}

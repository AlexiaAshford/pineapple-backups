package https

const (
	CatalogueDetailedByDivisionId = "chapter/get_updated_chapter_by_division_id"
	CatDivisionIdByBookId         = "book/get_division_list"
	CatBookDetailedById           = "book/get_info_by_id"
	CatSearchDetailed             = "bookcity/get_filter_search_book_list"
	CatChapterKeyByCid            = "chapter/get_chapter_cmd"
	CatContentDetailedByCid       = "chapter/get_cpt_ifm"
	CatRecommend                  = "bookcity/get_index_list"
	CatChangeRecommend            = "bookcity/change_recommend_exposure_books"
	CatWebSite                    = "https://app.hbooker.com/"
	CatUseGeetestSignup           = "signup/use_geetest"
	CatGeetestFirstRegister       = "signup/geetest_first_register"
	CatLoginByAccount             = "signup/login"
)

const (
	SFBookDetailedById      = "novels/%v?expand="
	SFWebSite               = "https://minipapi.sfacg.com/pas/mpapi/"
	SFAccountDetailedByApi  = "user"
	SFCatalogueDetailedById = "novels/%v/dirs?expand=originNeedFireMoney"
	SFContentDetailedByCid  = "Chaps/%v?expand=content&autoOrder=true"
	SFSearch                = "search/novels/result?q=%v&size=20&page=%v&expand="
	SFLogin                 = "sessions"
	//SFBookShelfDetailed       = "novels/%v/shelf?expand="
)

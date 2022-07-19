package HbookerAPI

const (
	CatalogueDetailedByDivisionId = "chapter/get_updated_chapter_by_division_id?division_id="
	DivisionIdByBookId            = "book/get_division_list?book_id="
	LoginByAccount                = "signup/login?login_name=%v&password=%v"
	BookDetailedById              = "book/get_info_by_id?book_id=%v"
	SearchDetailedByKeyword       = "bookcity/get_filter_search_book_list?count=10&page=%v&category_index=0&key=%v"
	//AccountDetailedByApi          = "user"
	//ContentDetailedByCid          = "Chaps/%v?expand=content&autoOrder=true"
)

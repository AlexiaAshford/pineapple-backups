package HbookerAPI

const (
	BookDetailedById              = "novels/%v?expand="
	AccountDetailedByApi          = "user"
	CatalogueDetailedByDivisionId = "/chapter/get_updated_chapter_by_division_id?division_id="
	DivisionIdByBookId            = "/book/get_division_list?book_id="
	ContentDetailedByCid          = "Chaps/%v?expand=content&autoOrder=true"
	SearchDetailedByKeyword       = "search/novels/result?q=%v&size=20&page=0&expand="
	LoginByAccount                = "sessions"
)

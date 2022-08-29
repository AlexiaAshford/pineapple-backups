package hbooker

import (
	"fmt"
	"github.com/gookit/color"
	"sf/cfg"
	req "sf/src/https"
	_struct "sf/struct"
	structs "sf/struct/hbooker_structs"
	"strconv"
	"time"
)

func GET_DIVISION(BookId string) []map[string]string {
	var division_info []map[string]string
	var chapter_index int
	response := req.Get("book/get_division_list", &structs.DivisionStruct{}, map[string]string{"book_id": BookId})
	for division_index, division := range response.(*structs.DivisionStruct).Data.DivisionList {
		fmt.Printf("第%v卷\t\t%v\n", division_index+1, division.DivisionName)
		for _, chapter := range GET_CATALOGUE(division.DivisionID) {
			chapter_index += 1
			division_info = append(division_info, map[string]string{
				"division_name":  division.DivisionName,
				"division_id":    division.DivisionID,
				"division_index": strconv.Itoa(division_index),
				"chapter_name":   chapter.ChapterTitle,
				"chapter_id":     chapter.ChapterID,
				"is_valid":       chapter.IsValid,
				"money":          chapter.AuthAccess,
				"chapter_index":  strconv.Itoa(chapter_index),
				"file_name":      cfg.Config_file_name(division_index, chapter_index, chapter.ChapterID),
			})
		}
	}
	return division_info
}

func GET_CATALOGUE(DivisionId string) []structs.ChapterList {
	params := map[string]string{"division_id": DivisionId}
	return req.Get("chapter/get_updated_chapter_by_division_id", &structs.ChapterStruct{}, params).(*structs.ChapterStruct).Data.ChapterList
}

func GET_BOOK_SHELF_INDEXES_INFORMATION(shelf_id, last_mod_time, direction string) (map[string]string, error) {
	params := map[string]string{"shelf_id": shelf_id, "last_mod_time": last_mod_time, "direction": direction}
	response := req.Get("bookshelf/get_shelf_book_list", &structs.DetailStruct{}, params).(*structs.DetailStruct)
	fmt.Println(response)
	return nil, nil
}

func GET_BOOK_SHELF_INFORMATION() (map[int][]map[string]string, error) {
	params, bookshelf_info := map[string]string{"expand": "novels"}, make(map[int][]map[string]string)
	fmt.Println(params)
	response := req.Get("bookshelf/get_shelf_list", &structs.DetailStruct{}, nil).(*structs.DetailStruct)
	fmt.Println(response)
	//if response.Status.HTTPCode != 200 {
	//	return nil, fmt.Errorf(response.Status.Msg.(string))
	//}
	//for index, value := range response.Data {
	//	fmt.Println("bookshelf index:", index, "\t\t\tbookshelf name:", value.Name)
	//	var bookshelf_info_list []map[string]string
	//	for _, book := range value.Expand.Novels {
	//		bookshelf_info_list = append(bookshelf_info_list,
	//			map[string]string{"novel_name": book.NovelName, "novel_id": strconv.Itoa(book.NovelID)},
	//		)
	//	}
	//	bookshelf_info[index] = bookshelf_info_list
	return bookshelf_info, nil
}
func GET_BOOK_INFORMATION(bid string) (_struct.Books, error) {
	params := map[string]string{"book_id": bid}
	response := req.Get("book/get_info_by_id", &structs.DetailStruct{}, params).(*structs.DetailStruct)
	if response.Code == "100000" {
		return _struct.Books{
			NovelName:  cfg.RegexpName(response.Data.BookInfo.BookName),
			NovelID:    response.Data.BookInfo.BookID,
			NovelCover: response.Data.BookInfo.Cover,
			AuthorName: response.Data.BookInfo.AuthorName,
			CharCount:  response.Data.BookInfo.TotalWordCount,
			MarkCount:  response.Data.BookInfo.UpdateStatus,
			SignStatus: response.Data.BookInfo.IsPaid,
		}, nil
	} else {
		return _struct.Books{}, fmt.Errorf(response.Tip.(string))
	}
}

func GET_SEARCH(KeyWord string, page int) *structs.SearchStruct {
	params := map[string]string{"count": "10", "page": strconv.Itoa(page), "category_index": "0", "key": KeyWord}
	return req.Get("bookcity/get_filter_search_book_list", &structs.SearchStruct{}, params).(*structs.SearchStruct)
}

func Login(account, password string) {
	params := map[string]string{"login_name": account, "password": password}
	response := req.Get("signup/login", &structs.LoginStruct{}, params)
	result := response.(*structs.LoginStruct)
	if result.Code == "100000" {
		cfg.Apps.Cat.Params.LoginToken = result.Data.LoginToken
		cfg.Apps.Cat.Params.Account = result.Data.ReaderInfo.Account
		cfg.SaveJson()
	} else {
		fmt.Println("Login failed!")
	}
}

func UseGeetest() *structs.GeetestStruct {
	return req.Get("signup/use_geetest", &structs.GeetestStruct{}, nil).(*structs.GeetestStruct)
}

func GeetestRegister(userID string) (string, string) {
	params := map[string]string{"t": strconv.FormatInt(time.Now().UnixNano()/1e6, 10), "user_id": userID}
	response, _ := req.Request("POST", req.QueryParams("signup/geetest_first_register", params))
	result := req.JsonUnmarshal(response, &structs.GeetestChallenge{}).(*structs.GeetestChallenge)
	return result.Challenge, result.Gt
}
func TestGeetest(userID string) {
	UseGeetest()
	challenge, gt := GeetestRegister(userID)
	status, CaptchaType, errorDetail := GetFullBG(&Geetest{GT: gt, Challenge: challenge})
	fmt.Println(status, CaptchaType, errorDetail)
	if status == "success" {
		color.Infoln("验证码类型：", CaptchaType, "")
	} else {
		color.Errorln("获取图片失败 Error: ", status, " ErrorDetail:", errorDetail)
		TestGeetest(userID)
	}
}

func GetRecommend() *structs.RecommendStruct {
	params := map[string]string{"theme_type": "NORMAL", "tab_type": "200"}
	return req.Get("bookcity/get_index_list", &structs.RecommendStruct{}, params).(*structs.RecommendStruct)
}
func GetChangeRecommend() []structs.ChangeBookList {
	params := map[string]string{"book_id": "100250589,100283902,100186621,100287528,100309123,100325245", "from_module_name": "长篇好书"}
	return req.Get("bookcity/change_recommend_exposure_books", &structs.ChangeRecommendStruct{}, params).(*structs.ChangeRecommendStruct).Data.BookList
}

func GetKeyByCid(chapterId string) string {
	params := map[string]string{"chapter_id": chapterId}
	return req.Get("chapter/get_chapter_cmd", &structs.KeyStruct{}, params).(*structs.KeyStruct).Data.Command
}

func GetContent(chapterId string, ChapterKey string) *structs.ContentStruct {
	params := map[string]string{"chapter_id": chapterId, "chapter_command": ChapterKey}
	return req.Get("chapter/get_cpt_ifm", &structs.ContentStruct{}, params).(*structs.ContentStruct)
}

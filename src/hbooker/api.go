package hbooker

import (
	"fmt"
	"github.com/gookit/color"
	"sf/cfg"
	"sf/src/hbooker/Encrypt"
	req "sf/src/https"
	_struct "sf/struct"
	structs "sf/struct/hbooker_structs"
	"sf/struct/hbooker_structs/bookshelf"
	"sf/struct/hbooker_structs/division"
	"strconv"
	"time"
)

func GET_DIVISION(BookId string) []map[string]string {
	var (
		division_info_list []map[string]string
		chapter_index      int
	)
	response := req.Get(GET_DIVISION_LIST, &division.DivisionList{}, map[string]string{"book_id": BookId})
	for division_index, division_info := range response.(*division.DivisionList).Data.DivisionList {
		fmt.Printf("第%v卷\t\t%v\n", division_index+1, division_info.DivisionName)
		for _, chapter := range GET_CATALOGUE(division_info.DivisionID) {
			chapter_index += 1
			division_info_list = append(division_info_list, map[string]string{
				"is_valid":       chapter.IsValid,
				"chapter_id":     chapter.ChapterID,
				"money":          chapter.AuthAccess,
				"chapter_name":   chapter.ChapterTitle,
				"division_name":  division_info.DivisionName,
				"division_id":    division_info.DivisionID,
				"division_index": strconv.Itoa(division_index),
				"chapter_index":  strconv.Itoa(chapter_index),
				"file_name":      cfg.Config_file_name(division_index, chapter_index, chapter.ChapterID),
			})
		}
	}
	return division_info_list
}

func GET_CATALOGUE(DivisionId string) []structs.ChapterList {
	params := map[string]string{"division_id": DivisionId}
	return req.Get(GET_CHAPTER_UPDATE, &structs.ChapterStruct{}, params).(*structs.ChapterStruct).Data.ChapterList
}

func GET_BOOK_SHELF_INDEXES_INFORMATION(shelf_id string) ([]map[string]string, error) {
	params := map[string]string{"shelf_id": shelf_id, "last_mod_time": "0", "direction": "prev"}
	response := req.Get(BOOKSHELF_GET_SHELF_BOOK_LIST, &bookshelf.GetShelfBookList{}, params).(*bookshelf.GetShelfBookList)
	if response.Code != "100000" {
		return nil, fmt.Errorf(response.Tip.(string))
	}
	var bookshelf_info_list []map[string]string
	for _, book := range response.Data.BookList {
		bookshelf_info_list = append(bookshelf_info_list,
			map[string]string{"novel_name": book.BookInfo.BookName, "novel_id": book.BookInfo.BookID},
		)
	}

	return bookshelf_info_list, nil
}

func GET_BOOK_SHELF_INFORMATION() (map[int][]map[string]string, error) {
	bookshelf_info := make(map[int][]map[string]string)
	response := req.Get(BOOKSHELF_GET_SHELF_LIST, &bookshelf.GetShelfList{}, nil).(*bookshelf.GetShelfList)
	if response.Code != "100000" {
		return nil, fmt.Errorf(response.Tip.(string))
	}
	for index, value := range response.Data.ShelfList {
		fmt.Println("bookshelf index:", index, "\t\t\tbookshelf name:", value.ShelfName)
		if bookshelf_info_list, err := GET_BOOK_SHELF_INDEXES_INFORMATION(value.ShelfID); err == nil {
			bookshelf_info[index] = bookshelf_info_list
		} else {
			fmt.Println("ShelfID:", value.ShelfID, "\terr:", err)
		}
	}
	return bookshelf_info, nil
}
func GET_BOOK_INFORMATION(bid string) (_struct.Books, error) {
	params := map[string]string{"book_id": bid}
	response := req.Get(BOOK_GET_INFO_BY_ID, &structs.DetailStruct{}, params).(*structs.DetailStruct)
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
	return req.Get(BOOKCITY_GET_FILTER_LIST, &structs.SearchStruct{}, params).(*structs.SearchStruct)
}

func Login(account, password string) {
	params := map[string]string{"login_name": account, "password": password}
	response := req.Get(MY_SIGN_LOGIN, &structs.LoginStruct{}, params)
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
	return req.Get(BOOKCITY_RECOMMEND_DATA, &structs.RecommendStruct{}, params).(*structs.RecommendStruct)
}
func GetChangeRecommend() []structs.ChangeBookList {
	params := map[string]string{"book_id": "100250589,100283902,100186621,100287528,100309123,100325245", "from_module_name": "长篇好书"}
	return req.Get(GET_CHANGE_RECOMMEND, &structs.ChangeRecommendStruct{}, params).(*structs.ChangeRecommendStruct).Data.BookList
}

func GetKeyByCid(chapterId string) string {
	params := map[string]string{"chapter_id": chapterId}
	return req.Get(GET_CHAPTER_KEY, &structs.KeyStruct{}, params).(*structs.KeyStruct).Data.Command
}

func GET_CHAPTER_CONTENT(chapterId, chapter_key string) string {
	params := map[string]string{"chapter_id": chapterId, "chapter_command": chapter_key}
	response := req.Get(GET_CPT_IFM, &structs.ContentStruct{}, params).(*structs.ContentStruct)
	if response != nil && response.Code == "100000" {
		chapter_info := response.Data.ChapterInfo
		content := string(Encrypt.Decode(chapter_info.TxtContent, chapter_key))
		content_title := fmt.Sprintf("%v: %v", chapter_info.ChapterTitle, chapter_info.Uptime)
		return content_title + "\n\n" + cfg.StandardContent(content)
	} else {
		fmt.Println("download failed! chapterId:", chapterId, "error:", response.Tip)
	}
	return ""
}

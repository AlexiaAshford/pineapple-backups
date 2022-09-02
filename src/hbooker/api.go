package hbooker

import (
	"fmt"
	"github.com/VeronicaAlexia/pineapple-backups/config"
	"github.com/VeronicaAlexia/pineapple-backups/src/hbooker/Encrypt"
	req "github.com/VeronicaAlexia/pineapple-backups/src/https"
	_struct "github.com/VeronicaAlexia/pineapple-backups/struct"
	structs "github.com/VeronicaAlexia/pineapple-backups/struct/hbooker_structs"
	"github.com/VeronicaAlexia/pineapple-backups/struct/hbooker_structs/bookshelf"
	"github.com/VeronicaAlexia/pineapple-backups/struct/hbooker_structs/division"
	"github.com/gookit/color"
	"strconv"
	"time"
)

func GET_DIVISION(BookId string) []map[string]string {
	var chapter_index int
	var division_info_list []map[string]string
	var s = new(division.DivisionList)
	req.Get(new(req.Context).Init(GET_DIVISION_LIST).Query("book_id", BookId).QueryToString(), s)
	for division_index, division_info := range s.Data.DivisionList {
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
				"file_name":      config.FileCacheName(division_index, chapter_index, chapter.ChapterID),
			})
		}
	}
	return division_info_list
}

func GET_CATALOGUE(DivisionId string) []structs.ChapterList {
	s := new(structs.ChapterStruct)
	req.Get(new(req.Context).Init(GET_CHAPTER_UPDATE).Query("division_id", DivisionId).QueryToString(), s)
	return s.Data.ChapterList
}

func GET_BOOK_SHELF_INDEXES_INFORMATION(shelf_id string) ([]map[string]string, error) {
	s := new(bookshelf.GetShelfBookList)
	req.Get(new(req.Context).Init(BOOKSHELF_GET_SHELF_BOOK_LIST).Query("shelf_id", shelf_id).
		Query("direction", "prev").Query("last_mod_time", "0").QueryToString(), s)
	if s.Code != "100000" {
		return nil, fmt.Errorf(s.Tip.(string))
	}
	var bookshelf_info_list []map[string]string
	for _, book := range s.Data.BookList {
		bookshelf_info_list = append(bookshelf_info_list,
			map[string]string{"novel_name": book.BookInfo.BookName, "novel_id": book.BookInfo.BookID},
		)
	}

	return bookshelf_info_list, nil
}

func GET_BOOK_SHELF_INFORMATION() (map[int][]map[string]string, error) {
	s := new(bookshelf.GetShelfList)
	bookshelf_info := make(map[int][]map[string]string)
	req.Get(new(req.Context).Init(BOOKSHELF_GET_SHELF_LIST).QueryToString(), s)
	if s.Code != "100000" {
		return nil, fmt.Errorf(s.Tip.(string))
	}
	for index, value := range s.Data.ShelfList {
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
	s := new(structs.DetailStruct)
	req.Get(new(req.Context).Init(BOOK_GET_INFO_BY_ID).Query("book_id", bid).QueryToString(), s)
	if s.Code == "100000" {
		return _struct.Books{
			NovelName:  config.RegexpName(s.Data.BookInfo.BookName),
			NovelID:    s.Data.BookInfo.BookID,
			NovelCover: s.Data.BookInfo.Cover,
			AuthorName: s.Data.BookInfo.AuthorName,
			CharCount:  s.Data.BookInfo.TotalWordCount,
			MarkCount:  s.Data.BookInfo.UpdateStatus,
			SignStatus: s.Data.BookInfo.IsPaid,
		}, nil
	} else {
		return _struct.Books{}, fmt.Errorf(s.Tip.(string))
	}
}

func GET_SEARCH(KeyWord string, page int) *structs.SearchStruct {
	s := new(structs.SearchStruct)
	req.Get(new(req.Context).Init(BOOKCITY_GET_FILTER_LIST).Query("count", "10").
		Query("page", strconv.Itoa(page)).Query("category_index", "0").Query("key", KeyWord).
		QueryToString(), s)
	return s
}

func Login(account, password string) {
	s := new(structs.LoginStruct)
	req.Get(new(req.Context).Init(MY_SIGN_LOGIN).Query("login_name", account).
		Query("password", password).QueryToString(), s)
	if s.Code == "100000" {
		config.Apps.Cat.Params.LoginToken = s.Data.LoginToken
		config.Apps.Cat.Params.Account = s.Data.ReaderInfo.Account
		config.SaveJson()
	} else {
		fmt.Println("Login failed!")
	}
}

func UseGeetest() *structs.GeetestStruct {
	return req.Get("signup/use_geetest", &structs.GeetestStruct{}).(*structs.GeetestStruct)
}

func GeetestRegister(userID string) (string, string) {
	s := new(structs.GeetestChallenge)
	response, _ := req.Request("POST", new(req.Context).Init("signup/geetest_register").
		Query("t", strconv.FormatInt(time.Now().UnixNano()/1e6, 10)).
		Query("user_id", userID).QueryToString())
	req.JsonUnmarshal(response, s)
	return s.Challenge, s.Gt
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
	s := new(structs.RecommendStruct)
	req.Get(new(req.Context).Init(BOOKCITY_RECOMMEND_DATA).Query("theme_type", "NORMAL").
		Query("tab_type", "200").QueryToString(), s)
	return s
}

func GetChangeRecommend() []structs.ChangeBookList {
	s := new(structs.ChangeRecommendStruct)
	req.Get(new(req.Context).Init(GET_CHANGE_RECOMMEND).
		Query("book_id", "100250589,100283902,100186621,100287528,100309123,100325245").QueryToString(), s)
	return s.Data.BookList
}

func GetKeyByCid(chapterId string) string {
	s := new(structs.KeyStruct)
	req.Get(new(req.Context).Init(GET_CHAPTER_KEY).Query("chapter_id", chapterId).QueryToString(), s)
	return s.Data.Command
}

func GET_CHAPTER_CONTENT(chapterId, chapter_key string) string {
	s := new(structs.ContentStruct)
	req.Get(new(req.Context).Init(GET_CPT_IFM).Query("chapter_id", chapterId).
		Query("chapter_command", chapter_key).QueryToString(), s)
	if s != nil && s.Code == "100000" {
		chapter_info := s.Data.ChapterInfo
		content := string(Encrypt.Decode(chapter_info.TxtContent, chapter_key))
		content_title := fmt.Sprintf("%v: %v", chapter_info.ChapterTitle, chapter_info.Uptime)
		return content_title + "\n\n" + config.StandardContent(content)
	} else {
		fmt.Println("download failed! chapterId:", chapterId, "error:", s.Tip)
	}
	return ""
}

func mains() {
	Login("account", "password")
	TestGeetest("123")
	GetRecommend()
	GetChangeRecommend()
}

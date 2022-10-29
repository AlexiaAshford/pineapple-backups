package hbooker

import (
	"fmt"
	"github.com/VeronicaAlexia/pineapple-backups/config"
	config_file "github.com/VeronicaAlexia/pineapple-backups/config/file"
	"github.com/VeronicaAlexia/pineapple-backups/config/tool"
	encryption2 "github.com/VeronicaAlexia/pineapple-backups/pkg/encryption"
	"github.com/VeronicaAlexia/pineapple-backups/pkg/request"
	_struct "github.com/VeronicaAlexia/pineapple-backups/struct"
	structs "github.com/VeronicaAlexia/pineapple-backups/struct/hbooker_structs"
	"github.com/VeronicaAlexia/pineapple-backups/struct/hbooker_structs/bookshelf"
	"github.com/VeronicaAlexia/pineapple-backups/struct/hbooker_structs/division"
	"github.com/gookit/color"
	"strconv"
	"strings"
	"time"
)

func GET_DIVISION(BookId string) []map[string]string {
	var chapter_index int
	var division_info_list []map[string]string
	request.NewHttpUtils(GET_DIVISION_LIST, "POST").Add("book_id", BookId).NewRequests().Unmarshal(&division.VolumeList)
	for division_index, division_info := range division.VolumeList.Data.DivisionList {
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
				"file_name":      config_file.FileCacheName(division_index, chapter_index, chapter.ChapterID),
			})
		}
	}
	return division_info_list
}

func GET_CATALOGUE(DivisionId string) []structs.ChapterList {
	request.NewHttpUtils(GET_CHAPTER_UPDATE, "POST").Add("division_id", DivisionId).NewRequests().Unmarshal(&structs.Chapter)
	return structs.Chapter.Data.ChapterList
}

func GET_BOOK_SHELF_INDEXES_INFORMATION(shelf_id string) ([]map[string]string, error) {
	request.NewHttpUtils(BOOKSHELF_GET_SHELF_BOOK_LIST, "POST").Add("shelf_id", shelf_id).Add("direction", "prev").
		Add("last_mod_time", "0").NewRequests().Unmarshal(&bookshelf.BookList)
	if bookshelf.BookList.Code != "100000" {
		return nil, fmt.Errorf(bookshelf.BookList.Tip.(string))
	}
	var bookshelf_info_list []map[string]string
	for _, book := range bookshelf.BookList.Data.BookList {
		bookshelf_info_list = append(bookshelf_info_list,
			map[string]string{"novel_name": book.BookInfo.BookName, "novel_id": book.BookInfo.BookID},
		)
	}
	return bookshelf_info_list, nil
}

func GET_BOOK_SHELF_INFORMATION() (map[int][]map[string]string, error) {
	bookshelf_info := make(map[int][]map[string]string)
	request.NewHttpUtils(BOOKSHELF_GET_SHELF_LIST, "POST").NewRequests().Unmarshal(&bookshelf.GetShelfList)
	if bookshelf.GetShelfList.Code != "100000" {
		return nil, fmt.Errorf(bookshelf.GetShelfList.Tip.(string))
	}
	for index, value := range bookshelf.GetShelfList.Data.ShelfList {
		fmt.Println("书架号:", index, "\t\t\t书架名:", value.ShelfName)
		if bookshelf_info_list, err := GET_BOOK_SHELF_INDEXES_INFORMATION(value.ShelfID); err == nil {
			bookshelf_info[index] = bookshelf_info_list
		} else {
			fmt.Println("ShelfID:", value.ShelfID, "\terr:", err)
		}
	}
	return bookshelf_info, nil
}
func GET_BOOK_INFORMATION(bid string) error {
	request.NewHttpUtils(BOOK_GET_INFO_BY_ID, "POST").
		Add("book_id", bid).NewRequests().Unmarshal(&structs.Detail)
	if structs.Detail.Code == "100000" {
		config.Current.Book = _struct.Books{
			NovelName:  tool.RegexpName(structs.Detail.Data.BookInfo.BookName),
			NovelID:    structs.Detail.Data.BookInfo.BookID,
			NovelCover: structs.Detail.Data.BookInfo.Cover,
			AuthorName: structs.Detail.Data.BookInfo.AuthorName,
			CharCount:  structs.Detail.Data.BookInfo.TotalWordCount,
			MarkCount:  structs.Detail.Data.BookInfo.UpdateStatus,
			SignStatus: structs.Detail.Data.BookInfo.IsPaid,
		}
		return nil
	}
	return fmt.Errorf(structs.Detail.Tip.(string))
}

func GET_SEARCH(KeyWord string, page int) {
	request.NewHttpUtils(BOOKCITY_GET_FILTER_LIST, "POST").Add("count", "10").
		Add("page", strconv.Itoa(page)).Add("category_index", "0").Add("key", KeyWord).NewRequests().Unmarshal(&structs.Search)

}

func GET_LOGIN_TOKEN(account, password string) {
	// hbooker new version add GEETEST verification, if you enter the wrong information or log in multiple times, GEETEST verification will be triggered.
	// IP address may need to log in again after a few hours to avoid triggering verification, you can try to change the IP to avoid triggering verification.
	GET_USE_GEETEST()
	request.NewHttpUtils(MY_SIGN_LOGIN, "POST").Add("login_name", account).
		Add("password", password).NewRequests().Unmarshal(&structs.Login)
	if structs.Login.Code == "100000" {
		config.Apps.Cat.Params.LoginToken = structs.Login.Data.LoginToken
		config.Apps.Cat.Params.Account = structs.Login.Data.ReaderInfo.Account
		config.SaveJson()
	} else {
		fmt.Println("Login failed!", structs.Login.Tip)
	}
}

func GET_USE_GEETEST() {
	request.NewHttpUtils(USE_GEETEST, "POST").NewRequests().Unmarshal(&structs.Geetest)
}

func GET_GEETEST_REGISTER(UserID string) (string, string) {
	request.NewHttpUtils(GEETEST_REGISTER, "POST").Add("user_id", UserID).
		Add("t", strconv.FormatInt(time.Now().UnixNano()/1e6, 10)).NewRequests().Unmarshal(&structs.Challenge)
	return structs.Challenge.Challenge, structs.Challenge.Gt
}
func TEST_GEETEST(userID string) {
	GET_USE_GEETEST()
	challenge, gt := GET_GEETEST_REGISTER(userID)
	status, CaptchaType, errorDetail := encryption2.GetFullBG(&encryption2.Geetest{GT: gt, Challenge: challenge})
	fmt.Println(status, CaptchaType, errorDetail)
	if status == "success" {
		color.Infoln("验证码类型：", CaptchaType, "")
		encryption2.Slide(&encryption2.Geetest{GT: gt, Challenge: challenge})
	} else {
		color.Errorln("获取图片失败 Error: ", status, " ErrorDetail:", errorDetail)
		TEST_GEETEST(userID)
	}
}

func GET_KET_BY_CHAPTER_ID(chapterId string) string {
	request.NewHttpUtils(GET_CHAPTER_KEY, "POST").Add("chapter_id", chapterId).NewRequests().Unmarshal(&structs.Key)
	return structs.Key.Data.Command
}

func GET_CHAPTER_CONTENT(chapterId, chapter_key string) string {
	request.NewHttpUtils(GET_CPT_IFM, "POST").Add("chapter_id", chapterId).
		Add("chapter_command", chapter_key).NewRequests().Unmarshal(&structs.Content)
	if structs.Content.Code == "100000" {
		chapter_info := structs.Content.Data.ChapterInfo
		content := string(encryption2.Decode(chapter_info.TxtContent, chapter_key))
		content_title := fmt.Sprintf("%v: %v", chapter_info.ChapterTitle, chapter_info.Uptime)
		return content_title + "\n\n" + tool.StandardContent(strings.Split(content, "\n"))
	} else {
		fmt.Println("download failed! chapterId:", chapterId, "error:", structs.Content.Tip)
	}
	return ""
}

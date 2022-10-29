package boluobao

import (
	"fmt"
	"github.com/VeronicaAlexia/pineapple-backups/config"
	"github.com/VeronicaAlexia/pineapple-backups/pkg/file"
	"github.com/VeronicaAlexia/pineapple-backups/pkg/request"
	"github.com/VeronicaAlexia/pineapple-backups/pkg/tools"
	_struct "github.com/VeronicaAlexia/pineapple-backups/struct"
	"github.com/VeronicaAlexia/pineapple-backups/struct/sfacg_structs"
	"github.com/VeronicaAlexia/pineapple-backups/struct/sfacg_structs/bookshelf"
	url_ "net/url"
	"os"
	"strconv"
	"strings"
)

func GET_BOOK_INFORMATION(NovelId string) error {
	request.NewHttpUtils("novels/"+NovelId, "GET").Add("expand", "intro,tags,sysTags").NewRequests().Unmarshal(&sfacg_structs.BookInfo)
	if sfacg_structs.BookInfo.Status.HTTPCode == 200 && sfacg_structs.BookInfo.Data.NovelName != "" {
		config.Current.Book = _struct.Books{
			NovelCover: sfacg_structs.BookInfo.Data.NovelCover,
			AuthorName: sfacg_structs.BookInfo.Data.AuthorName,
			SignStatus: sfacg_structs.BookInfo.Data.SignStatus,
			NovelID:    strconv.Itoa(sfacg_structs.BookInfo.Data.NovelID),
			CharCount:  strconv.Itoa(sfacg_structs.BookInfo.Data.CharCount),
			MarkCount:  strconv.Itoa(sfacg_structs.BookInfo.Data.MarkCount),
			NovelName:  tools.RegexpName(sfacg_structs.BookInfo.Data.NovelName),
		}
		return nil
	} else {
		if sfacg_structs.BookInfo.Status.Msg != nil {
			return fmt.Errorf(sfacg_structs.BookInfo.Status.Msg.(string))
		}
		return fmt.Errorf("book is not found")
	}
}

func GET_ACCOUNT_INFORMATION() *sfacg_structs.Account {
	return request.Get("user", &sfacg_structs.Account{}).(*sfacg_structs.Account)
}

func GET_BOOK_SHELF_INFORMATION() (map[int][]map[string]string, error) {
	s := new(bookshelf.InfoData)
	bookshelf_info := make(map[int][]map[string]string)
	request.Get(new(request.Context).Init("user/Pockets").Query("expand", "novels").QueryToString(), s)
	if s.Status.HTTPCode != 200 {
		return nil, fmt.Errorf(s.Status.Msg.(string))
	}
	for index, value := range s.Data {
		fmt.Println("书架号:", index, "\t\t\t书架名:", value.Name)
		var bookshelf_info_list []map[string]string
		for _, book := range value.Expand.Novels {
			bookshelf_info_list = append(bookshelf_info_list,
				map[string]string{"novel_name": book.NovelName, "novel_id": strconv.Itoa(book.NovelID)},
			)
		}
		bookshelf_info[index] = bookshelf_info_list
	}
	return bookshelf_info, nil
}

func GET_CATALOGUE(NovelID string) []map[string]string {
	var chapter_index int
	var division_info []map[string]string
	s := new(sfacg_structs.Catalogue)
	request.Get(new(request.Context).Init(fmt.Sprintf("novels/%v/dirs", NovelID)).
		Query("expand", "originNeedFireMoney").QueryToString(), s)

	for division_index, division := range s.Data.VolumeList {
		fmt.Printf("第%v卷\t\t%v\n", division_index+1, division.Title)
		for _, chapter := range division.ChapterList {

			chapter_index += 1
			division_info = append(division_info, map[string]string{
				"division_name":  division.Title,
				"division_id":    strconv.Itoa(division.VolumeID),
				"division_index": strconv.Itoa(division_index),
				"chapter_name":   chapter.Title,
				"chapter_id":     strconv.Itoa(chapter.ChapID),
				"chapter_index":  strconv.Itoa(chapter_index),
				"money":          strconv.Itoa(chapter.OriginNeedFireMoney),
				"file_name":      file.NameSetting(chapter.VolumeID, chapter.ChapOrder, chapter.ChapID),
			})
		}
	}
	return division_info

}

func GET_CHAPTER_CONTENT(chapter_id string) string {
	s := new(sfacg_structs.Content)
	request.Get(new(request.Context).Init("Chaps/"+chapter_id).Query("expand", "content").QueryToString(), s)
	if s != nil && s.Status.HTTPCode == 200 {
		content_title := fmt.Sprintf("%v: %v", s.Data.Title, s.Data.AddTime)
		return content_title + "\n" + tools.StandardContent(strings.Split(s.Data.Expand.Content, "\n"))

	} else {
		fmt.Println("download failed! chapterId:", chapter_id, "error:", s.Status.Msg)
	}
	return ""
}

func GET_SEARCH(keyword string, page int) *sfacg_structs.Search {
	s := new(sfacg_structs.Search)
	request.Get(new(request.Context).Init("search/novels/result").Query("q", url_.QueryEscape(keyword)).
		Query("size", "20").Query("page", strconv.Itoa(page)).QueryToString(), s)
	return s

}

func LOGIN_ACCOUNT(username, password string) {
	CookieJar := request.NewHttpUtils("sessions", "POST").Add("username", username).
		Add("password", password).NewRequests().Unmarshal(&sfacg_structs.Login).GetCookie()
	for _, cookie := range CookieJar {
		sfacg_structs.Login.Cookie += cookie.Name + "=" + cookie.Value + ";"
	}

	if sfacg_structs.Login.Status.HTTPCode == 200 {
		config.Apps.Sfacg.Cookie = sfacg_structs.Login.Cookie
		config.Apps.Sfacg.UserName = username
		config.Apps.Sfacg.Password = password
		config.SaveJson()
	} else {
		if message := sfacg_structs.Login.Status.Msg.(string); message == "用户名密码不匹配" {
			fmt.Println(message)
			os.Exit(0)
		} else {
			fmt.Println("login failed! error:", message)
		}
	}
}

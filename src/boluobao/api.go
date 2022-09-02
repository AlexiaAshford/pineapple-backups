package boluobao

import (
	"fmt"
	"github.com/VeronicaAlexia/pineapple-backups/config"
	req "github.com/VeronicaAlexia/pineapple-backups/src/https"
	_struct "github.com/VeronicaAlexia/pineapple-backups/struct"
	"github.com/VeronicaAlexia/pineapple-backups/struct/sfacg_structs"
	"github.com/VeronicaAlexia/pineapple-backups/struct/sfacg_structs/bookshelf"
	url_ "net/url"
	"strconv"
)

func GET_BOOK_INFORMATION(NovelId string) (_struct.Books, error) {
	params := map[string]string{"expand": "intro,tags,sysTags,totalNeedFireMoney,originTotalNeedFireMoney"}
	response := req.Get("novels/"+NovelId, &sfacg_structs.BookInfo{}, params).(*sfacg_structs.BookInfo)
	if response.Status.HTTPCode == 200 && response.Data.NovelName != "" {
		return _struct.Books{
			NovelName:  config.RegexpName(response.Data.NovelName),
			NovelID:    strconv.Itoa(response.Data.NovelID),
			NovelCover: response.Data.NovelCover,
			AuthorName: response.Data.AuthorName,
			CharCount:  strconv.Itoa(response.Data.CharCount),
			MarkCount:  strconv.Itoa(response.Data.MarkCount),
			SignStatus: response.Data.SignStatus,
		}, nil
	} else {
		if response.Status.Msg != nil {
			return _struct.Books{}, fmt.Errorf(response.Status.Msg.(string))
		} else {
			return _struct.Books{}, fmt.Errorf("book is not found")
		}
	}

}

func GET_ACCOUNT_INFORMATION() *sfacg_structs.Account {
	return req.Get("user", &sfacg_structs.Account{}, nil).(*sfacg_structs.Account)
}

func GET_BOOK_SHELF_INFORMATION() (map[int][]map[string]string, error) {
	params, bookshelf_info := map[string]string{"expand": "novels"}, make(map[int][]map[string]string)
	response := req.Get("user/Pockets", &bookshelf.InfoData{}, params).(*bookshelf.InfoData)
	if response.Status.HTTPCode != 200 {
		return nil, fmt.Errorf(response.Status.Msg.(string))
	}
	for index, value := range response.Data {
		fmt.Println("bookshelf index:", index, "\t\t\tbookshelf name:", value.Name)
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
	var division_info []map[string]string
	var chapter_index int
	response := req.Get(fmt.Sprintf("novels/%v/dirs", NovelID), &sfacg_structs.Catalogue{}, map[string]string{"expand": "originNeedFireMoney"})
	for division_index, division := range response.(*sfacg_structs.Catalogue).Data.VolumeList {
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
				"file_name":      config.FileCacheName(division_index, chapter_index, strconv.Itoa(chapter.ChapID)),
			})
		}
	}
	return division_info

}

func GET_CHAPTER_CONTENT(chapter_id string) string {
	params := map[string]string{"expand": "content"}
	response := req.Get("Chaps/"+chapter_id, &sfacg_structs.Content{}, params).(*sfacg_structs.Content)
	if response != nil && response.Status.HTTPCode == 200 {
		content_title := fmt.Sprintf("%v: %v", response.Data.Title, response.Data.AddTime)
		return content_title + "\n" + config.StandardContent(response.Data.Expand.Content)

	} else {
		fmt.Println("download failed! chapterId:", chapter_id, "error:", response.Status.Msg)
	}
	return ""
}

func GET_SEARCH(keyword string, page int) *sfacg_structs.Search {
	params := map[string]string{"q": url_.QueryEscape(keyword), "size": "20", "page": strconv.Itoa(page)}
	return req.Get("search/novels/result", &sfacg_structs.Search{}, params).(*sfacg_structs.Search)

}

func LOGIN_ACCOUNT(username, password string) *sfacg_structs.Login {
	params := fmt.Sprintf(`{"username":"%s", "password": "%s"}`, username, password)
	response, Cookie := req.Login(req.SET_URL("sessions", nil), []byte(params))
	for _, cookie := range Cookie {
		response.Cookie += cookie.Name + "=" + cookie.Value + ";"
	}
	return response
}

package HbookerAPI

//func GetCatalog(BookInfo *simplejson.Json) [][]string {
//	var chapterList [][]string
//	for _, division := range GetDivision(BookInfo.Get("book_id").MustString()) {
//		divisionInfo, _ := division.(map[string]interface{})
//		for _, chapter := range GetChapters(fmt.Sprintf("%v", divisionInfo["division_id"])) {
//			chapInfo, _ := chapter.(map[string]interface{})
//			cid := fmt.Sprintf("%v", chapInfo["chapter_id"])
//			if config.Isfile(path.Join("Cache", BookInfo.Get("book_name").MustString(), cid+".txt")) {
//				continue
//			} else if fmt.Sprintf("%v", chapInfo["auth_access"]) != "1" {
//				continue
//			} else {
//				chapterList = append(chapterList, []string{cid, fmt.Sprintf("%v", chapInfo["chapter_title"])})
//			}
//		}
//	}
//	return chapterList
//}
//
//func GetDivision(bid string) []interface{} {
//	response := util.Get("/book/get_division_list", map[string]string{"book_id": bid})
//	if response != nil && response.Get("code").MustInt() == 100000 {
//		return response.Get("data").Get("division_list").MustArray()
//	}
//	fmt.Println("", response.Get("tip").MustString())
//	return nil
//}
//
//func GetChapters(did string) []interface{} {
//	response := util.Get("/chapter/get_updated_chapter_by_division_id", map[string]string{"division_id": did})
//	if response != nil && response.Get("code").MustInt() == 100000 {
//		return response.Get("data").Get("chapter_list").MustArray()
//	}
//	return nil
//}

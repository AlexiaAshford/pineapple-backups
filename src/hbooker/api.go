package HbookerAPI

import (
	"encoding/json"
	"fmt"
	req "sf/src/hbooker/request"
	structs "sf/structural/hbooker_structs"
)

func GetDivisionIdByBookId(BookId string) []structs.DivisionList {
	var result structs.DivisionStruct
	if err := json.Unmarshal(req.Get(DivisionIdByBookId+BookId, 0), &result); err != nil {
		fmt.Println("json unmarshal error:", err)
		return nil
	}
	return result.Data.DivisionList
	//var chapterList []structs.ChapterList
	//for index, division := range result.Data.DivisionList {
	//	fmt.Println("index:", index, "\t\tdivision:", division.DivisionName)
	//	catalogueInfo := GetCatalogueByDivisionId(division.DivisionID)
	//	if catalogueInfo == nil {
	//		for _, chapter := range catalogueInfo {
	//			if chapter.IsValid == "1" {
	//				chapterList = append(chapterList, chapter)
	//			}
	//		}
	//	} else {
	//		return nil
	//	}
	//}
	//return chapterList
}

func GetCatalogueByDivisionId(DivisionId string) []structs.ChapterList {
	var result structs.ChapterStruct
	if err := json.Unmarshal(req.Get(CatalogueDetailedByDivisionId+DivisionId, 0), &result); err != nil {
		fmt.Println("json unmarshal error:", err)
		return nil
	}
	return result.Data.ChapterList
}

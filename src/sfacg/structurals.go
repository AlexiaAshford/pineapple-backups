package sfacg

import (
	"net/http"
	"time"
)

type UserStruct struct {
	Status struct {
		HTTPCode  int    `json:"httpCode"`
		ErrorCode int    `json:"errorCode"`
		MsgType   int    `json:"msgType"`
		Msg       string `json:"msg"`
	} `json:"status"`
	Data struct {
		UserName     string `json:"userName"`
		NickName     string `json:"nickName"`
		Email        string `json:"email"`
		AccountID    int    `json:"accountId"`
		RoleName     string `json:"roleName"`
		FireCoin     int    `json:"fireCoin"`
		Avatar       string `json:"avatar"`
		IsAuthor     bool   `json:"isAuthor"`
		PhoneNum     string `json:"phoneNum"`
		RegisterDate string `json:"registerDate"`
	} `json:"data"`
}

type FavsJson struct {
	Status struct {
		HTTPCode  int    `json:"httpCode"`
		ErrorCode int    `json:"errorCode"`
		MsgType   int    `json:"msgType"`
		Msg       string `json:"msg"`
	} `json:"status"`
}

type MoneyStruct struct {
	Status struct {
		HTTPCode  int         `json:"httpCode"`
		ErrorCode int         `json:"errorCode"`
		MsgType   int         `json:"msgType"`
		Msg       interface{} `json:"msg"`
	} `json:"status"`
	Data struct {
		RmbCost         int `json:"rmbCost"`
		FireMoneyUsed   int `json:"fireMoneyUsed"`
		FireMoneyRemain int `json:"fireMoneyRemain"`
		VipLevel        int `json:"vipLevel"`
		CouponsRemain   int `json:"couponsRemain"`
	} `json:"data"`
}

type CookieList struct {
	Cookie []struct {
		Name       string      `json:"Name"`
		Value      string      `json:"Value"`
		Path       string      `json:"Path"`
		Domain     string      `json:"Domain"`
		Expires    time.Time   `json:"Expires"`
		RawExpires string      `json:"RawExpires"`
		MaxAge     int         `json:"MaxAge"`
		Secure     bool        `json:"Secure"`
		HTTPOnly   bool        `json:"HttpOnly"`
		SameSite   int         `json:"SameSite"`
		Raw        string      `json:"Raw"`
		Unparsed   interface{} `json:"Unparsed"`
	} `json:"Cookie"`
}

type Cookie struct {
	Cookie []*http.Cookie
}

type SearchJson struct {
	Status struct {
		HTTPCode  int         `json:"httpCode"`
		ErrorCode int         `json:"errorCode"`
		MsgType   int         `json:"msgType"`
		Msg       interface{} `json:"msg"`
	} `json:"status"`
	Data struct {
		Novels []struct {
			AuthorID       int     `json:"authorId"`
			LastUpdateTime string  `json:"lastUpdateTime"`
			MarkCount      int     `json:"markCount"`
			NovelCover     string  `json:"novelCover"`
			BgBanner       string  `json:"bgBanner"`
			NovelID        int     `json:"novelId"`
			NovelName      string  `json:"novelName"`
			Point          float64 `json:"point"`
			IsFinish       bool    `json:"isFinish"`
			AuthorName     string  `json:"authorName"`
			CharCount      int     `json:"charCount"`
			ViewTimes      int     `json:"viewTimes"`
			TypeID         int     `json:"typeId"`
			AllowDown      bool    `json:"allowDown"`
			SignStatus     string  `json:"signStatus"`
			CategoryID     int     `json:"categoryId"`
			Expand         struct {
				ChapterCount       int      `json:"chapterCount"`
				BigBgBanner        string   `json:"bigBgBanner"`
				BigNovelCover      string   `json:"bigNovelCover"`
				TypeName           string   `json:"typeName"`
				Intro              string   `json:"intro"`
				Fav                int      `json:"fav"`
				Tags               []string `json:"tags"`
				SignLevel          string   `json:"signLevel"`
				TotalNeedFireMoney int      `json:"totalNeedFireMoney"`
			} `json:"expand"`
			Weight int `json:"weight"`
		} `json:"novels"`
		Comics []interface{} `json:"comics"`
		Albums []interface{} `json:"albums"`
	} `json:"data"`
}

type Pool struct {
	BookID string
}

type BookShelfJson struct {
	Status struct {
		HTTPCode  int         `json:"httpCode"`
		ErrorCode int         `json:"errorCode"`
		MsgType   int         `json:"msgType"`
		Msg       interface{} `json:"msg"`
	} `json:"status"`
	Data []struct {
		AccountID  int    `json:"accountId"`
		PocketID   int    `json:"pocketId"`
		Name       string `json:"name"`
		TypeID     int    `json:"typeId"`
		CreateTime string `json:"createTime"`
		IsFull     bool   `json:"isFull"`
		CanModify  bool   `json:"canModify"`
		Expand     struct {
			Novels []struct {
				AuthorID       int    `json:"authorId"`
				LastUpdateTime string `json:"lastUpdateTime"`
				MarkCount      int    `json:"markCount"`
				NovelCover     string `json:"novelCover"`
				BgBanner       string `json:"bgBanner"`
				NovelID        int    `json:"novelId"`
				NovelName      string `json:"novelName"`
				Point          int    `json:"point"`
				IsFinish       bool   `json:"isFinish"`
				AuthorName     string `json:"authorName"`
				CharCount      int    `json:"charCount"`
				ViewTimes      int    `json:"viewTimes"`
				TypeID         int    `json:"typeId"`
				AllowDown      bool   `json:"allowDown"`
				SignStatus     string `json:"signStatus"`
				CategoryID     int    `json:"categoryId"`
				Expand         struct {
					Discount           int    `json:"discount"`
					DiscountExpireDate string `json:"discountExpireDate"`
				} `json:"expand"`
				IsSticky       bool        `json:"isSticky"`
				StickyDateTime interface{} `json:"stickyDateTime"`
				MarkDateTime   string      `json:"markDateTime"`
			} `json:"novels"`
		} `json:"expand"`
	} `json:"data"`
}

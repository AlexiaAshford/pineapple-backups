package https

import (
	"encoding/json"
	"fmt"
	"github.com/VeronicaAlexia/pineapple-backups/config"
	"github.com/VeronicaAlexia/pineapple-backups/src/encryption"
	"io"
	"net/http"
	"strings"
)

func JsonUnmarshal(response []byte, Struct any) any {
	err := json.Unmarshal(response, Struct)
	if err != nil {
		fmt.Println(config.Error("json unmarshal", err, 18))
	}
	return Struct
}

func SET_URL(url string) string {
	switch config.Vars.AppType {
	case "cat":
		return CatWebSite + strings.ReplaceAll(url, CatWebSite, "")
	case "sfacg":
		return SFWebSite + strings.ReplaceAll(url, SFWebSite, "")
	case "happybooker":
		return HappyWebSite + strings.ReplaceAll(url, HappyWebSite, "")
	default:
		return url
	}
}

func Request(url string) []byte {
	method := "GET"
	if config.Vars.AppType == "cat" || strings.Contains(url, "session") {
		if !strings.Contains(url, "jpg") {
			method = "POST"
		}
	}
	request, _ := http.NewRequest(method, url, nil)
	SET_THE_HEADERS(request)
	if response, ok := http.DefaultClient.Do(request); ok == nil {
		result_body, _ := io.ReadAll(response.Body)
		if config.Vars.AppType == "cat" && !strings.Contains(url, "jpg") {
			return encryption.Decode(string(result_body), "")
		} else {
			return result_body
		}
	} else {
		fmt.Println(config.Error(method+":"+url, ok, 67))
	}

	return nil
}

func Get(url string, structural any) any {
	if result := Request(SET_URL(url)); result != nil {
		return JsonUnmarshal(result, structural)
	}
	return nil
}

//func Login(url string, dataJson []byte) (*sfacg_structs.Login, []*http.Cookie) {
//	request, err := http.NewRequest("POST", SET_URL(url), bytes.NewBuffer(dataJson))
//	if err != nil {
//		fmt.Println(config.Error("Login session", err, 39))
//		return nil, nil
//	}
//	SET_THE_HEADERS(request)
//	response, ok := http.DefaultClient.Do(request)
//	if ok != nil {
//		return nil, nil
//	}
//	body, _ := io.ReadAll(response.Body)
//	return JsonUnmarshal(body, &sfacg_structs.Login{}).(*sfacg_structs.Login), response.Cookies()
//}

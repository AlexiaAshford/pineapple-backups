package config

import (
	"encoding/json"
	"fmt"
	"github.com/AlexiaVeronica/boluobaoLib"
	"github.com/AlexiaVeronica/boluobaoLib/boluobaomodel"
	"github.com/AlexiaVeronica/hbookerLib"
	"github.com/AlexiaVeronica/hbookerLib/hbookermodel"
	"github.com/VeronicaAlexia/pineapple-backups/pkg/tools"
	"github.com/google/uuid"
	"os"
	"path"
	"sync"
)

func UpdateConfig() { // update config.json if necessary
	changeVar := false
	if Vars.MaxRetry == 0 || Vars.MaxRetry >= 10 {
		Vars.MaxRetry = 5 // retry times when failed
		changeVar = true
	}
	if Vars.DeviceId == "" {
		Vars.DeviceId = uuid.New().String()
		changeVar = true
	}
	if Vars.ConfigName == "" || Vars.OutputName == "" || Vars.CoverFile == "" {
		Vars.ConfigName, Vars.OutputName, Vars.CoverFile = "cache", "save", "cover"
		changeVar = true
	}
	if Apps.Hbooker.DeviceToken == "" || Apps.Hbooker.AppVersion == "" {
		Apps.Hbooker.DeviceToken, Apps.Hbooker.AppVersion = "ciweimao_", "2.9.291"
		changeVar = true
	}
	Exist([]string{Vars.ConfigName, Vars.OutputName, path.Join(Vars.ConfigName, Vars.CoverFile)})
	if changeVar {
		SaveJson()
	}
}

var APP = struct {
	Hbooker *Hbooker
	SFacg   *SFacg
}{}

type Hbooker struct {
	Client   *hbookerLib.Client
	BookInfo *hbookermodel.BookInfo
}
type SFacg struct {
	Client   *boluobaoLib.Client
	BookInfo *boluobaomodel.BookInfoData
}

func Exist(fileName any) bool {
	switch fileName.(type) {
	case string:
		_, err := os.Stat(fileName.(string))
		if os.IsNotExist(err) {
			return false
		} else {
			return true
		}
	case []string:
		for _, v := range fileName.([]string) {
			if !Exist(v) {
				tools.Mkdir(v)
			}
		}
	}
	return false
}

var FileLock = &sync.Mutex{}

func ReadConfig(fileName string) []byte {
	if fileName == "" {
		fileName = "./config.json"
	}
	// del ioutil and use os
	if data, err := os.ReadFile(fileName); err == nil {
		return data
	} else {
		fmt.Println("ReadConfig:", err)
	}
	return nil
}

func LoadJson() {
	FileLock.Lock()
	defer FileLock.Unlock()
	if ok := json.Unmarshal(ReadConfig(""), &Apps); ok != nil {
		fmt.Println("Load:", ok)
	}

}

func SaveJson() {
	//Apps.Config = Vars
	if save, ok := json.MarshalIndent(Apps, "", "    "); ok == nil {
		// del ioutil and use os
		if err := os.WriteFile("config.json", save, 0777); err != nil {
			fmt.Println("SaveJson:", err)
		}
		LoadJson()
	} else {
		fmt.Println("SaveJson:", ok)
	}
}

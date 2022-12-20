package config

import (
	"encoding/json"
	"fmt"
	"github.com/VeronicaAlexia/pineapple-backups/pkg/tools"
	"github.com/VeronicaAlexia/pineapple-backups/struct"
	"github.com/google/uuid"
	"os"
	"sync"
)

var (
	Vars    = _struct.MyConfigPro{}
	Apps    = _struct.MyAppPro{}
	Current = _struct.MyBookInfoJsonPro{}
)

func UpdateConfig() { // update config.json if necessary
	changeVar := false
	if Vars.ThreadNum == 0 || Vars.ThreadNum >= 64 {
		Vars.ThreadNum = 32 // default value is 32 thread
		changeVar = true
	}
	if Vars.MaxRetry == 0 || Vars.MaxRetry >= 10 {
		Vars.MaxRetry = 5 // retry times when failed
		changeVar = true
	}
	if Vars.AppType == "" {
		// default app type
		Vars.AppType = "cat"
		changeVar = true
	}
	if Vars.AppKey == "" {
		Vars.AppKey = "FMLxgOdsfxmN!Dt4"
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
	if Apps.Cat.Params.DeviceToken == "" || Apps.Cat.Params.AppVersion == "" {
		Apps.Cat.Params.DeviceToken, Apps.Cat.Params.AppVersion = "ciweimao_", "2.9.291"
		changeVar = true
	}
	Exist([]string{Vars.ConfigName, Vars.OutputName, Vars.CoverFile})
	if changeVar {
		SaveJson()
	}
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
	Apps.Config = Vars
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

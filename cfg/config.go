package cfg

import (
	"encoding/json"
	"fmt"
	"os"
	"sf/structural"
	"sync"
)

var (
	Vars        = structural.MyConfigPro{}
	Apps        = structural.MyAppPro{}
	CurrentBook = structural.MyBookInfoJsonPro{}
)

func UpdateConfig() bool { // update config.json if necessary
	changeVar := false
	if Vars.MaxThreadNumber == 0 || Vars.MaxThreadNumber >= 64 {
		Vars.MaxThreadNumber = 32 // default value is 32 thread
		changeVar = true
	}
	if Vars.MaxRetry == 0 || Vars.MaxRetry >= 10 {
		Vars.MaxRetry = 5 // retry times when failed
		changeVar = true
	}
	if Vars.AppType == "" {
		Vars.AppType = "sfacg"
		changeVar = true
	}
	if Apps.Sfacg.UserAgent == "" || Apps.Cat.UserAgent == "" {
		Apps.Sfacg.UserAgent = "minip_novel/1.0.70(android;11)/wxmp"
		Apps.Cat.UserAgent = "Android com.kuangxiangciweimao.novel 2.9.291"
		changeVar = true
	}
	if Vars.ConfigFile == "" || Vars.SaveFile == "" {
		Vars.ConfigFile, Vars.SaveFile = "cache", "save"
		changeVar = true
	}
	if Apps.Cat.Params.DeviceToken == "" || Apps.Cat.Params.AppVersion == "" {
		Apps.Cat.Params.DeviceToken, Apps.Cat.Params.AppVersion = "ciweimao_", "2.9.291"
		changeVar = true
	}
	Exist([]string{Vars.ConfigFile, Vars.SaveFile})
	return changeVar
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
				Mkdir(v)
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

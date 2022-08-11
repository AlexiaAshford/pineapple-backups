package cfg

import (
	"encoding/json"
	"fmt"
	"os"
	"sf/structural"
	"sync"
)

var Vars = structural.MyConfigPro{}
var Apps = structural.MyAppPro{}
var BookConfig = structural.MyBookInfoJsonPro{}

func updateConfig() {
	Load()
	if Vars.MaxRetry == 0 || Vars.MaxRetry > 10 {
		Vars.MaxRetry = 5
	}
	if Vars.MaxThreadNumber == 0 || Vars.MaxThreadNumber >= 64 {
		Vars.MaxThreadNumber = 32
	}
	if Apps.Sfacg.UserAgent == "" {
		Apps.Sfacg.UserAgent = "minip_novel/1.0.70(android;11)/wxmp"
		fmt.Println("Sfacg.UserAgent is empty, set to default value:", Apps.Sfacg.UserAgent)
	}
	if Vars.ConfigFile == "" {
		Vars.ConfigFile = "cache"
	}
	if Vars.AppType == "" {
		Vars.ConfigFile = "sfacg"
	}
	if Vars.SaveFile == "" {
		Vars.SaveFile = "save"
	}
	if Apps.Cat.UserAgent == "" {
		Apps.Cat.UserAgent = "Android com.kuangxiangciweimao.novel 2.9.290"
	}
	if !CheckFileExist(Vars.ConfigFile) {
		Mkdir(Vars.ConfigFile)
	}
	if !CheckFileExist(Vars.SaveFile) {
		Mkdir(Vars.SaveFile)
	}
	SaveJson()
}

func ConfigInit() {
	if !CheckFileExist("./config.json") || FileSize("./config.json") == 0 {
		Apps.Sfacg.UserAgent = "minip_novel/1.0.70(android;11)/wxmp"
		Apps.Cat.Params.DeviceToken = "ciweimao_"
		Apps.Cat.Params.AppVersion = "2.9.290" // hbooker app version
		Apps.Cat.UserAgent = "Android com.kuangxiangciweimao.novel 2.9.290"
		SaveJson()
	}
	updateConfig()
}

func CheckFileExist(fileName string) bool {
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return false
	}
	return true
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

func Load() {
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
		Load()
	} else {
		fmt.Println("SaveJson:", ok)
	}
}

package cfg

import (
	"encoding/json"
	"fmt"
	"os"
	"sf/structural"
	"sync"
)

var Vars = structural.MyJsonPro{}
var BookConfig = structural.MyBookInfoJsonPro{}

func updateConfig() {
	Load()
	if Vars.MaxRetry == 0 || Vars.MaxRetry > 10 {
		Vars.MaxRetry = 5
	}
	if Vars.MaxThreadNumber == 0 || Vars.MaxThreadNumber >= 64 {
		Vars.MaxThreadNumber = 32
	}
	if Vars.Sfacg.UserAgent == "" {
		Vars.Sfacg.UserAgent = "minip_novel/1.0.70(android;11)/wxmp"
		fmt.Println("Sfacg.UserAgent is empty, set to default value:", Vars.Sfacg.UserAgent)
	}
	if Vars.ConfigFile == "" {
		Vars.ConfigFile = "cache"
		fmt.Println("ConfigFile is empty, use default cache")
	}
	if Vars.AppType == "" {
		Vars.ConfigFile = "sfacg"
	}
	if Vars.SaveFile == "" {
		Vars.SaveFile = "save"
		fmt.Println("SaveFile is empty, use default save")
	}
	if Vars.Cat.UserAgent == "" {
		Vars.Cat.UserAgent = "Android com.kuangxiangciweimao.novel 2.9.290"
		fmt.Println("UserAgent is empty, use default Android com.kuangxiangciweimao.novel 2.9.290")
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
		Vars.SaveFile = "save"
		Vars.ConfigFile = "cache"
		Vars.AppType = "sfacg"
		Vars.MaxThreadNumber = 32
		Vars.MaxRetry = 5 // retry times when failed
		Vars.Sfacg.UserAgent = "minip_novel/1.0.70(android;11)/wxmp"
		Vars.Cat.Params.DeviceToken = "ciweimao_"
		Vars.Cat.Params.AppVersion = "2.9.290" // hbooker app version
		Vars.Cat.UserAgent = "Android com.kuangxiangciweimao.novel 2.9.290"
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
	if ok := json.Unmarshal(ReadConfig(""), &Vars); ok != nil {
		fmt.Println("Load:", ok)
	}

}

func SaveJson() {
	if save, ok := json.MarshalIndent(Vars, "", "    "); ok == nil {
		// del ioutil and use os
		if err := os.WriteFile("config.json", save, 0777); err != nil {
			fmt.Println("SaveJson:", err)
		}
		Load()
	} else {
		fmt.Println("SaveJson:", ok)
	}
}

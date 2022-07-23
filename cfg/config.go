package cfg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sf/structural"
	"sync"
)

var Vars = structural.MyJsonPro{}

func ConfigInit() {
	if !CheckFileExist("./config.json") || FileSize("./config.json") == 0 {
		Vars.SaveFile = "save"
		Vars.ConfigFile = "cache"
		Vars.Cat.UserAgent = "Android com.kuangxiangciweimao.novel 2.9.290"
		Vars.Cat.Params.DeviceToken = "ciweimao_"
		Vars.Cat.Params.AppVersion = "2.9.290"
		SaveJson()
	}
	Load()
	if Vars.ConfigFile == "" {
		Vars.ConfigFile = "cache"
		fmt.Println("ConfigFile is empty, use default cache")
	}
	if Vars.SaveFile == "" {
		Vars.SaveFile = "save"
		fmt.Println("SaveFile is empty, use default save")
	}
	if Vars.Cat.UserAgent == "" {
		Vars.Cat.UserAgent = "Android com.kuangxiangciweimao.novel 2.9.290"
		fmt.Println("UserAgent is empty, use default Android com.kuangxiangciweimao.novel 2.9.290")
	}
	//if !CheckFileExist(Vars.ConfigFile) {
	//	fmt.Println("ConfigFile not exist, create it now ...")
	//	Mkdir(Vars.ConfigFile)
	//}
	if !CheckFileExist(Vars.SaveFile) {
		Mkdir(Vars.SaveFile)
	}
}

func CheckFileExist(fileName string) bool {
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

var FileLock = &sync.Mutex{}

func Load() {
	FileLock.Lock()
	defer FileLock.Unlock()
	if data, err := ioutil.ReadFile("config.json"); err == nil {
		if ok := json.Unmarshal(data, &Vars); ok != nil {
			fmt.Println("Load:", ok)
		}
	} else {
		fmt.Println("Load:", err)
	}
}

func SaveJson() {
	if save, ok := json.MarshalIndent(Vars, "", "    "); ok == nil {
		if err := ioutil.WriteFile("config.json", save, 0777); err != nil {
			fmt.Println("SaveJson:", err)
		}
		Load()
	} else {
		fmt.Println("SaveJson:", ok)
	}
}

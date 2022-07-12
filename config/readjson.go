package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

type MyJsonPro struct {
	ConfigFile    string
	SaveFile      string
	UserName      string
	Password      string
	authKey       string
	Authorization string
	DeviceToken   string
	UserAgent     string
	CatToken      string
	Account       string
	Cookie        string
}

var Var = MyJsonPro{}

func NewMyJsonPro() {
	if !CheckFileExist("config.json") {
		Var.SaveFile = "save"
		Var.ConfigFile = "cache"
		SaveJson()
	}
	Load()
	Mkdir(Var.SaveFile)
	Mkdir(Var.ConfigFile)
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
		if ok := json.Unmarshal(data, &Var); ok != nil {
			fmt.Println("Load:", ok)
		}
	} else {
		fmt.Println("Load:", err)
	}
}

func SaveJson() {
	if save, ok := json.MarshalIndent(Var, "", "    "); ok == nil {
		if err := ioutil.WriteFile("config.json", save, 0777); err != nil {
			fmt.Println("SaveJson:", err)
		}
	} else {
		fmt.Println("SaveJson:", ok)
	}
}

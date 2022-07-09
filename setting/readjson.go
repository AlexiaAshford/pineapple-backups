package setting

import (
	"encoding/json"
	"github.com/google/uuid"
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
	Cookie        map[string]string
}

var (
	Var       = MyJsonPro{}
	ConfigDir = "config.json"
)

func NewMyJsonPro() {
	if !CheckFileExist(ConfigDir) {
		Var.SaveFile = "save"
		Var.ConfigFile = "config"
		Var.DeviceToken = uuid.New().String()
		Var.Authorization = "Basic YXBpdXNlcjozcyMxLXl0NmUqQWN2QHFlcg=="
		Var.UserAgent = "boluobao/4.8.66(iOS;15.4.1)/appStore/" + Var.DeviceToken
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
	if data, err := ioutil.ReadFile(ConfigDir); err == nil {
		if ok := json.Unmarshal(data, &Var); ok != nil {
			println("Load:", ok)
		}
	} else {
		println("Load:", err)
	}
}

func SaveJson() {
	if save, ok := json.MarshalIndent(Var, "", "    "); ok == nil {
		if err := ioutil.WriteFile(ConfigDir, save, 0777); err != nil {
			println("SaveJson:", err)
		}
	} else {
		println("SaveJson:", ok)
	}
}

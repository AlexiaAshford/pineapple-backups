package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"sync"

	"github.com/AlexiaVeronica/boluobaoLib"
	"github.com/AlexiaVeronica/boluobaoLib/boluobaomodel"
	"github.com/AlexiaVeronica/hbookerLib"
	"github.com/AlexiaVeronica/hbookerLib/hbookermodel"
	"github.com/AlexiaVeronica/pineapple-backups/pkg/tools"
	"github.com/google/uuid"
)

var (
	Apps     = AppConfig{}
	Vars     = &Apps.Config
	FileLock = &sync.Mutex{}
)

type AppConfig struct {
	Hbooker HbookerCommonParams `json:"common_params"`
	Sfacg   BoluobaoConfig      `json:"sfacg_config"`
	Config  ScriptConfig        `json:"script_config"`
}

var HelpMessage = []string{"input help to see the command list:",
	"input quit to quit",
	"input download <bookid/url> to download book",
	"input search <keyword> to search book",
	"input show to show config",
	"input update config to update config by config.json",
	"input login <account> <password> to login account",
	"input app <app app keyword> to change app type",
	"input max <thread> to change max thread number",
	"you can input command like this: download <bookid/url>\n\n",
}

type ScriptConfig struct {
	ConfigName string `json:"config_name"`
	OutputName string `json:"output_name"`
	CoverFile  string `json:"cover_file"`
	DeviceId   string `json:"device_id"`
	ThreadNum  int    `json:"thread_num"`
	MaxRetry   int    `json:"max_retry"`
	Epub       bool   `json:"epub"`
}

type BoluobaoConfig struct {
	Cookie string `json:"cookie"`
}

type HbookerCommonParams struct {
	LoginToken string `json:"login_token"`
	Account    string `json:"account"`
}

type Hbooker struct {
	Client   *hbookerLib.Client
	BookInfo *hbookermodel.BookInfo
}

type SFacg struct {
	Client   *boluobaoLib.Client
	BookInfo *boluobaomodel.BookInfoData
}

var APP = struct {
	Hbooker *Hbooker
	SFacg   *SFacg
}{}

func UpdateConfig() {
	changeVar := false

	if Vars.MaxRetry <= 0 || Vars.MaxRetry >= 10 {
		Vars.MaxRetry = 5
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

	EnsureDirectoriesExist([]string{Vars.ConfigName, Vars.OutputName, path.Join(Vars.ConfigName, Vars.CoverFile)})

	if changeVar {
		SaveConfig()
	}
}

func EnsureDirectoriesExist(dirNames []string) {
	for _, dir := range dirNames {
		if !DirectoryExists(dir) {
			tools.Mkdir(dir)
		}
	}
}

func DirectoryExists(dirName string) bool {
	_, err := os.Stat(dirName)
	return !os.IsNotExist(err)
}

func ReadConfig(fileName string) ([]byte, error) {
	if fileName == "" {
		fileName = "./config.json"
	}
	return os.ReadFile(fileName)
}

func LoadConfig() {
	FileLock.Lock()
	defer FileLock.Unlock()

	data, err := ReadConfig("")
	if err != nil {
		fmt.Println("ReadConfig:", err)
		return
	}

	if err := json.Unmarshal(data, &Apps); err != nil {
		fmt.Println("Load:", err)
	}
}

func SaveConfig() {
	FileLock.Lock()
	defer FileLock.Unlock()

	data, err := json.MarshalIndent(Apps, "", "    ")
	if err != nil {
		fmt.Println("SaveConfig:", err)
		return
	}

	if err := os.WriteFile("config.json", data, 0777); err != nil {
		fmt.Println("SaveConfig:", err)
	}

	LoadConfig()
}

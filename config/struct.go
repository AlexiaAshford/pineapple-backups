package config

var Apps = AppConfig{}
var Vars = &Apps.Config

type AppConfig struct {
	Hbooker HbookerCommonParams `json:"common_params"`
	Sfacg   BoluobaoConfig      `json:"sfacg_config"`
	Config  ScriptConfig        `json:"script_config"`
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

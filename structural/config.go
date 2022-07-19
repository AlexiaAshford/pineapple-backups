package structural

type MyJsonPro struct {
	ConfigFile string `json:"ConfigFile"`
	SaveFile   string `json:"SaveFile"`
	Sfacg      struct {
		UserName string `json:"UserName"`
		Password string `json:"Password"`
		Cookie   string `json:"Cookie"`
	} `json:"sfacg"`
	Cat struct {
		Token     string `json:"Token"`
		Account   string `json:"Account"`
		UserAgent string `json:"UserAgent"`
	} `json:"cat"`
}

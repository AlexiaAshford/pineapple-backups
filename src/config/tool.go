package config

import (
	"fmt"
	"os"
	"path"
	"regexp"
)

func RegexpName(Name string) string {
	return regexp.MustCompile(`[\\/:*?"<>|]`).ReplaceAllString(Name, "")
}

func Mkdir(filePath string) {
	dirPath, _ := os.Getwd()
	if err := os.MkdirAll(path.Join(dirPath, filePath), os.ModePerm); err != nil {
		fmt.Println(err)
	}
}

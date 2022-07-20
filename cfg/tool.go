package cfg

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

// input init
func InputInt(introduction string) int {
	var input int
	// if search keyword is not empty, search book and download
	fmt.Printf(introduction)
	if _, err := fmt.Scanln(&input); err == nil {
		return input
	} else {
		fmt.Println(err)
		InputInt(introduction)
	}
	fmt.Println("something wrong, return 0")
	return 0
}

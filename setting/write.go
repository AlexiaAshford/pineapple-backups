package setting

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

func ReadFile(filePath string) string {
	if f, err := ioutil.ReadFile(filePath); err == nil {
		return string(f)
	} else {
		fmt.Println("ReadFile", err)
	}
	return ""
}

func Mkdir(filePath string) {
	dirPath, _ := os.Getwd()
	if err := os.MkdirAll(path.Join(dirPath, filePath), os.ModePerm); err != nil {
		fmt.Println(err)
	}
}
func Isfile(filePath string) bool {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return true
	} else {
		return false
	}
}

func NewFile(filePath string) {
	if Isfile(filePath) {
		file2, _ := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0766)
		if err := file2.Close(); err != nil {
			fmt.Println(err)
		}
	}
}

func Writes(filePath string, content string) {
	file, err := os.OpenFile(filePath, os.O_WRONLY, 0644)
	defer func(file *os.File) {
		fmt.Println(err)
	}(file)
	if err != nil {
		fmt.Println(err)
	} else {
		n, _ := file.Seek(0, 2)
		if _, err = file.WriteAt([]byte(content), n); err != nil {
			println(err)
		}
	}

}

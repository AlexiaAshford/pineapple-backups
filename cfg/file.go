package cfg

import (
	"fmt"
	"io"
	"os"
)

// WriteFile write content to file with file name
func WriteFile(fileName string, content string, perm os.FileMode) error {
	if perm == 0644 {
		if _, err := os.Create(fileName); err != nil {
			fmt.Println("Create file error:", err)
		}
	}
	if f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND, perm); err == nil {
		defer func(f *os.File) {
			err = f.Close()
			if err != nil {
				fmt.Println("file close failed. err: " + err.Error())
			}
		}(f)
		if _, ok := f.WriteString(content); ok != nil {
			fmt.Println("file write failed. err: " + ok.Error())
			return ok
		}
	} else {
		fmt.Println("file open failed. err: " + err.Error())
		return err
	}
	return nil
}
func Write(Path string, content string, permMode string) string {
	var perm os.FileMode
	if permMode == "w" {
		perm = 0644
	} else if permMode == "a" {
		perm = 0666
	} else if permMode == "r" {
		return ReadFile(Path)
	} else {

		panic("permMode error")
	}
	for i := 0; i < Vars.MaxRetry; i++ {
		if WriteFile(Path, content, perm) == nil {
			break
		} else {
			fmt.Println("write file error, try again:", i)
		}
	}
	return ""
}
func FileSize(FilePath string) int {
	if file, err := os.Open(FilePath); err != nil {
		fmt.Println("open file error:", err)
	} else {
		sum := 0
		buf := make([]byte, 2014)
		for {
			n, err := file.Read(buf)
			sum += n
			if err == io.EOF {
				break
			}
		}
		return sum
	}
	return 0
}

// ReadFile read file content to string
func ReadFile(filePath string) string {
	if content, err := os.ReadFile(filePath); err == nil {
		return string(content)
	} else {
		fmt.Println("ReadFile:", err)
	}
	return ""
}

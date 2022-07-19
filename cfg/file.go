package cfg

import (
	"fmt"
	"os"
)

// WriteFile write content to file with file name
func WriteFile(fileName string, content string, perm os.FileMode) error {
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

package tool

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func RegexpName(Name string) string {
	return regexp.MustCompile(`[\\/:*?"<>|]`).ReplaceAllString(Name, "")
}

func JsonString(jsonStruct any) string {
	if jsonInfo, err := json.MarshalIndent(jsonStruct, "", "    "); err == nil {
		return string(jsonInfo)
	} else {
		log.Println(err)
	}
	return ""
}

func StandardContent(contentList []string) string {
	content := "" // clear content string
	for _, s := range contentList {
		if s != "" {
			content += "\n" + strings.ReplaceAll(s, " ", "")
		}
	}
	return content
}
func TestList(List []string, testString string) bool {
	for _, s := range List {
		if s == testString {
			return true
		}
	}
	return false
}
func TestIntList(List []int, testString string) bool {
	for _, s := range List {
		if strconv.Itoa(s) == testString {
			return true
		}
	}
	return false
}
func GetFileName(dirname string) []string {
	var file_list []string
	f, err := os.Open(dirname)
	if err != nil {
		log.Fatal(err)
	}
	if list, ok := f.Readdir(-1); ok == nil {
		_ = f.Close()
		sort.Slice(list, func(i, j int) bool { return list[i].Name() < list[j].Name() })
		for _, v := range list {
			file_list = append(file_list, v.Name())
		}
		return file_list
	} else {
		log.Fatal(ok)
	}
	return nil
}

//func ColorPrint(s string, i int) {
//	//set color and print
//	kernel32 := syscall.NewLazyDLL("kernel32.dll")
//	proc := kernel32.NewProc("SetConsoleTextAttribute")
//	handle, _, _ := proc.Call(uintptr(syscall.Stdout), uintptr(i))
//	fmt.Print(s)
//	handle, _, _ = proc.Call(uintptr(syscall.Stdout), uintptr(7))
//	CloseHandle := kernel32.NewProc("CloseHandle")
//	_, _, _ = CloseHandle.Call(handle)
//}

func get_working_directory() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	} else {
		return dir
	}
	return ""
}

func Mkdir(filePath string) string {
	file_path := path.Join(get_working_directory(), filePath)
	if err := os.MkdirAll(file_path, os.ModePerm); err != nil {
		fmt.Println(err)
	}
	return file_path
}

// InputStr input str
func InputStr(introduction string) string {
	var input string
	// if search keyword is not empty, search book and download
	fmt.Printf(introduction)
	if _, err := fmt.Scanln(&input); err == nil {
		if input != "" {
			return input
		}
	}
	return InputStr(">")
}

func GET(prompt string) []string {
	compile, _ := regexp.Compile(`\s+`)
	inputs := compile.Split(strings.TrimSpace(Input(prompt)), -1)
	if len(inputs) > 0 && inputs[0] != "" {
		return inputs
	}
	return nil
}

func Input(prompt string) string {
	for {
		fmt.Printf(prompt)
		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err == nil {
			return input
		}
	}
}
func IsNum(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// StrToInt string to int
func StrToInt(str string) int {
	if i, err := strconv.Atoi(str); err == nil {
		return i
	}
	return 0
}

func FormatJson(jsonString []byte) {
	var str bytes.Buffer
	if err := json.Indent(&str, jsonString, "", "    "); err == nil {
		fmt.Println(str.String())
	} else {
		log.Fatalln(err)
	}
}

// InputInt input int
func InputInt(introduction string, max_indexes int) int {
	var input int
	// if search keyword is not empty, search book and download
	fmt.Printf(introduction)
	if _, err := fmt.Scanln(&input); err == nil {
		for {
			if input >= max_indexes {
				fmt.Println("you input index is out of range, please input again:")
				return InputInt(">", max_indexes)
			} else {
				return input
			}
		}
	} else {
		return InputInt(">", max_indexes)
	}
}

//func TestKeyword(Text string, keyword any) bool {
//	switch keyword.(type) {
//	case string:
//		return strings.Contains(Text, keyword.(string))
//	case int:
//		return strings.Contains(Text, strconv.Itoa(keyword.(int)))
//	default:
//		panic("keyword type error")
//	}
//}

// strconv.FormatBool()
//func FormatBool(b bool) string {
//	if b {
//		return "true"
//	}
//	return "false"
//}

package cfg

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
	"syscall"
)

func RegexpName(Name string) string {
	return regexp.MustCompile(`[\\/:*?"<>|]`).ReplaceAllString(Name, "")
}
func StandardContent(content string) string {
	content_list := strings.Split(content, "\n")
	content = "" // clear content string
	for _, s := range content_list {
		if s != "" {
			content += "\n" + strings.ReplaceAll(s, " ", "")
		}
	}
	return content
}
func FormatBool(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
func Config_file_name(index, chapter_index, ChapID any) string {
	index = StrToInt(fmt.Sprintf("%d", index))
	return fmt.Sprintf("%05d", index) + "-" + fmt.Sprintf("%05d", chapter_index) + "-" +
		fmt.Sprintf("%v", ChapID) + ".txt"
}
func TestList(List []string, testString string) bool {
	for _, s := range List {
		if s == testString {
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

func ColorPrint(s string, i int) {
	//set color and print
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	proc := kernel32.NewProc("SetConsoleTextAttribute")
	handle, _, _ := proc.Call(uintptr(syscall.Stdout), uintptr(i))
	fmt.Print(s)
	handle, _, _ = proc.Call(uintptr(syscall.Stdout), uintptr(7))
	CloseHandle := kernel32.NewProc("CloseHandle")
	_, _, _ = CloseHandle.Call(handle)
}

func ExtractBookID(url string) string {
	if url != "" {
		bookID := regexp.MustCompile(`(\d+)`).FindStringSubmatch(url)
		if len(bookID) > 1 {
			return bookID[1]
		}
	}
	return ""
}

func Mkdir(filePath string) {
	dirPath, _ := os.Getwd()
	if err := os.MkdirAll(path.Join(dirPath, filePath), os.ModePerm); err != nil {
		fmt.Println(err)
	}
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

func Input(prompt string) string {
	for {
		fmt.Printf(prompt)
		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err == nil {
			return input
		}
		return ""
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
	} else {
		return 0
	}
}

func FormatJson(jsonString []byte) {
	var str bytes.Buffer
	if err := json.Indent(&str, jsonString, "", "    "); err == nil {
		fmt.Println(str.String())
	} else {
		log.Fatalln(err)
	}
}
func TestKeyword(Text string, keyword any) bool {
	switch keyword.(type) {
	case string:
		return strings.Contains(Text, keyword.(string))
	case int:
		return strings.Contains(Text, strconv.Itoa(keyword.(int)))
	default:
		panic("keyword type error")
	}
}

// input int
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

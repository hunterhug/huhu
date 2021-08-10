package util

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func CurDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return filepath.Dir(filename)
}

func SaveToFile(filepath string, content []byte) error {
	err := ioutil.WriteFile(filepath, content, 0777)
	return err
}

func ReadFromFile(filepath string) ([]byte, error) {
	return ioutil.ReadFile(filepath)
}

func MakeDir(fileDir string) error {
	return os.MkdirAll(fileDir, 0777)
}

func MakeDirByFile(filepath string) error {
	temp := strings.Split(filepath, "/")
	if len(temp) <= 2 {
		return errors.New("please input complete file name like ./dir/filename or /home/dir/filename")
	}
	dirPath := strings.Join(temp[0:len(temp)-1], "/")
	return MakeDir(dirPath)
}

func FileExist(filename string) bool {
	fi, err := os.Stat(filename)
	if err != nil {
		return false
	}
	if fi.IsDir() {
		return false
	} else {
		return true
	}
}

func ListDir(dirPth string, suffix string) (files []string, err error) {
	files = make([]string, 0, 10)
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}
	suffix = strings.ToUpper(suffix)
	for _, fi := range dir {
		if fi.IsDir() {
			continue
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) { //匹配文件
			files = append(files, dirPth+"/"+fi.Name())
		}
	}
	return files, nil
}

func ValidFileName(filename string) string {
	patterns := []string{
		" ", "#01#",
		"\\", "#02#",
		"/", "#03#",
		":", "#04#",
		"\"", "#05#",
		"?", "#06#",
		"<", "#07#",
		">", "#08#",
		"|", "#09#",
	}
	r := strings.NewReplacer(patterns...)
	dudu := r.Replace(filename)
	return dudu
}

func ValidBackFileName(filename string) string {
	patterns := []string{
		"#01#", " ",
		"#02#", "\\",
		"#03#", "/",
		"#04#", ":",
		"#05#", "\"",
		"#06#", "?",
		"#07#", "<",
		"#08#", ">",
		"#09#", "|",
	}
	r := strings.NewReplacer(patterns...)
	dudu := r.Replace(filename)
	return dudu
}

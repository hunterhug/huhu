package util

import (
	"errors"
	"strconv"
	"strings"
)

func SI(s string) (i int, e error) {
	i, e = strconv.Atoi(s)
	return
}

func IS(i int) string {
	return strconv.Itoa(i)
}

func ToLower(s string) string {
	return strings.ToLower(s)
}

func DevideStringList(files []string, num int) (map[int][]string, error) {
	length := len(files)
	split := map[int][]string{}
	if num <= 0 {
		return split, errors.New("num must not negtive")
	}
	if num > length {
		num = length
		//return split, errors.New("num must not bigger than the list length")
	}
	process := length / num
	for i := 0; i < num; i++ {
		// slice inside has a refer, so must do this append
		//split[i]=files[i*process : (i+1)*process] wrong!
		split[i] = append(split[i], files[i*process:(i+1)*process]...)
	}
	remain := files[num*process:]
	for i := 0; i < len(remain); i++ {
		split[i%num] = append(split[i%num], remain[i])
	}
	return split, nil
}

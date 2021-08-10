package util

import (
	"encoding/json"
)

// JsonBack 最好的方法
func JsonBack(s []byte) ([]byte, error) {
	temp := new(interface{})
	json.Unmarshal(s, temp)
	result, err := json.Marshal(temp)
	return result, err
}

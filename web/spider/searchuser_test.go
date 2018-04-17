package spider

import (
	"testing"
	//"github.com/hunterhug/marmot/miner"
	"github.com/hunterhug/huhu/src"
	"fmt"
)

func init() {
	//miner.SetLogLevel("debug")
}

func TestSearchUser(t *testing.T) {
	e := src.SetCookie("../../cookie1.txt")
	if e != nil {
		fmt.Println(e.Error())
		return
	}
	StartSpider("hunterhug", 0)
}

func TestSearchUser2(t *testing.T) {
	e := src.SetCookie("../../cookie1.txt")
	if e != nil {
		fmt.Println(e.Error())
		return
	}
	SearchUser("excited-vczh")
}

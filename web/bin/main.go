package main

import (
	"github.com/hunterhug/huhu/src"
	"fmt"
	"github.com/hunterhug/huhu/web/spider"
	"github.com/hunterhug/parrot/util"
)

func main() {
	CurDir, _ := util.GetBinaryCurrentPath()
	e := src.SetCookie(CurDir + "/cookie1.txt")
	if e != nil {
		fmt.Println(e.Error())
		return
	}
	spider.StartSpider("excited-vczh", 0)

}

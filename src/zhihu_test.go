package src

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

func TestLogin(t *testing.T) {
	//u := Input("请输入手机/邮箱:", "569929309@qq.com")
	//p := Input("请输入密码:", "txxxx6")
	u := "ddd"
	p := "ddd"
	body, err := Login(u, p)
	if err != nil {
		fmt.Printf("错误：%s，请手动关闭\n", err.Error())
		time.Sleep(50 * time.Second)
		os.Exit(0)
	}
	if strings.Contains(string(body), "验证") {
		fmt.Println(string(body))
		fmt.Println("错误，请手动写cookie.txt")
	} else {
		fmt.Println(string(body))
	}
}

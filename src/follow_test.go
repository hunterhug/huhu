package src

import (
	"fmt"
	"strings"
	"testing"
)

// catch粉丝
func TestCatchUser(t *testing.T) {
	e := SetCookie("../cookie1.txt")
	if e != nil {
		fmt.Println(e.Error())
		return
	}
	a, e := CatchUser(false, "hunterhug", 20, 0)
	if e != nil {
		fmt.Println(e.Error())
	} else {
		//r, e := Baba.JsonToString()
		//if e != nil {
		//	fmt.Println(e.Error())
		//}
		//fmt.Println(r)
		rr := ParseUser(a)
		for _, v := range rr.Data {
			fmt.Printf("%v,%v,https://www.zhihu.com/people/%v\n", strings.Replace(v.Name, ",", ".", -1), v.Gender, v.UrlToken)
		}
	}
}

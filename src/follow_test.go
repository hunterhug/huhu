/*
   Created by jinhan on 17-8-14.
   Tip:
   Update:
*/
package src

import (
	"fmt"
	"strings"
	"testing"
)

// catch粉丝
func TestCatchUser(t *testing.T) {
	e := SetCookie("/home/jinhan/cookie.txt")
	if e != nil {
		fmt.Println(e.Error())
	}
	a, e := CatchUser(true, "hunterhug", 20, 0)
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

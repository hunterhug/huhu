package core

import (
	"fmt"
	"testing"

	"github.com/hunterhug/huhu/util"
)

func TestOneOutputHtml(t *testing.T) {
	a, e := util.ReadFromFile("../data/28467579/shi-xun-wang-zi-60316405/shi-xun-wang-zi-60316405的回答.html")
	if e != nil {
		fmt.Println(e.Error())
	} else {
		fmt.Println(OneOutputHtml(string(a)))
	}
}

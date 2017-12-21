/*
   Created by jinhan on 17-9-4.
   Tip:
   Update:
*/
package src

import (
	"fmt"
	"github.com/hunterhug/parrot/util"
	"testing"
)

func TestOneOutputHtml(t *testing.T) {
	a, e := util.ReadfromFile("../data/28467579/shi-xun-wang-zi-60316405/shi-xun-wang-zi-60316405的回答.html")
	if e != nil {
		fmt.Println(e.Error())
	} else {
		fmt.Println(OneOutputHtml(string(a)))
	}
}

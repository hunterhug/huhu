/*
   Created by jinhan on 17-10-20.
   Tip:
   Update:
*/
package src

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/go_tool/util"
	"github.com/hunterhug/marmot/expert"
)

func TestCatchPeopleAnswer(t *testing.T) {
	who := "da-xiong-nu-da-xiong-nu"
	page := 1
	d, e := CatchPeopleAnswer(who, page)
	if e != nil {
		fmt.Println(e.Error())
	} else {
		util.SaveToFile(filepath.Join(util.CurDir(), "data/oquestion.html"), d)

		doc, _ := expert.QueryBytes(d)
		text, ok := doc.Find("div#data").Attr("data-state")
		if ok {
			util.SaveToFile(filepath.Join(util.CurDir(), "data/oquestion.json"), []byte(text))
		}
		fmt.Println(string(text))
	}
}

func TestParsePeopleAnswer(t *testing.T) {
	who := "liu-yuan-ming-89"
	page := 1
	d, e := CatchPeopleAnswer(who, page)
	if e != nil {
		fmt.Println(e.Error())
	} else {
		fmt.Printf("%#v\n", ParsePeopleAnswer(d))
	}
}

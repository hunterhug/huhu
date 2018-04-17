package src

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/hunterhug/marmot/expert"
	"github.com/hunterhug/parrot/util"
)

func TestCatchPeopleAnswer(t *testing.T) {
	who := "hunterhug"
	page := 1
	d, e := CatchPeopleAnswer(who, page)
	if e != nil {
		fmt.Println(e.Error())
	} else {
		util.SaveToFile(filepath.Join(util.CurDir(), "data/oquestion2.html"), d)

		doc, _ := expert.QueryBytes(d)
		text, ok := doc.Find("div#data").Attr("data-state")
		if ok {
			util.SaveToFile(filepath.Join(util.CurDir(), "data/oquestion2.json"), []byte(text))
		}
		fmt.Println(string(text))
	}
}

func TestParsePeopleAnswer(t *testing.T) {
	who := "hunterhug"
	page := 1
	d, e := CatchPeopleAnswer(who, page)
	if e != nil {
		fmt.Println(e.Error())
	} else {
		uinfo := ParsePeopleAnswer(d)
		fmt.Printf("%v\n%#v", uinfo.Entities.Users, uinfo.Entities.Answers)
	}
}

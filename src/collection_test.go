package src

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/hunterhug/parrot/util"
)

func TestCatchCoolection(t *testing.T) {
	e := SetCookie("/home/jinhan/cookie.txt")
	if e != nil {
		fmt.Println(e.Error())
	}
	b, e := CatchCoolection(78172986, 2)
	if e != nil {
		fmt.Println(e.Error())
	} else {
		util.SaveToFile(filepath.Join(util.CurDir(), "data/collection.html"), []byte(b))
	}
}

func TestParseCollection(t *testing.T) {
	body, err := util.ReadfromFile(filepath.Join(util.CurDir(), "data/collection.html"))
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("%#v.", ParseCollection(body))

}

func TestCatchAllCollection(t *testing.T) {
	fmt.Printf("%#v,", CatchAllCollection(78172986))
}

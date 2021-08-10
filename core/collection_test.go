package core

import (
	"fmt"
	"github.com/hunterhug/marmot/miner"
	"path/filepath"
	"testing"

	"github.com/hunterhug/huhu/util"
)

func TestCatchCollection(t *testing.T) {
	miner.SetLogLevel(miner.DEBUG)
	util.MakeDir(filepath.Join(util.CurDir(), "data"))
	b, e := CatchCollection(78172986, 2)
	if e != nil {
		fmt.Println(e.Error())
	} else {
		util.SaveToFile(filepath.Join(util.CurDir(), "data/collection.json"), b)
	}
}

func TestParseCollection(t *testing.T) {
	body, err := util.ReadFromFile(filepath.Join(util.CurDir(), "data/collection.json"))
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("%#v.", ParseCollectionNew(body))

}

func TestCatchAllCollection(t *testing.T) {
	miner.SetLogLevel(miner.DEBUG)
	c := CatchAllCollection(78172986)
	fmt.Println(len(c))
	fmt.Printf("%#v", c)
}

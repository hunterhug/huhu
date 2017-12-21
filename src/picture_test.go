/*
   Created by jinhan on 17-10-20.
   Tip:
   Update:
*/
package src

import (
	"github.com/hunterhug/parrot/util"
	"path/filepath"
	"testing"
)

func TestSavePicture(t *testing.T) {
	body, err := util.ReadfromFile(filepath.Join(util.CurDir(), "data/question.html"))
	dir := filepath.Join(util.CurDir(), "data/00/00")
	if err != nil {
		panic(err.Error())
	}
	SetSavePicture(true, 50)
	SavePicture(dir, body)
}

/*
Copyright 2017 by GoSpider author.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License
*/
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

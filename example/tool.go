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
package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/hunterhug/parrot/util"
)

func main() {
	ss := `	<li><a class="page-link" href="/zhihu/%s-html/1.html">%s</a></li>`
	fs, err := util.ListDir(filepath.Join("..", "data"), ".xx")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		for _, f := range fs {
			f = strings.Replace(f, "../data/", "", -1)
			dudu := strings.Split(f, "-")
			fmt.Println(fmt.Sprintf(ss, dudu[0], strings.Replace(dudu[1], ".xx", "", -1)))
		}
	}
}

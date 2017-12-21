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

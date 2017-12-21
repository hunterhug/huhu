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
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/hunterhug/marmot/expert"
)

// 抓取收藏夹第几页列表
func CatchCoolection(id, page int) ([]byte, error) {
	Baba.SetUrl(fmt.Sprintf(CollectionUrl, id, page))
	return Baba.Get()
}

// 抓取全部收藏夹页数,并返回问题ID和标题
func CatchAllCollection(id int) map[string]string {
	returns := map[string]string{}
	i := 1
	for {
		body, err := CatchCoolection(id, i)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		fmt.Printf("抓取收藏夹第%d页\n", i)
		i = i + 1
		maps := ParseCollection(body)
		if len(maps) == 0 {
			break
		}
		for id, q := range maps {
			returns[id] = q
		}
	}
	return returns
}

// 解析收藏夹，返回问题ID和标题
func ParseCollection(body []byte) map[string]string {
	returns := map[string]string{}
	doc, _ := expert.QueryBytes(body)
	//zm-item-title
	doc.Find(".zm-item-title").Each(func(num int, node *goquery.Selection) {
		qa := node.Find("a")
		q, ok := qa.Attr("href")
		if ok {
			returns[strings.Replace(q, "/question/", "", -1)] = qa.Text()
		}
	})
	return returns
}

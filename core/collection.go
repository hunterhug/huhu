package core

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hunterhug/marmot/expert"
	"github.com/hunterhug/marmot/util/goquery"
)

// CatchCollection 抓取收藏夹第几页列表
func CatchCollection(id, page int) ([]byte, error) {
	w := Baba.Clone()
	w.SetUrl(fmt.Sprintf(CollectionUrlNew, id, (page-1)*400))
	return w.Get()
}

// CatchAllCollection 抓取全部收藏夹页数,并返回问题ID和标题
func CatchAllCollection(id int) map[string]string {
	returns := map[string]string{}
	i := 1
	for {
		body, err := CatchCollection(id, i)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		fmt.Printf("抓取收藏夹第%d页\n", i)
		i = i + 1
		maps := ParseCollectionNew(body)
		if len(maps) == 0 {
			break
		}
		for mid, q := range maps {
			returns[mid] = q
		}
	}
	return returns
}

// ParseCollection 解析收藏夹，返回问题ID和标题
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

type Collection struct {
	Data []CollectionOne `json:"data"`
}

type CollectionOne struct {
	Content CollectionOneContent `json:"content"`
}

type CollectionOneContent struct {
	Id       int64        `json:"id"`
	Question QuestionInfo `json:"question"`
}

func ParseCollectionNew(body []byte) map[string]string {
	m := map[string]string{}
	collect := new(Collection)
	json.Unmarshal(body, collect)
	for _, v := range collect.Data {
		if v.Content.Question.Qid != 0 {
			m[fmt.Sprintf("%d", v.Content.Question.Qid)] = v.Content.Question.Title
		}
	}

	return m
}

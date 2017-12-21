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

func TestCatchAnswer(t *testing.T) {
	e := SetCookie("/home/jinhan/cookie.txt")
	if e != nil {
		fmt.Println(e.Error())
	}
	b, e := CatchAnswer(Question("28467579"), 1, 1)
	if e != nil {
		fmt.Println(e.Error())
		data, e1 := JsonBack(b)
		fmt.Println(string(data), e1)
	} else {
		e := util.SaveToFile(filepath.Join(util.CurDir(), "data/question.json"), []byte(b))
		fmt.Printf("%#v", e)
	}
}

func TestStructAnswer(t *testing.T) {
	body, err := util.ReadfromFile(filepath.Join(util.CurDir(), "data/question.json"))
	if err != nil {
		panic(err.Error())
	}
	temp, err := StructAnswer(body)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("总数:%#v\n", temp.Page.Totals)
		fmt.Printf("%#v\n", temp.Page.IsEnd)
		fmt.Printf("%#v\n", temp.Page.IsStart)
		fmt.Printf("下一个:%#v\n", temp.Page.NextUrl)
		fmt.Printf("回答：%#v\n", temp.Data[0].Content)
	}

}

func TestOutputHtml(t *testing.T) {
	body, err := util.ReadfromFile(filepath.Join(util.CurDir(), "data/question.json"))
	if err != nil {
		panic(err.Error())
	}
	temp, err := StructAnswer(body)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		answer := temp.Data[0]
		qid, title, aid, who, html := OutputHtml(answer)
		fmt.Println(qid, aid, who, title)
		util.SaveToFile(filepath.Join(util.CurDir(), "data/question.html"), []byte(html))
	}

}

func TestCatchOneAnswer(t *testing.T) {
	e := SetCookie("/home/jinhan/cookie.txt")
	if e != nil {
		fmt.Println(e.Error())
	}
	//https://www.zhihu.com/question/22030591/answer/22392511
	//51829486/answer/197012416
	Qid := "51829486"
	Aid := "197012416"
	data, e := CatchOneAnswer(Qid, Aid)
	if e != nil {
		fmt.Println(e.Error())
	} else {
		fmt.Println(string(data))
		util.SaveToFile(filepath.Join(util.CurDir(), "data/onequestion.html"), data)
	}
}

func TestParseOneAnswer(t *testing.T) {
	a, _ := util.ReadfromFile(filepath.Join(util.CurDir(), "data/onequestion.html"))
	r := ParseOneAnswer(a)
	fmt.Printf("%#v", r)
}

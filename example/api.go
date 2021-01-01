package main

import (
	"fmt"
	"strings"

	zhihu "github.com/hunterhug/huhu/src"
)

// API使用说明
func main() {
	//  1. 设置爬虫暂停时间，可选
	zhihu.SetWaitTime(1)

	// 2. 调试模式设置为debug，可选
	zhihu.SetLogLevel("info")

	// 3. 需先传入cookie，必须
	e := zhihu.SetCookie("./cookie.txt")
	if e != nil {
		panic(e.Error())
	}

	// 4.构建问题，url差页数
	q := zhihu.Question("28467579")
	fmt.Println(q)

	// 5.抓取问题回答，按页数，传入页数是为了补齐url，策略是循环抓，直到抓不到可认为页数已完
	page := 1
	limit := 20
	body, e := zhihu.CatchAnswer(q, limit, page)
	if e != nil {
		fmt.Println(e.Error())
		return
	}
	if strings.Contains(string(body), "error") { //可能cookie失效
		b, _ := zhihu.JsonBack(body)
		fmt.Println(string(b))
	}

	// 6.结构化回答
	answers, e := zhihu.StructAnswer(body)
	if e != nil {
		fmt.Println(e.Error())
	} else {
		// 就不打出来了
		//fmt.Printf("%#v\n", answers.Page)
		//fmt.Printf("%#v\n", answers.Data)
	}

	for _, a := range answers.Data {
		fmt.Printf("a:%#v\n", a)
	}
	// 7. 选择OutputHtml不要防盗链，因为回答输出的html经过了处理，所以我们进行过滤出好东西
	zhihu.SetPublishToWeb(false)
	qid, aid, t, who, html := zhihu.OutputHtml(answers.Data[0])
	fmt.Println(qid)
	fmt.Println(aid)
	fmt.Println(t)
	fmt.Println(who)

	// 8. 抓图片, 并发进程数10
	zhihu.SetSavePicture(false, 10)
	zhihu.SavePicture("test", []byte(html))

	// 9. 抓集合，第2页
	b, e := zhihu.CatchCoolection(78172986, 2)
	if e != nil {
		fmt.Println(e.Error())
	} else {
		// 解析集合
		fmt.Printf("%#v", zhihu.ParseCollection(b))
	}
}

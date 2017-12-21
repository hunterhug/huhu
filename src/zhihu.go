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

// 知乎系列爬虫
package src

import (
	"fmt"
	"strings"

	"github.com/hunterhug/marmot/miner"
	"github.com/hunterhug/parrot/util"
) // 第一步：引入库

//"fmt"

var (
	// 各种链接
	QuestionUrl   = "https://www.zhihu.com/question/%d"
	PeopleUrl     = "https://www.zhihu.com/people/%s"
	AnswerUrl     = "https://www.zhihu.com/question/%d/answer/%d"
	CollectionUrl = "https://www.zhihu.com/collection/%d?page=%d"

	// 一只小爬虫
	Baba *miner.Worker

	// 图片临时爬虫
	PTemp *miner.Worker

	// 知乎防盗链，要加一个js
	PublishToWeb = false
	// 抓取图片？
	CatchP = false
	Debug  = "info"
)

func init() {
	// 第一步：可选设置全局
	miner.SetLogLevel(Debug)   // 设置全局爬虫日志，可不设置，设置debug可打印出http请求轨迹
	miner.SetGlobalTimeout(60) // 爬虫超时时间，可不设置，默认超长时间

	// 第二步： 新建一个爬虫对象，nil表示不使用代理IP，可选代理
	spiders, err := miner.New(nil)

	if err != nil {
		panic(err)
	}
	Baba = spiders
	Baba.SetUa(miner.RandomUa())
	Baba.SetWaitTime(1)

	PTemp, _ = miner.New(nil)
	PTemp.SetUa(miner.RandomUa())
}

// 设置爬虫调试日志级别，开发可用:debug,info
func SetLogLevel(level string) {
	miner.SetLogLevel(level)

}

// 设置爬虫暂停时间
func SetWaitTime(w int) {
	Baba.SetWaitTime(w)
}

// 输出HTML选择防盗链方式
func SetPublishToWeb(put bool) {
	PublishToWeb = put
}

// 登录，验证码突破不了，请采用SetCookie
func Login(email, password string) ([]byte, error) {
	if strings.Contains(email, "@") {
		Baba.SetUrl("https://www.zhihu.com/login/email").SetRefer("https://www.zhihu.com/").SetUa(miner.RandomUa())
		Baba.SetFormParm("email", email).SetFormParm("password", password)
	} else {
		Baba.SetUrl("https://www.zhihu.com/login/phone_num").SetRefer("https://www.zhihu.com/").SetUa(miner.RandomUa())
		Baba.SetFormParm("phone_num", email).SetFormParm("password", password)
	}
	body, err := Baba.Post()

	// 清除Post的数据，方便下次使用
	Baba.Clear()

	if err != nil {
		return []byte("网路错误..."), err
	}
	return JsonBack(body)
}

// 设置cookie，需传入文件位置，文件中放cookie
func SetCookie(file string) error {
	haha, err := util.ReadfromFile(file)
	if err != nil {
		return err
	}
	cookie := string(haha)
	cookie = strings.Replace(cookie, " ", "", -1)
	cookie = strings.Replace(cookie, "\n", "", -1)
	cookie = strings.Replace(cookie, "\r", "", -1)
	Baba.SetHeaderParm("Cookie", strings.TrimSpace(cookie))
	return nil
}

// 谨慎使用,关注某人
func FollowWho(who string) ([]byte, error) {
	Baba.SetUrl(fmt.Sprintf("https://www.zhihu.com/api/v4/members/%s/followers", who))
	return Baba.Post()
}

func Follow(who string) {
	//Baba.SetUrl(fmt.Sprintf("https://www.zhihu.com/api/v4/members/%s/followers", who))
	//Baba.Post()
}

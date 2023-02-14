package core

import (
	"fmt"
	"strings"

	"github.com/hunterhug/huhu/util"
	"github.com/hunterhug/marmot/miner"
)

var (
	QuestionUrl      = "https://www.zhihu.com/question/%d"
	PeopleUrl        = "https://www.zhihu.com/people/%s"
	AnswerUrl        = "https://www.zhihu.com/question/%d/answer/%d"
	CollectionUrl    = "https://www.zhihu.com/collection/%d?page=%d"
	CollectionUrlNew = "https://www.zhihu.com/api/v4/collections/%d/items?limit=400&offset=%d"

	// Baba 一只小爬虫
	Baba *miner.Worker

	// PTemp 图片临时爬虫
	PTemp *miner.Worker

	// PublishToWeb 知乎防盗链，要加一个js
	PublishToWeb = false

	// CatchP 抓取图片？
	CatchP   = false
	LogLevel = miner.INFO
)

func init() {
	// 第一步：可选设置全局
	miner.SetLogLevel(LogLevel) // 设置全局爬虫日志，可不设置，设置debug可打印出http请求轨迹
	miner.SetGlobalTimeout(60)  // 爬虫超时时间，可不设置，默认超长时间

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

// SetLogLevel 设置爬虫调试日志级别
func SetLogLevel(level miner.Level) {
	miner.SetLogLevel(level)
}

// SetWaitTime 设置爬虫暂停时间
func SetWaitTime(w int) {
	Baba.SetWaitTime(w)
}

// SetPublishToWeb 输出HTML选择防盗链方式
func SetPublishToWeb(put bool) {
	PublishToWeb = put
}

// Login 登录，验证码突破不了，请采用SetCookie
func Login(email, password string) ([]byte, error) {
	if strings.Contains(email, "@") {
		Baba.SetUrl("https://www.zhihu.com/login/email").SetRefer("https://www.zhihu.com/").SetUa(miner.RandomUa())
		Baba.SetFormParam("email", email).SetFormParam("password", password)
	} else {
		Baba.SetUrl("https://www.zhihu.com/login/phone_num").SetRefer("https://www.zhihu.com/").SetUa(miner.RandomUa())
		Baba.SetFormParam("phone_num", email).SetFormParam("password", password)
	}
	body, err := Baba.Post()

	// 清除Post的数据，方便下次使用
	Baba.Clear()

	if err != nil {
		return []byte("网路错误..."), err
	}
	return JsonBack(body)
}

func SetCookie(file string) error {
	haha, err := util.ReadFromFile(file)
	if err != nil {
		return err
	}
	cookie := string(haha)
	cookie = strings.Replace(cookie, " ", "", -1)
	cookie = strings.Replace(cookie, "\n", "", -1)
	cookie = strings.Replace(cookie, "\r", "", -1)
	Baba.SetHeaderParam("Cookie", strings.TrimSpace(cookie))
	Baba.SetUrl("https://www.zhihu.com").Get()
	return nil
}

// FollowWho 谨慎使用,关注某人
func FollowWho(who string) ([]byte, error) {
	Baba.SetUrl(fmt.Sprintf("https://www.zhihu.com/api/v4/members/%s/followers", who))
	return Baba.Post()
}

func Follow(who string) {
	//Baba.SetUrl(fmt.Sprintf("https://www.zhihu.com/api/v4/members/%s/followers", who))
	//Baba.Post()
}

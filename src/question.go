package src

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/hunterhug/marmot/expert"
)

var (
	// 问题链接
	Qurl = "https://www.zhihu.com/api/v4/questions/%s/answers?"

	// 单答案链接
	Q1url     = "https://www.zhihu.com/question/%s/answer/%s"
	Qurlquery = "sort_by=default&include=data[*].%s"
	// 各种参数，问题获取到的JSON字段意思
	Qurlparm = []string{
		"is_normal",           // 是否正常
		"is_collapsed",        // 是否折叠
		"collapse_reason",     // 折叠理由
		"is_sticky",           // 无
		"collapsed_by",        // 5
		"suggest_edit",        // 建议
		"comment_count",       //评论数
		"can_comment",         // 能否评论
		"content",             //内容 重要
		"editable_content",    //5
		"voteup_count",        // 投票数?
		"reshipment_settings", //?
		"comment_permission",  // 可否评论
		"mark_infos",          //5
		"created_time",
		"updated_time",
		"relationship.is_authorized,is_author,voting,is_thanked,is_nothelp", // 关系？
		"upvoted_followees;data[*].author.follower_count,badge[?(type=best_answerer)].topics",
	}
)

type Answer struct {
	Page PageInfo   `json:"paging"`
	Data []DataInfo `json:"data"`
}

type PageInfo struct {
	IsEnd   bool   `json:"is_end"`
	Totals  int    `json:"totals"`
	PreUrl  string `json:"previous"`
	IsStart bool   `json:"is_start"`
	NextUrl string `json:"next"`
}

type DataInfo struct {
	Excerpt    string       `json:"excerpt"`
	Author     AuthorInfo   `json:"author"`
	Question   QuestionInfo `json:"question"`
	Content    string       `json:"content"`
	Aid        int          `json:"id"`
	CreateTime int          `json:"created_time"`
	UpdateTime int          `json:"updated_time"`
}

type AuthorInfo struct {
	About    string `json:"headline"`
	UrlToken string `json:"url_token"`
	Name     string `json:"name"`
	Sex      int    `json:"gender"`
	Image    string `json:"avatar_url_template"`
}
type QuestionInfo struct {
	CreateTime int    `json:"created`
	Title      string `json:"title"`
	UpdateTime int    `json:"updated_time"`
	Qid        int    `json:"id"`
}

var tempParm = fmt.Sprintf(Qurlquery, strings.Join(Qurlparm, ","))

// 构造问题链接，返回url,你可以通过Qurlparm拼出另外一个url
func Question(id string) string {
	return fmt.Sprintf(Qurl, id) + tempParm + "&limit=%d&offset=%d"
}

// 抓答案，需传入限制和页数,每次最多抓20个答案
func CatchAnswer(url string, limit, page int) ([]byte, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 1
	}
	if limit > 20 {
		limit = 20
	}
	uurl := fmt.Sprintf(url, limit, (page-1)*limit)

	//fmt.Println(uurl)
	Baba.SetUrl(uurl)

	body, err := Baba.Get()
	if err != nil {

	} else {
		if strings.Contains(string(body), "AuthenticationInvalid") {
			data, _ := JsonBack(body)
			return data, errors.New("CookiePASS")
		}
	}
	return body, err
}

// 抓单个答案，需传入问题ID和答案ID 鸡肋功能，弃用！
func CatchOneAnswer(Qid, Aid string) ([]byte, error) {
	Baba.SetUrl(fmt.Sprintf(Q1url, Qid, Aid))
	body, err := Baba.Get()
	if err != nil {

	} else {
		if strings.Contains(string(body), "AuthenticationInvalid") {
			data, _ := JsonBack(body)
			return data, errors.New("CookiePASS")
		}
	}
	return body, err
}

// 解析单个答案，待完善
func ParseOneAnswer(data []byte) map[string]string {
	a, e := expert.QueryBytes(data)
	if e != nil {
		println(e.Error())
	}
	b := a.Find(".CopyrightRichText-richText")
	c, e := b.Html()
	if e != nil {
		println(e.Error())
	}
	d := map[string]string{"text": c}
	return d
}

// 结构化回答，返回一个结构体
func StructAnswer(body []byte) (*Answer, error) {
	temp := new(Answer)
	err := json.Unmarshal(body, temp)
	return temp, err
}

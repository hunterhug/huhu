package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hunterhug/marmot/expert"
	"regexp"
	"strings"
)

var PeopleAnswer = "https://www.zhihu.com/people/%s/answers?page=%d"

// PeopleAnswerSS 单个用户回答
type PeopleAnswerSS struct {
	Entities    OutEntities `json:"entities"`
	CurrentUser string      `json:"currentUser"`
}

type OutEntities struct {
	Users   map[string]UserInfo  `json:"users"`
	Answers map[string]OutAnswer `json:"answers"`
	// todo
}

type UserInfo struct {
	FollowingCount    int64               `json:"followingCount"`    // 关注了
	FollowerCount     int64               `json:"followerCount"`     // 关注者
	UrlToken          string              `json:"urlToken"`          // 唯一标识
	Name              string              `json:"name"`              // 名字
	AnswerCount       int64               `json:"answerCount"`       // 回答问题数
	QuestionCount     int                 `json:"questionCount"`     // 提问数
	Gender            int                 `json:"gender"`            // 性别
	ThankedCount      int64               `json:"thankedCount"`      // 被感谢
	FavoritedCount    int64               `json:"favoritedCount"`    // 被收藏
	VoteupCount       int64               `json:"voteupCount"`       // 被赞成
	Headline          string              `json:"headline"`          // 简介小
	Description       string              `json:"description"`       // 简介大
	AvatarUrlTemplate string              `json:"avatarUrlTemplate"` // 头像  xll
	Business          UserInfoBusiness    `json:"business"`          // 行业
	Locations         []UserInfoLocations `json:"locations"`         // 地址
}

type UserInfoBusiness struct {
	Name string `json:"name"`
}

type UserInfoLocations struct {
	Name string `json:"name"`
}

type OutAnswer struct {
	Content     string `json:"content"` // 内容
	CreatedTime int    `json:"createdTime"`
	UpdateTime  int    `json:"updatedTime"`
	//Author       map[string]interface{} `json:"author"`
	Question     OutAnswerQuestionInfo `json:"question"`
	VoteupCount  int                   `json:"voteupCount"`  // 投票数: 赞同
	CommentCount int                   `json:"commentCount"` // 评论数
	Url          string                `json:"url"`          // 网址
}

type OutAnswerQuestionInfo struct {
	Created     int64  `json:"created"`
	Url         string `json:"url"`
	Title       string `json:"title"`
	Id          int64  `json:"id"`
	Typedudu    string `json:"type"`
	UpdatedTime int64  `json:"updatedTime"`
}

// CatchPeopleAnswer 获取一个人的一页回答, who为用户标志, 如:da-xiong-nu-da-xiong-nu page为页数
func CatchPeopleAnswer(who string, page int) ([]byte, error) {
	url := fmt.Sprintf(PeopleAnswer, who, page)
	Baba.SetUrl(url)
	b, e := Baba.Get()
	if strings.Contains(string(b), "你似乎来到了没有知识存在的荒原") {
		e = errors.New("not exist this page")
	}

	if strings.Contains(string(b), "EmptyState-image") && strings.Contains(string(b), "还没有回答") {
		e = errors.New("empty this page")
	}

	return b, e
}

// ParsePeopleAnswer 解析获取的回答, 返回的是一个结构体
func ParsePeopleAnswer(data []byte) PeopleAnswerSS {
	r := PeopleAnswerSS{}
	doc, _ := expert.QueryBytes(data)
	text, ok := doc.Find("div#data").Attr("data-state")
	if ok {
		//fmt.Println(string(text))
		e := json.Unmarshal([]byte(text), &r)
		if e != nil {
			fmt.Println(e.Error())
		}
	}
	return r
}

// Todo
// 获取一个人的所有回答, 由以上函数封装(内存占用由该用户回答数决定), 返回带有页数的map
func CatchPeopleAllAnswer(who string) map[int]PeopleAnswerSS {
	//page := 1
	//d, e := CatchPeopleAnswer(who, page)
	//if e != nil {
	//	fmt.Println(e.Error())
	//	return
	//} else {
	//	uinfo := ParsePeopleAnswer(d)
	//}
	return nil
}

func ReplacePeopleOneAnswerOuput(s string, markdown bool) string {
	r1, _ := regexp.Compile("<noscript>.*?</noscript>")
	s = r1.ReplaceAllString(s, "")
	r, _ := regexp.Compile(`src="(.*?)"`)
	s = r.ReplaceAllString(s, "")

	s = strings.Replace(s, "data-actualsrc", "src", -1)
	s = strings.Replace(s, "data-original", "src", -1)

	if markdown {
		//s = html2md.Convert(s)
	}
	return s
}

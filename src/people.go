package src

import (
	"encoding/json"
	"fmt"
	"github.com/hunterhug/marmot/expert"
	"strings"
	"errors"
	"regexp"
	"github.com/lunny/html2md"
)

// https://www.zhihu.com/people/da-xiong-nu-da-xiong-nu/answers
var PeopleAnswer = "https://www.zhihu.com/people/%s/answers?page=%d"

// 单个用户回答
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

/*
"244181927": {
        "suggestEdit": {
          "status": false,
          "reason": "",
          "title": "",
          "url": "",
          "unnormalDetails": {},
          "tip": ""
        },
        "relationship": {
          "upvotedFollowees": [],
          "isAuthor": false,
          "isNothelp": false,
          "isAuthorized": false,
          "voting": 0,
          "isThanked": false
        },
        "markInfos": [],
        "excerpt": "这个要看体质吧，我187的这个身高就是在15岁稳定的没有再长，我爸189我妈177也算稳定遗传 但是我认识一个191的妹子，他爸172他妈178他上初中那年才150，初中毕业就奔180了，真的可以说是看着他长大的，高中毕业后188，然后一直稳定到20岁，之后他开始健身，…",
        "annotationAction": [],
        "adminClosedComment": false,
        "collapsedBy": "nobody",
        "createdTime": 1507952175,
        "id": 244181927,
        "voteupCount": 13,
        "collapseReason": "",
        "isCollapsed": false,
        "author": {
          "avatarUrlTemplate": "https://pic3.zhimg.com/50/v2-55112c6c5de90df0ef30207f3b8fbf16_hd.jpg",
          "name": "萌妹咂",
          "headline": "一个17岁G的187CM的小可爱(划掉小可爱)一个不知名小网站签约作者",
          "type": "people",
          "userType": "people",
          "urlToken": "da-xiong-nu-da-xiong-nu",
          "isAdvertiser": false,
          "avatarUrl": "https://pic3.zhimg.com/50/v2-55112c6c5de90df0ef30207f3b8fbf16_hd.jpg",
          "url": "http://www.zhihu.com/people/decc8cea88a3063dc1a2cd2f04a0fcfa",
          "gender": 0,
          "badge": [],
          "id": "decc8cea88a3063dc1a2cd2f04a0fcfa",
          "isOrg": false
        },
        "url": "http://www.zhihu.com/answers/244181927",
        "commentPermission": "all",
        "canComment": {
          "status": true,
          "reason": ""
        },
        "question": {
          "questionType": "normal",
          "created": 1507708436,
          "url": "http://www.zhihu.com/questions/66496003",
          "title": "16岁还能能长高吗？",
          "type": "question",
          "id": 66496003,
          "updatedTime": 1507708436
        },
        "updatedTime": 1507952176,
        "content": "这个要看体质吧，我187的这个身高就是在15岁稳定的没有再长，我爸189我妈177也算稳定遗传<br>但是我认识一个191的妹子，他爸172他妈178他上初中那年才150，初中毕业就奔180了，真的可以说是看着他长大的，高中毕业后188，然后一直稳定到20岁，之后他开始健身，游泳什么的，现在22已经191.3了，所以身高这种东西真的说不准",
        "commentCount": 28,
        "extras": "",
        "reshipmentSettings": "allowed",
        "rewardInfo": {
          "rewardMemberCount": 0,
          "isRewardable": false,
          "rewardTotalMoney": 0,
          "canOpenReward": false,
          "tagline": ""
        },
        "isCopyable": true,
        "type": "answer",
        "thumbnail": "",
        "isNormal": true
      }
*/
type OutAnswer struct {
	Content     string `json:"content"` // 内容
	CreatedTime int    `json:"createdTime"`
	UpdateTime  int    `json:"updatedTime"`
	//Author       map[string]interface{} `json:"author"`
	Question     OutAnswerQuestionInfo `json:"question"`
	VoteupCount  int                   `json:"voteup_count"`  // 投票数: 赞同
	CommentCount int                   `json:"comment_count"` // 评论数
	Url          string                `json:"url"`           // 网址
}

type OutAnswerQuestionInfo struct {
	Created     int64  `json:"created"`
	Url         string `json:"url"`
	Title       string `json:"title"`
	Id          int64  `json:"id"`
	Typedudu    string `json:"type"`
	UpdatedTime int64  `json:"updatedTime"`
}

// 获取一个人的一页回答, who为用户标志, 如:da-xiong-nu-da-xiong-nu page为页数
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

// 解析获取的回答, 返回的是一个结构体
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
		s = html2md.Convert(s)
	}
	return s
}

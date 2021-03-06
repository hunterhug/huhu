package src

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var (
	// 谁关注
	// 粉丝
	followersrurl = "https://www.zhihu.com/api/v4/members/%s/followers?"

	// 关注谁
	// 偶像
	followeesurl = "https://www.zhihu.com/api/v4/members/%s/followees?include="
	fparm        = "data[*].answer_count,articles_count,gender,follower_count,is_followed,is_following,badge[?(type=best_answerer)].topics"
)

type FollowData struct {
	Page PageInfo           `json:"paging"`
	Data []FollowerDataInfo `json:"data"`
}


type FollowerDataInfo struct {
	IsFollowed        bool   `json:"is_followed"`
	AvatarUrlTemplate string `json:"avatar_url_template"`
	UserType          string `json:"user_type"`
	AnswerCount       int    `json:"answer_count"`
	IsFollowing       bool   `json:"is_following"`
	Headline          string `json:"headline"`
	UrlToken          string `json:"url_token"`
	Id                string `json:"id"`
	ArticlesCount     int    `json:"articles_count"`
	Type              string `json:"type"`
	Name              string `json:"name"`
	Url               string `json:"url"`
	Gender            int    `json:"gender"`
	IsAdvertiser      bool   `json:"is_advertiser"`
	AvatarUrl         string `json:"avatar_url"`
	IsOrg             bool   `json:"is_org"`
	FollowerCount     int    `json:"follower_count"`
}

func followees(token string, limit, offset int) string {
	return fmt.Sprintf(followeesurl, token) + fparm + fmt.Sprintf("&limit=%d&offset=%d", limit, offset)
}

func followers(token string, limit, offset int) string {
	return fmt.Sprintf(followersrurl, token) + fparm + fmt.Sprintf("&limit=%d&offset=%d", limit, offset)
}

// 抓取用户：fensi 抓取你的粉丝，否则，抓取你的偶像，token为用户：如https://www.zhihu.com/people/hunterhug中的hunterhug,limit限制最多20条
func CatchUser(fensi bool, token string, limit, offset int) ([]byte, error) {
	if limit < 0 {
		limit = -limit
	}
	if limit >= 20 {
		limit = 20
	}
	url := ""
	if fensi {
		url = followers(token, limit, offset)
	} else {
		url = followees(token, limit, offset)
	}
	Baba.SetUrl(url)
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

// 解析用户
func ParseUser(data []byte) FollowData {
	var r FollowData
	json.Unmarshal(data, &r)
	return r
}

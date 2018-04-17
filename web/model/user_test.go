package model

import (
	"testing"
	"fmt"
	"github.com/hunterhug/huhu/src"
	"strings"
)

func TestInsertOrUpdateUser(t *testing.T) {
	who := "hunterhug"
	page := 1
	d, e := src.CatchPeopleAnswer(who, page)
	if e != nil {
		fmt.Println(e.Error())
	} else {
		uinfo := src.ParsePeopleAnswer(d)
		us := uinfo.Entities.Users
		if len(us) > 0 {
			for _, v := range us {

				u := User{}
				u.AvatarUrl = strings.Replace(v.AvatarUrlTemplate, "{size}", "xll", -1)
				u.AnswerCount = v.AnswerCount
				u.Business = v.Business.Name
				u.Description = v.Description
				u.FavoritedCount = v.FavoritedCount
				u.FollowerCount = v.FollowerCount
				u.FollowingCount = v.FollowingCount
				u.Gender = v.Gender
				u.Headline = v.Headline
				u.Name = v.Name
				u.ThankedCount = v.ThankedCount
				u.QuestionCount = v.QuestionCount
				u.VoteupCount = v.VoteupCount
				u.AnswerCount = v.AnswerCount
				u.Id = v.UrlToken
				location := []string{}
				for _, l := range v.Locations {
					location = append(location, l.Name)
				}
				u.Locations = strings.Join(location, ",")

				InsertOrUpdateUser(&u)
			}
		}
	}
}

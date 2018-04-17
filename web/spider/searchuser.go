package spider

import (
	"github.com/hunterhug/huhu/src"
	"fmt"
	"strings"
	"github.com/hunterhug/huhu/web/model"
	"sync"
	"github.com/hunterhug/marmot/miner"
	"github.com/hunterhug/parrot/util"
)

var GlobalUserHeap = Heap{
	DoingUser: map[string]bool{},
	DoneUser:  map[string]bool{}}

type Heap struct {
	DoneUser  map[string]bool
	DoingUser map[string]bool
	Num       int64
	//Root     *DiyHeap
	mux sync.RWMutex
}

//type DiyHeap struct {
//	Head  *DiyToken
//	Tail  *DiyToken
//	Total int64
//}
//
//type DiyToken struct {
//	Token string
//	Next  *DiyToken
//}

func (h *Heap) Add(userId string) bool {
	h.mux.Lock()
	defer h.mux.Unlock()
	if exist, _ := h.DoneUser[userId]; exist {
		return true
	} else {
		fmt.Println(userId + " adding")
		if ee, _ := h.DoingUser[userId]; !ee {
			h.Num = h.Num + 1
			h.DoingUser[userId] = true
		}

	}
	return false
}

func (h *Heap) Delete(userId string) {
	h.mux.Lock()
	defer h.mux.Unlock()
	h.DoneUser[userId] = true
	if ee, _ := h.DoingUser[userId]; ee {
		fmt.Println(userId + " deleteing")
		delete(h.DoingUser, userId)
		if h.Num > 0 {
			h.Num = h.Num - 1
		}
	}
}

func (h *Heap) GetNum() int64 {
	h.mux.Lock()
	defer h.mux.Unlock()
	return h.Num
}

// 深度搜索
func StartSpider(who string, layer int) {
	fmt.Printf("==layer:%v, listnum:%v==\n", layer, GlobalUserHeap.GetNum())
	ll := UserList(who)
	for _, k := range ll {
		fmt.Printf("%v->%v\n", who, k)
		SearchUser(k)
	}
	for _, k := range ll {
		GlobalUserHeap.Delete(k)
		StartSpider(k, layer+1)
	}
}

func UserList(who string) []string {
	r := []string{}
	offset := 0
	for {
		if len(r) > 200 {
			break
		}
		a, e := src.CatchUser(false, who, 20, offset)
		if e != nil {
			fmt.Println(e.Error())
		} else {
			rr := src.ParseUser(a)
			if len(rr.Data) == 0 {
				break
			}
			for _, v := range rr.Data {
				GlobalUserHeap.Add(v.UrlToken)
				r = append(r, v.UrlToken)
			}
			offset = offset + 20
		}
	}

	return r
}

var CurDir string

func init() {
	CurDir, _ = util.GetBinaryCurrentPath()
	util.MakeDir(CurDir + "/pic")

}

func SearchUser(who string) {
	page := 1
	d, e := src.CatchPeopleAnswer(who, page)
	if e != nil {
		fmt.Println(e.Error())
	} else {
		uinfo := src.ParsePeopleAnswer(d)
		us := uinfo.Entities.Users
		if len(us) > 0 {
			for _, v := range us {
				if uinfo.CurrentUser == v.UrlToken {
					continue
				}
				u := model.User{}
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
				model.InsertOrUpdateUser(&u)

				pp, _ := miner.New(nil)
				braw, e := pp.SetUa(miner.RandomUa()).SetUrl(u.AvatarUrl).Get()
				if e == nil {
					fdir := CurDir + "/pic/" + u.Id + ".jpg"
					if util.FileExist(fdir) {
						return
					}
					err := util.SaveToFile(fdir, braw)
					if err == nil {
						fmt.Println(fdir)
					} else {
						fmt.Println(fdir + ":" + err.Error())
					}
				}
			}
		}
	}
}

package model

import (
	"fmt"
	"github.com/hunterhug/parrot/store/mydb"
	"time"
)

type User struct {
	FollowingCount int64                // 关注了
	FollowerCount  int64                // 关注者
	Id             string `xorm:"pk"`   // 唯一标识
	Name           string               // 名字
	AnswerCount    int64                // 回答问题数
	QuestionCount  int                  // 提问数
	Gender         int                  // 性别
	ThankedCount   int64                // 被感谢
	FavoritedCount int64                // 被收藏
	VoteupCount    int64                // 被赞成
	Headline       string `xorm:"text"` // 简介小
	Description    string `xorm:"text"` // 简介大
	AvatarUrl      string               // 头像  xll
	LocalAvatarUrl string               // 本地头像
	Business       string               // 行业
	Locations      string               // 地址
	DbCreateTime   int64
	DbUpdateTime   int64
}

var GlobalDBClient *mydb.MyDb

func init() {
	config := mydb.MyDbConfig{
		DriverName: mydb.MYSQL,
		DbConfig: mydb.DbConfig{
			Host: "127.0.0.1",
			User: "root",
			Pass: "6833066",
			Name: "zhihu",
		},
		MaxOpenConns: 10,
		MaxIdleConns: 5,
		//Debug:        true,
	}

	db, err := mydb.NewDb(config)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		err := db.Ping()
		if err != nil {
			fmt.Println("d" + err.Error())
		}
		GlobalDBClient = db
	}

	//GlobalDBClient.DropTables(User{})
	if ok, _ := GlobalDBClient.IsTableExist(User{}); !ok {
		e := GlobalDBClient.CreateTables(User{})
		if e != nil {
			fmt.Println(e.Error())
		}
	}
}

func InsertOrUpdateUser(u *User) {
	if u == nil {
		return
	}
	if ok, _ := GlobalDBClient.Client.ID(u.Id).Get(&User{}); ok {
		u.DbUpdateTime = time.Now().UTC().Unix()
		_, err := GlobalDBClient.Client.ID(u.Id).Update(u)
		if err != nil {
			fmt.Println(err.Error())
		}
		return
	}
	u.DbCreateTime = time.Now().UTC().Unix()
	_, err := GlobalDBClient.InsertOne(u)
	if err != nil {
		fmt.Println(err.Error())
	}
}

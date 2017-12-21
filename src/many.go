package src

import (
	"fmt"
	"os"
	"strings"

	"github.com/hunterhug/parrot/util"
)

// collect id and boss way catch? Conitue catch?
func ManyAlone(id string, Boss bool, Conitue bool, SaveP bool, Limit int) {
	SetSavePicture(SaveP, 10)
	collectids := id
	if id == "" {
		//78172986
		collectids = Input("萌萌：请输入集合ID:", "")
	}
	collectid, e := util.SI(collectids)
	if e != nil {
		fmt.Println("收藏夹ID错误")
		return
	}

	// 收藏夹已抓
	newcatch := true
	savecfff := collectids + ".txt" // 收藏夹续抓标志
	cxx, _ := util.ReadfromFile(savecfff)
	cxx1 := strings.Split(string(cxx), "\n")

	qids := map[string]string{}

	for _, v := range cxx1 {
		tttt := strings.Split(v, "-")
		if len(tttt) != 2 {
			continue
		}
		qids[tttt[0]] = v
	}

	if len(qids) > 0 {
		fmt.Printf("总计有%d个剩余问题:\n", len(qids))

		if Conitue {
			newcatch = false
		}
	}

	if newcatch {
		qids = CatchAllCollection(collectid)
		if len(qids) == 0 {
			fmt.Println("收藏夹下没有问题！")
			return
		}
		fmt.Printf("总计有%d个问题:\n", len(qids))
		s := []string{}
		for id, qa := range qids {
			fmt.Printf("ID:%s，Answer:%s\n", id, qa)
			temppp := fmt.Sprintf("%s-%s", id, strings.Replace(qa, ",", ".", -1))
			s = append(s, temppp)
			qids[id] = temppp
		}
		util.SaveToFile(savecfff, []byte(strings.Join(s, "\n")))
	}

	// 抓過的刪除掉
	txtmap := map[string]string{}
	for k, v := range qids {
		txtmap[k] = v
	}

	for id, _ := range qids {
		page := 1
		q := Question(id)
		//fmt.Println(q)

		// 第一个答案
		body, err := CatchAnswer(q, 1, page)
		fmt.Println("预抓取第一个回答！")
		if err != nil {
			fmt.Println("问题预抓取出错:" + id + "-" + err.Error())
			if strings.Contains(err.Error(), "CookiePASS") {
				a := []string{}
				for _, v := range txtmap {
					a = append(a, v)
				}
				util.SaveToFile(savecfff, []byte(strings.Join(a, "\n")))
				fmt.Println("cookie.txt失效!重新填写")
				os.Exit(1)
			}
			continue
		}

		temp, err := StructAnswer(body)
		if err != nil {
			fmt.Println("b" + err.Error())
			s, _ := JsonBack(body)
			fmt.Println(string(s))
			continue
		}
		if len(temp.Data) == 0 {
			delete(txtmap, id)
			fmt.Println("没有答案！")
			continue
		}

		fmt.Println("开始处理答案:" + temp.Data[0].Excerpt)
		qid, aid, title, who, html := OutputHtml(temp.Data[0])
		fmt.Println("哦，这个问题是:" + title)
		if util.FileExist(fmt.Sprintf("data/%d-%s.xx", qid, util.ValidFileName(title))) {
			fmt.Printf("已经存在：%s,跳过！\n", fmt.Sprintf("data/%d-%s.xx", qid, util.ValidFileName(title)))
			delete(txtmap, id)
			continue
		}

		filename := fmt.Sprintf("data/%d/%s-%d/%s-%d的回答.html", qid, who, aid, who, aid)
		util.MakeDirByFile(filename)
		if PublishToWeb {
			util.SaveToFile(fmt.Sprintf("data/%d/%s", qid, JsName), []byte(Js))
		}

		err = util.SaveToFile(filename, []byte(OneOutputHtml(html)))
		// html
		util.MakeDir(fmt.Sprintf("data/%d-html", qid))
		link := ""
		if page == 1 {
			link = fmt.Sprintf(`<a href="%d.html" style="float:right">Next下一页</a>`, page+1)
		} else {
			link = fmt.Sprintf(`<a href="%d.html" style="float:left">Pre上一页</a><a href="%d.html" style="float:right">Next下一页</a>`, page-1, page+1)
		}
		html = strings.Replace(html, "###link###", link, -1)

		if Boss {
			util.SaveToFile(fmt.Sprintf("data/%d-html/%d.html", qid, page), []byte(BossOutputHtml(qid, who, aid, html)))
		} else {
			util.SaveToFile(fmt.Sprintf("data/%d-html/%d.html", qid, page), []byte(html))
		}

		if err == nil {
			fmt.Println("保存答案成功:" + filename)
		} else {
			fmt.Println("保存答案失败:" + err.Error())
			continue
		}
		SavePicture(fmt.Sprintf("data/%d/%s-%d", qid, who, aid), []byte(html))

		for {
			if temp.Page.IsEnd {
				fmt.Println("回答已经结束！")
				break
			}
			//fmt.Println(temp.Page.NextUrl)
			if page+1 > Limit {
				fmt.Println("萌萌：答案超出个数了哦，哦耶~")
				break
			}
			body, err = CatchAnswer(q, 1, page+1)
			if err != nil {
				fmt.Println("抓取答案失败：" + id + "-" + err.Error())
				if strings.Contains(err.Error(), "CookiePASS") {
					a := []string{}
					for _, v := range txtmap {
						a = append(a, v)
					}
					util.SaveToFile(savecfff, []byte(strings.Join(a, "\n")))
					fmt.Println("cookie.txt失效!重新填写")
					os.Exit(1)
				}
				continue
			} else {
				page = page + 1
			}
			//util.SaveToFile("data/question.json", body)

			temp1, err := StructAnswer(body)
			if err != nil {
				fmt.Printf("%s:%s\n", err.Error(), string(body))
				break
			}
			if len(temp1.Data) == 0 {
				fmt.Println("没有答案！")
				s, _ := JsonBack(body)
				fmt.Println(string(s))
				break
			}

			// 成功后再来
			temp = temp1

			fmt.Println("开始处理答案:" + temp.Data[0].Excerpt)
			qid, aid, _, who, html := OutputHtml(temp.Data[0])
			filename := fmt.Sprintf("data/%d/%s-%d/%s-%d的回答.html", qid, who, aid, who, aid)
			util.MakeDirByFile(filename)
			err = util.SaveToFile(filename, []byte(OneOutputHtml(html)))
			// html
			util.MakeDir(fmt.Sprintf("data/%d-html", qid))
			link := ""
			if page == 1 {
				link = fmt.Sprintf(`<a href="%d.html" style="float:right">Next下一页</a>`, page+1)
			} else {
				link = fmt.Sprintf(`<a href="%d.html" style="float:left">Pre上一页</a><a href="%d.html" style="float:right">Next下一页</a>`, page-1, page+1)
			}
			html = strings.Replace(html, "###link###", link, -1)

			if Boss {
				util.SaveToFile(fmt.Sprintf("data/%d-html/%d.html", qid, page), []byte(BossOutputHtml(qid, who, aid, html)))
			} else {
				util.SaveToFile(fmt.Sprintf("data/%d-html/%d.html", qid, page), []byte(html))
			}

			if err == nil {
				fmt.Println("保存答案成功:" + filename)
			} else {
				fmt.Println("保存答案失败:", err.Error())
				continue
			}
			SavePicture(fmt.Sprintf("data/%d/%s-%d", qid, who, aid), []byte(html))
		}

		util.SaveToFile(fmt.Sprintf("data/%d-%s.xx", qid, util.ValidFileName(title)), []byte(""))
		delete(txtmap, id)

		// 每个问题都要保存一次,防止出错
		a := []string{}
		for _, v := range txtmap {
			a = append(a, v)
		}
		util.SaveToFile(savecfff, []byte(strings.Join(a, "\n")))
		fmt.Println("写入一次文件!")
	}
}

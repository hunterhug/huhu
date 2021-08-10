package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	zhihu "github.com/hunterhug/huhu/core"
	"github.com/hunterhug/huhu/util"
)

var Limit = 520
var Follow = false
var Boss = false

// 抓取一个问题的全部信息
// 每次只抓取一个答案
func help() {
	fmt.Println(`
	-----------------
	知乎问题信息小助手

	支持:
	1. 从收藏夹 https://www.zhihu.com/collection/78172986 批量获取很多问题答案
	2. 从问题 https://www.zhihu.com/question/28853910 批量获取一个问题很多答案

	请您按提示操作（Enter）！答案保存在 data 文件夹下，例如打开 28467579-html/1.html ！

	如果什么都没抓到请往 exe 同级目录 cookie.txt 增加 cookie。

	联系：https://github.com/hunterhug/huhu
	-----------------
	`)
}

var (
	datadir = flag.String("c", "cookie.txt", "cookie file position")
	x       = flag.Bool("x", false, "flag mode")
	mode    = flag.String("m", "2", "catch mode(1:question,2:collect,3)")
	id      = flag.String("i", "x", "id: must be num")
	l       = flag.Int("l", 300, "limit")
)

// 应该替换为本地照片！已经做了
func main() {
	if !flag.Parsed() {
		flag.Parse()
	}
	help()
	err := zhihu.SetCookie(*datadir)
	if err != nil {
		fmt.Println("请您一定要保证目录下cookie.txt存在哦：" + err.Error())
		time.Sleep(50 * time.Second)
		os.Exit(0)
	}

	if *x {
		switch *mode {
		case "2":
			// 文件夹命令行模式
			zhihu.ManyAlone(*id, true, true, true, *l)
		default:
			fmt.Println("mode error, now only support collect 2")
		}

		return
	}

	js := strings.ToLower(zhihu.Input(`萌萌：你有几种选项, 你的决定命运着图片链接是否被替换?

	选项:
	1. X(默认): 下载回答的同时图片会下载到电脑里，可以永久观看
	2. Y: 下载回答，但不下载图片，图片还是在知乎的服务器上
	3. Z: 将所有问题集合起来，在当前目录下生成一个 index.html

	请选择(可以直接回车):`, "x"))

	if strings.Contains(js, "y") ||
		strings.Contains(js, "x") ||
		strings.Contains(js, "z") {
	} else {
		js = "x"
	}
	if strings.Contains(js, "z") {
		zhihu.Tool()
		return
	}

	//if strings.Contains(js, "y") {
	//zhihu.PublishToWeb = true
	//zhihu.InitJs()
	//util.SaveToFile("data/"+zhihu.JsName, []byte(zhihu.Js))
	//}

	if strings.Contains(js, "x") {
		Boss = true
		zhihu.SetSavePicture(true, 10)
	} else {
		tu := strings.ToLower(zhihu.Input("萌萌：不试试抓取图片吗Y/N(默认N)", "n"))
		if strings.Contains(tu, "y") {
			zhihu.SetSavePicture(true, 10)
		}
	}
	choice := zhihu.Input("萌萌：从收藏夹获取回答按1，从问题获取回答按2(默认)", "2")
	for {
		ll := zhihu.Input("萌萌说亲爱的，因为每个问题的回答实在太多了，请限制获取的回答个数:30（默认)", "30")
		temp, errr := util.SI(ll)
		if errr != nil {
			fmt.Println("萌萌表示输入的应该是数字哦")
			continue
		}
		if temp <= 0 || temp > 500 {
			fmt.Println("萌萌表示不抓取，哼")
			continue
		}
		Limit = temp
		break
	}
	//ff := util.ToLower(zhihu.Input("酱~关注下答案中的小姐姐/小哥哥吧，默认N（Y/N）", "n"))
	ff := "n"
	if strings.Contains(ff, "y") {
		Follow = true
	}
	if choice == "1" {
		Many()
	} else {
		Base()
	}

}

func Base() {
	for {
		page := 1
		//28467579
		id := zhihu.Input("萌萌：请输入问题ID(默认28467579):", "28467579")
		q := zhihu.Question(id)
		//fmt.Println(q)

		// 第一个答案
		body, err := zhihu.CatchAnswer(q, 1, page)
		fmt.Println("预抓取第一个回答！")
		if err != nil {
			fmt.Println("a" + err.Error())
			continue
		}

		temp, err := zhihu.StructAnswer(body)
		if err != nil {
			fmt.Println("b" + err.Error())
			s, _ := zhihu.JsonBack(body)
			fmt.Println(string(s))
			continue
		}
		if len(temp.Data) == 0 {
			fmt.Println("没有答案！")
			break
		}

		fmt.Println("开始处理答案:" + temp.Data[0].Excerpt)
		qid, aid, title, who, html := zhihu.OutputHtml(temp.Data[0])
		fmt.Println("哦，这个问题是:" + title)
		if util.FileExist(fmt.Sprintf("data/%d-%s.xx", qid, util.ValidFileName(title))) {
			fmt.Printf("已经存在：%s,抓取请手动删除它！\n", fmt.Sprintf("data/%d-%s.xx", qid, util.ValidFileName(title)))
			break
		}
		filename := fmt.Sprintf("data/%d/%s-%d/%s-%d的回答.html", qid, who, aid, who, aid)
		util.MakeDirByFile(filename)
		if zhihu.PublishToWeb {
			util.SaveToFile(fmt.Sprintf("data/%d/%s", qid, zhihu.JsName), []byte(zhihu.Js))
		}
		util.SaveToFile(fmt.Sprintf("data/%d-%s.xx", qid, util.ValidFileName(title)), []byte(""))
		err = util.SaveToFile(filename, []byte(zhihu.OneOutputHtml(html)))

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
			util.SaveToFile(fmt.Sprintf("data/%d-html/%d.html", qid, page), []byte(zhihu.BossOutputHtml(qid, who, aid, html)))
		} else {
			util.SaveToFile(fmt.Sprintf("data/%d-html/%d.html", qid, page), []byte(html))
		}
		if Follow {
			zhihu.Follow(who)
		}
		if err == nil {
			fmt.Println("保存答案成功:" + filename)
		} else {
			fmt.Println("保存答案失败:" + err.Error())
			continue
		}
		zhihu.SavePicture(fmt.Sprintf("data/%d/%s-%d", qid, who, aid), []byte(html))

		all := util.ToLower(zhihu.Input("批量抓取答案，默认N(Y/N)", "N"))
		for {
			if temp.Page.IsEnd {
				fmt.Println("回答已经结束！")
				break
			}
			if strings.Contains(all, "n") {
				yes := util.ToLower(zhihu.Input("抓取下一个答案，默认Y(Y/N)", "Y"))
				if strings.Contains(yes, "n") {
					break
				}
			}
			//fmt.Println(temp.Page.NextUrl)
			if page+1 > Limit {
				fmt.Println("萌萌：答案超出个数了哦，哦耶~")
				break
			}
			body, err = zhihu.CatchAnswer(q, 1, page+1)
			if err != nil {
				fmt.Println("抓取答案失败：" + err.Error())
				continue
			} else {
				page = page + 1
			}
			//util.SaveToFile("data/question.json", body)

			temp1, err := zhihu.StructAnswer(body)
			if err != nil {
				fmt.Printf("%s:%s\n", err.Error(), string(body))
				break
			}
			if len(temp1.Data) == 0 {
				fmt.Println("没有答案！")
				s, _ := zhihu.JsonBack(body)
				fmt.Println(string(s))
				break
			}

			// 成功后再来
			temp = temp1

			fmt.Println("开始处理答案:" + temp.Data[0].Excerpt)
			qid, aid, _, who, html := zhihu.OutputHtml(temp.Data[0])
			filename := fmt.Sprintf("data/%d/%s-%d/%s-%d的回答.html", qid, who, aid, who, aid)
			util.MakeDirByFile(filename)
			err = util.SaveToFile(filename, []byte(zhihu.OneOutputHtml(html)))
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
				util.SaveToFile(fmt.Sprintf("data/%d-html/%d.html", qid, page), []byte(zhihu.BossOutputHtml(qid, who, aid, html)))
			} else {
				util.SaveToFile(fmt.Sprintf("data/%d-html/%d.html", qid, page), []byte(html))
			}
			if Follow {
				zhihu.Follow(who)
			}
			if err == nil {
				fmt.Println("保存答案成功:" + filename)
			} else {
				fmt.Println("保存答案失败:", err.Error())
				continue
			}
			zhihu.SavePicture(fmt.Sprintf("data/%d/%s-%d", qid, who, aid), []byte(html))
		}
	}
}

func Many() {
	for {
		//78172986
		collectids := zhihu.Input("萌萌：请输入集合ID:(默认78172986):", "78172986")
		collectid, e := util.SI(collectids)
		if e != nil {
			fmt.Println("收藏夹ID错误")
			continue
		}

		// 收藏夹已抓
		newcatch := true
		savecfff := collectids + ".txt" // 收藏夹续抓标志
		cxx, _ := util.ReadFromFile(savecfff)
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

			cxx2 := util.ToLower(zhihu.Input("上次这个收藏夹没抓完, 继续抓按Y, 默认Y(Y/N)?", "Y"))
			if strings.Contains(cxx2, "y") {
				newcatch = false
			}
		}
		god := util.ToLower(zhihu.Input("开启上帝模式吗(一路抓到底)，默认N(Y/N)?", "N"))
		skip := false
		if strings.Contains(god, "y") {
			skip = true
		}

		if newcatch {
			qids = zhihu.CatchAllCollection(collectid)
			if len(qids) == 0 {
				fmt.Println("收藏夹下没有问题！")
				continue
			}
			fmt.Printf("总计有%d个问题:\n", len(qids))
			var s []string
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
			q := zhihu.Question(id)
			//fmt.Println(q)

			// 第一个答案
			body, err := zhihu.CatchAnswer(q, 1, page)
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

			temp, err := zhihu.StructAnswer(body)
			if err != nil {
				fmt.Println("b" + err.Error())
				s, _ := zhihu.JsonBack(body)
				fmt.Println(string(s))
				continue
			}
			if len(temp.Data) == 0 {
				delete(txtmap, id)
				fmt.Println("没有答案！")
				continue
			}

			fmt.Println("开始处理答案:" + temp.Data[0].Excerpt)
			qid, aid, title, who, html := zhihu.OutputHtml(temp.Data[0])
			fmt.Println("哦，这个问题是:" + title)
			if util.FileExist(fmt.Sprintf("data/%d-%s.xx", qid, util.ValidFileName(title))) {
				fmt.Printf("已经存在：%s,跳过！\n", fmt.Sprintf("data/%d-%s.xx", qid, util.ValidFileName(title)))
				delete(txtmap, id)
				continue
			}

			if !skip {
				tiaotiao := util.ToLower(zhihu.Input("跳过这个问题吗，默认N(Y/N)?", "N"))
				if strings.Contains(tiaotiao, "y") {
					continue
				}
			}
			filename := fmt.Sprintf("data/%d/%s-%d/%s-%d的回答.html", qid, who, aid, who, aid)
			util.MakeDirByFile(filename)
			if zhihu.PublishToWeb {
				util.SaveToFile(fmt.Sprintf("data/%d/%s", qid, zhihu.JsName), []byte(zhihu.Js))
			}

			err = util.SaveToFile(filename, []byte(zhihu.OneOutputHtml(html)))
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
				util.SaveToFile(fmt.Sprintf("data/%d-html/%d.html", qid, page), []byte(zhihu.BossOutputHtml(qid, who, aid, html)))
			} else {
				util.SaveToFile(fmt.Sprintf("data/%d-html/%d.html", qid, page), []byte(html))
			}

			if Follow {
				zhihu.Follow(who)
			}
			if err == nil {
				fmt.Println("保存答案成功:" + filename)
			} else {
				fmt.Println("保存答案失败:" + err.Error())
				continue
			}
			zhihu.SavePicture(fmt.Sprintf("data/%d/%s-%d", qid, who, aid), []byte(html))

			all := "y"
			if !skip {
				all = util.ToLower(zhihu.Input("批量抓取这个问题的所有答案，默认N(Y/N)", "N"))
			}
			for {
				if temp.Page.IsEnd {
					fmt.Println("回答已经结束！")
					break
				}
				if strings.Contains(all, "n") {
					yes := util.ToLower(zhihu.Input("抓取下一个答案，默认Y(Y/N)", "Y"))
					if strings.Contains(yes, "n") {
						break
					}
				}
				//fmt.Println(temp.Page.NextUrl)
				if page+1 > Limit {
					fmt.Println("萌萌：答案超出个数了哦，哦耶~")
					break
				}
				body, err = zhihu.CatchAnswer(q, 1, page+1)
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

				temp1, err := zhihu.StructAnswer(body)
				if err != nil {
					fmt.Printf("%s:%s\n", err.Error(), string(body))
					break
				}
				if len(temp1.Data) == 0 {
					fmt.Println("没有答案！")
					s, _ := zhihu.JsonBack(body)
					fmt.Println(string(s))
					break
				}

				// 成功后再来
				temp = temp1

				fmt.Println("开始处理答案:" + temp.Data[0].Excerpt)
				qid, aid, _, who, html := zhihu.OutputHtml(temp.Data[0])
				filename := fmt.Sprintf("data/%d/%s-%d/%s-%d的回答.html", qid, who, aid, who, aid)
				util.MakeDirByFile(filename)
				err = util.SaveToFile(filename, []byte(zhihu.OneOutputHtml(html)))
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
					util.SaveToFile(fmt.Sprintf("data/%d-html/%d.html", qid, page), []byte(zhihu.BossOutputHtml(qid, who, aid, html)))
				} else {
					util.SaveToFile(fmt.Sprintf("data/%d-html/%d.html", qid, page), []byte(html))
				}

				if Follow {
					zhihu.Follow(who)
				}
				if err == nil {
					fmt.Println("保存答案成功:" + filename)
				} else {
					fmt.Println("保存答案失败:", err.Error())
					continue
				}
				zhihu.SavePicture(fmt.Sprintf("data/%d/%s-%d", qid, who, aid), []byte(html))
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
}

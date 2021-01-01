package src

import (
	"fmt"
	"strings"

	"github.com/hunterhug/marmot/expert"
	"github.com/hunterhug/marmot/miner"
	"github.com/hunterhug/marmot/util/goquery"
	"github.com/hunterhug/parrot/util"
)

var PSpiderPOOL []string

// 抓取图片前需要设置true
func SetSavePicture(catch bool, num int) {
	CatchP = catch
	for i := 0; i < num; i++ {
		sp, _ := miner.New(nil)
		sp.SetUa(miner.RandomUa())
		miner.Pool.Set(util.IS(i), sp)
		PSpiderPOOL = append(PSpiderPOOL, util.IS(i))
	}
}

// 抓取html中的图片，保存图片在dir下
func SavePicture(dir string, body []byte) {
	if !CatchP {
		return
	}
	util.MakeDir(dir)
	docm, err := expert.QueryBytes(body)
	if err != nil {
		fmt.Println(err.Error())
		return
	} else {
		//fmt.Println(string(content))
		piclist := []string{}
		docm.Find("img").Each(func(num int, node *goquery.Selection) {
			img, e := node.Attr("src")
			if e == false {
				img, e = node.Attr("data-src")
			}
			if e && img != "" {
				piclist = append(piclist, img)
			}
		})

		if len(piclist) == 0 {
			return
		}

		// 并发抓取图片
		//for _, img := range piclist {
		//	fmt.Println("并发抓取原始文件：" + img)
		//}

		// 高并发
		spidernum := len(PSpiderPOOL)

		// 分给几只爬虫
		xxx, _ := util.DevideStringList(piclist, spidernum)

		chs := make(chan int, len(piclist))
		for num, img := range xxx {
			spp, ok := miner.Pool.Get(util.IS(num))
			if !ok {
				spp = PTemp
			}
			go func(temps []string, spider2 *miner.Worker, num int) {
				for _, temp := range temps {
					//filename := util.ValidFileName(temp) + "." + "png" //  知乎图片默认后缀不知
					filename := strings.Replace(strings.Replace(util.ValidFileName(temp), "#", "_", -1), "jpg", "png", -1)
					if util.FileExist(dir + "/" + filename) {
						fmt.Println("文件存在：" + dir + "/" + filename)
						chs <- 0
					} else {
						//fmt.Println("下载:" + temp)
						spider2.SetUrl(temp)
						imgsrc, e := spider2.Get()
						if e != nil {
							fmt.Println("下载出错" + temp + ":" + e.Error())
							chs <- 0
							return
						}
						e = util.SaveToFile(dir+"/"+filename, imgsrc)
						if e == nil {
							fmt.Printf("PSP%d: 成功保存在%s/%s\n", num, dir, filename)
						}
						chs <- 1
						//util.Sleep(1)
						//fmt.Println("暂停1秒")
					}
				}
			}(img, spp, num)
		}

		for i := 0; i < len(piclist); i++ {
			<-chs
		}
	}
}

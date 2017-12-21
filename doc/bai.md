目录结构及获取的数据如下:

```

--- zhihu_windows_amd64.exe
--- zhihu_linux_x86_64
--- cookie.txt
--- data  收藏夹和回答生成的数据在data文件夹
     --- 27761934-如何让自拍的照片看上去像是别人拍的？.xx   *去重标志,如果要重新获取答案,请将`.xx`文件去掉
     --- 27761934  * 回答文件集
        ---zhi-zhi-zhi-41-89-167963702 * 一个用户的回答 包括图片
           --- zhi-zhi-zhi-41-89-167963702的回答.html (里面的图片链接都替换成本地链接)
           --- https###pic1.zhimg.com#v2-22407b227c9a7a19aa0057f38bf6e754_r.png  (已经不是这种样子了)
               https###pic1.zhimg.com#v2-7782ff69838c379173415458b97b5008_xll.jpg
               https###pic1.zhimg.com#v2-c41bf767819fbc61b3ff7bb4c2900884_r.jpg

        ---zhi-zhi-wei-zhi-zhi-36-38-164986419
        ---zhi-zhi-wei-zhi-zhi-hu-hu-wei-hu-hu-164880780

     --- 27761934-html  生成的html集,可以点击查看(可选择防盗链, 请用非火狐浏览器查看)
        --- 1.html
        --- 2.html

--- people   获取用户粉丝数据和所有回答在此文件夹下
--- index.html 点这个就KO了.
```

## 一.小白指南

> 可以下载EXE文件

Golang开发的爬虫，小白用户请下载[释出版本二进制](https://github.com/hunterhug/GoZhihu/releases)中的`zhihu_windows_amd64.exe`，并在同一目录下新建一个`cookie.txt`文件，

打开火狐浏览器后人工登录知乎，按F12，点击网络，刷新一下首页，然后点击第一个出现的`GET /`，找到消息头请求头，复制Cookie，然后粘贴到cookie.txt

![](cookie.png)

点击EXE后,可选JS解决防盗链（这个是你要发布到自己的网站如：[减肥成功是什么感觉？给生活带来哪些改变？](http://www.lenggirl.com/zhihu/26613082-html/1.html)）
我们自己本地看的话就不要选择防盗链了！回答个数已经限制不大于500个。如果没有答案证明Cookie失效，请重新按照上述方法手动修改`cookie.txt`。

你也可以全部图片保存在本地, 这样数据会巨大!

```

        -----------------
        知乎问题信息小助手

        支持:
        1. 从收藏夹https://www.zhihu.com/collection/78172986批量获取很多问题答案
        2. 从问题https://www.zhihu.com/question/28853910批量获取一个问题很多答案
        3. 从某个人https://www.zhihu.com/people/hunterhug批量获取粉丝/偶像和所有回答(待做)

        请您按提示操作（Enter）！答案保存在data或者people文件夹下！

        如果什么都没抓到请往exe同级目录cookie.txt,增加cookie，手动增加cookie见说明

        你亲爱的萌萌~ 努力工作中...
        陈白痴~~~

        联系: Github:hunterhug
        QQ: 459527502   Version: 1.1
        2017.6.29 写于大深圳
        -----------------

萌萌：你有几种选项, 你的决定命运着图片链接是否被替换?

        因为知乎防盗链，把生成的HTML放在你的网站上是看不见图片的！

        选项:
        1. N: 不防盗链(默认), 只能本地浏览器查看远程zhihu图片
        2. Y: JS解决防盗链, 引入JS方便查看远程zhihu图片
        3. X: HTML替换本地图片, 图片会保存, 可以永久观看
        4. Z: 打印抓取的问题html

        请选择:
n
萌萌：不试试抓取图片吗Y/N(默认N)
y
萌萌：从收藏夹获取回答按1，从问题获取回答按2(默认)
2
```

选择Z可以打出所有的问题, 汇集`index.html`

结果：

![](1.png)

![](2.png)

可以看网站 [我的知乎](http://zhihu.lenggirl.com/)
package main

import (
	"fmt"

	"github.com/Akegarasu/blivedm-go/client"
	"github.com/Akegarasu/blivedm-go/message"
	_ "github.com/Akegarasu/blivedm-go/utils"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)
	// 房间号
	c := client.NewClient(27793356)
	// 浏览器获取cookie 打开 https://api.bilibili.com/x/web-interface/nav，需要复制控制台网络请求中的cookie, 不能复制document.cookie
	// 由于 B站 反爬虫改版，现在需要使用已登陆账号的 Cookie 才可以正常获取弹幕。如果不设置 Cookie，获取到的弹幕昵称、UID都被限制。还有可能弹幕限流，无法获取到全部内容。
	c.SetCookie("")
	//弹幕事件
	c.OnDanmaku(func(danmaku *message.Danmaku) {
		if danmaku.Type == message.EmoticonDanmaku {
			fmt.Printf("[弹幕表情] %s：表情URL： %s\n", danmaku.Sender.Uname, danmaku.Emoticon.Url)
		} else {
			fmt.Printf("[弹幕] %s：%s\n", danmaku.Sender.Uname, danmaku.Content)
		}
	})
	// 醒目留言事件
	c.OnSuperChat(func(superChat *message.SuperChat) {
		fmt.Printf("[SC|%d元] %s: %s\n", superChat.Price, superChat.UserInfo.Uname, superChat.Message)
	})
	// 礼物事件
	c.OnGift(func(gift *message.Gift) {
		if gift.CoinType == "gold" {
			fmt.Printf("[礼物] %s 的 %s %d 个 共%.2f元\n", gift.Uname, gift.GiftName, gift.Num, float64(gift.Num*gift.Price)/1000)
		}
	})
	// 上舰事件
	c.OnGuardBuy(func(guardBuy *message.GuardBuy) {
		fmt.Printf("[大航海] %s 开通了 %d 等级的大航海，金额 %d 元\n", guardBuy.Username, guardBuy.GuardLevel, guardBuy.Price/1000)
	})
	// 监听自定义事件
	c.RegisterCustomEventHandler("STOP_LIVE_ROOM_LIST", func(s string) {
		// data := gjson.Get(s, "data").String()
		// fmt.Printf("STOP_LIVE_ROOM_LIST: %s\n", data)
	})

	err := c.Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("started")
	// 需要自行阻塞什么方法都可以
	select {}
}

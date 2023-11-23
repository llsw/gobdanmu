package main

import (
	"github.com/Akegarasu/blivedm-go/client"
	"github.com/Akegarasu/blivedm-go/message"
	_ "github.com/Akegarasu/blivedm-go/utils"
	log "github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/spf13/viper"
)

func loadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.SetDefault("room", 7777)
	viper.SetDefault("cookie", "")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("read config failed: %v", err)
	}
}

func main() {
	loadConfig()
	log.SetLevel(log.LevelDebug)
	// 房间号
	c := client.NewClient(viper.GetInt("room"))
	// 浏览器获取cookie 打开 https://api.bilibili.com/x/web-interface/nav，需要复制控制台网络请求中的cookie, 不能复制document.cookie
	// 由于 B站 反爬虫改版，现在需要使用已登陆账号的 Cookie 才可以正常获取弹幕。如果不设置 Cookie，获取到的弹幕昵称、UID都被限制。还有可能弹幕限流，无法获取到全部内容。
	c.SetCookie(viper.GetString("cookie"))
	//弹幕事件
	c.OnDanmaku(func(danmaku *message.Danmaku) {
		if danmaku.Type == message.EmoticonDanmaku {
			log.Infof("[弹幕表情] %s：表情URL： %s\n", danmaku.Sender.Uname, danmaku.Emoticon.Url)
		} else {
			log.Infof("[弹幕] %s：%s\n", danmaku.Sender.Uname, danmaku.Content)
		}
	})
	// 醒目留言事件
	c.OnSuperChat(func(superChat *message.SuperChat) {
		log.Infof("[SC|%d元] %s: %s\n", superChat.Price, superChat.UserInfo.Uname, superChat.Message)
	})
	// 礼物事件
	c.OnGift(func(gift *message.Gift) {
		if gift.CoinType == "gold" {
			log.Infof("[礼物] %s 的 %s %d 个 共%.2f元\n", gift.Uname, gift.GiftName, gift.Num, float64(gift.Num*gift.Price)/1000)
		}
	})
	// 上舰事件
	c.OnGuardBuy(func(guardBuy *message.GuardBuy) {
		log.Infof("[大航海] %s 开通了 %d 等级的大航海，金额 %d 元\n", guardBuy.Username, guardBuy.GuardLevel, guardBuy.Price/1000)
	})
	// 监听自定义事件
	c.RegisterCustomEventHandler("STOP_LIVE_ROOM_LIST", func(s string) {
		// data := gjson.Get(s, "data").String()
		// log.Infof("STOP_LIVE_ROOM_LIST: %s\n", data)
	})

	err := c.Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Info("started")
	// 需要自行阻塞什么方法都可以
	select {}
}

package main

import (
	"github.com/Akegarasu/blivedm-go/client"
	_ "github.com/Akegarasu/blivedm-go/utils"
	log "github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/llsw/gobdanmu/src/msg"
	//	logr "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
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
	// ogr.SetLevel(logr.DebugLevel)
	loadConfig()
	fileName := "/Users/xmk/git/go/gobdanmu/log.log"
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Errorf("log file:%s creae error: ", fileName, err.Error())
	} else {
		log.SetOutput(f)
	}
	log.SetLevel(log.LevelDebug)
	// 房间号
	c := client.NewClient(viper.GetInt("room"))
	// 浏览器获取cookie 打开 https://api.bilibili.com/x/web-interface/nav，需要复制控制台网络请求中的cookie, 不能复制document.cookie
	// 由于 B站 反爬虫改版，现在需要使用已登陆账号的 Cookie 才可以正常获取弹幕。如果不设置 Cookie，获取到的弹幕昵称、UID都被限制。还有可能弹幕限流，无法获取到全部内容。
	c.SetCookie(viper.GetString("cookie"))

	//弹幕事件
	c.OnDanmaku(msg.OnDanmu)
	// 醒目留言事件
	c.OnSuperChat(msg.OnSuperChat)
	// 礼物事件
	c.OnGift(msg.OnGift)
	// 上舰事件
	c.OnGuardBuy(msg.OnGuardBuy)
	// 欢迎进入直播间
	msg.OnWelcome()
	c.RegisterCustomEventHandler(msg.WELCOME, msg.GetEventHandler(msg.WELCOME))
	c.RegisterCustomEventHandler(msg.WELCOME_GUARD, msg.GetEventHandler(msg.WELCOME_GUARD))
	//进入直播间
	msg.OnInterActWord()
	c.RegisterCustomEventHandler(msg.INTERACT_WORD, msg.GetEventHandler(msg.INTERACT_WORD))
	//进入直播间
	msg.OnEnterEffect()
	c.RegisterCustomEventHandler(msg.ENTRY_EFFECT, msg.GetEventHandler(msg.ENTRY_EFFECT))
	// 先设置个默认事件
	msg.DefaultHandler(msg.STOP_LIVE_ROOM_LIST)
	// 监听自定义事件
	c.RegisterCustomEventHandler(msg.STOP_LIVE_ROOM_LIST, msg.GetEventHandler(msg.STOP_LIVE_ROOM_LIST))

	err = c.Start()
	if err != nil {
		log.Fatal(err)
	}
	// log.Info("started")
	// 需要自行阻塞什么方法都可以
	//select {}
	msg.Start()
}

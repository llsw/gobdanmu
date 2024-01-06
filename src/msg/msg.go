package msg

import (
	"fmt"
	"github.com/Akegarasu/blivedm-go/message"
	log "github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/gookit/event"
	//"github.com/tidwall/gjson"
)

const (
	STOP_LIVE_ROOM_LIST = "STOP_LIVE_ROOM_LIST"
	ON_DAN_MU           = "ON_DAN_MU"
	WELCOME             = "WELCOME"
	WELCOME_GUARD       = "WELCOME_GUARD"
)

// 弹幕事件
func OnDanmu(danmaku *message.Danmaku) {
	var dmsg *DanMuMsg
	if danmaku.Type == message.EmoticonDanmaku {
		dmsg = NewDanMuMsg(
			"[弹幕表情]",
			fmt.Sprintf(" %s", danmaku.Sender.Uname),
			danmaku.Emoticon.Url,
		)
		// log.Infof("[弹幕表情] %s：表情URL： %s\n", danmaku.Sender.Uname, danmaku.Emoticon.Url)
	} else {
		dmsg = NewDanMuMsg(
			"[弹幕]",
			fmt.Sprintf(" %s", danmaku.Sender.Uname),
			danmaku.Content,
		)
		// log.Infof("[弹幕] %s：%s\n", danmaku.Sender.Uname, danmaku.Content)
	}
	event.MustFire(ON_DAN_MU, event.M{"data": dmsg})
}

// 醒目留言事件
func OnSuperChat(superChat *message.SuperChat) {
	log.Infof("[SC|%d元] %s: %s\n", superChat.Price, superChat.UserInfo.Uname, superChat.Message)
}

// 礼物事件
func OnGift(gift *message.Gift) {
	if gift.CoinType == "gold" {
		dmsg := NewDanMuMsg(
			"[礼物]",
			fmt.Sprintf(" %s", gift.Uname),
			fmt.Sprintf(
				"%s %d 个 共%.2f元",
				gift.GiftName, gift.Num,
				float64(gift.Num*gift.Price)/1000,
			),
		)
		event.MustFire(ON_DAN_MU, event.M{"data": dmsg})
		// log.Infof("[礼物] %s 的 %s %d 个 共%.2f元\n", gift.Uname, gift.GiftName, gift.Num, float64(gift.Num*gift.Price)/1000)
	}
}

// 上舰事件
func OnGuardBuy(guardBuy *message.GuardBuy) {
	dmsg := NewDanMuMsg(
		"[大航海]",
		fmt.Sprintf(" %s", guardBuy.Username),
		fmt.Sprintf(
			"开通了 %d 等级的大航海，金额 %d 元\n",
			guardBuy.GuardLevel, guardBuy.Price/1000),
	)
	event.MustFire(ON_DAN_MU, event.M{"data": dmsg})
	// log.Infof("[大航海] %s 开通了 %d 等级的大航海，金额 %d 元\n", guardBuy.Username, guardBuy.GuardLevel, guardBuy.Price/1000)
}

func GetEventHandler(eventName string) func(data string) {
	return func(data string) {
		event.MustFire(eventName, event.M{"data": data})
	}
}

// 监听自定义事件
func DefaultHandler(eventName string) {
	event.On(eventName, event.ListenerFunc(func(e event.Event) error {
		if _, ok := e.Data()["data"]; ok {
			// data := gjson.Get(s.(string), "data").String()
			// 还是不要打印
			// log.Infof("%s: %s\n", eventName, data)
			return nil
		} else {
			return fmt.Errorf("event:%s no data", eventName)
		}
	}))
}

// 弹幕事件
// func OnWelcome(danmaku *message.Danmaku) {
// 	dmsg := NewDanMuMsg(
// 		"[弹幕表情]",
// 		fmt.Sprintf(" %s", danmaku.Sender.Uname),
// 		danmaku.Emoticon.Url,
// 	)
// }

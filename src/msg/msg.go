package msg

import (
	"fmt"
	"github.com/Akegarasu/blivedm-go/message"
	log "github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/gookit/event"
	"github.com/tidwall/gjson"
	"runtime/debug"
)

const (
	STOP_LIVE_ROOM_LIST = "STOP_LIVE_ROOM_LIST"
	ON_DAN_MU           = "ON_DAN_MU"
	WELCOME             = "WELCOME"
	WELCOME_GUARD       = "WELCOME_GUARD"
	INTERACT_WORD       = "INTERACT_WORD"
	ENTRY_EFFECT        = "ENTRY_EFFECT"
)

// 弹幕事件
func OnDanmu(danmaku *message.Danmaku) {
	var dmsg *DanMuMsg
	if danmaku.Type == message.EmoticonDanmaku {
		dmsg = NewDanMuMsg(
			"[弹幕表情]",
			fmt.Sprintf(" %s", danmaku.Sender.Uname),
			danmaku.Emoticon.Url,
			0,
		)
		// log.Infof("[弹幕表情] %s：表情URL： %s\n", danmaku.Sender.Uname, danmaku.Emoticon.Url)
	} else {
		dmsg = NewDanMuMsg(
			"[弹幕]",
			fmt.Sprintf(" %s", danmaku.Sender.Uname),
			danmaku.Content,
			0,
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
			0,
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
		0,
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

// ENTRY_EFFECT body:{"cmd":"ENTRY_EFFECT","data":{"id":135,"uid":297462432,"target_id":1402851423,"mock_effect":0,"face":"https://i0.hdslb.com/bfs/face/80b9e93131f7815ada27960858850391ccbf1097.jpg","privilege_type":0,"copy_writing":"欢迎 \u003c%梯度上升%\u003e 进入直播间","copy_color":"#000000","highlight_color":"#FFF100","priority":1,"basemap_url":"https://i0.hdslb.com/bfs/live/mlive/da6933ea70f31c4df63f4b68b735891284888357.png","show_avatar":1,"effective_time":1,"web_basemap_url":"https://i0.hdslb.com/bfs/live/mlive/da6933ea70f31c4df63f4b68b735891284888357.png","web_effective_time":2,"web_effect_close":0,"web_close_time":900,"business":3,"copy_writing_v2":"欢迎 \u003c^icon^\u003e \u003c%梯度上升%\u003e 进入直播间","icon_list":[1],"max_delay_time":7,"trigger_time":1704534340507624467,"identities":22,"effect_silent_time":0,"effective_time_new":0,"web_dynamic_url_webp":"","web_dynamic_url_apng":"","mobile_dynamic_url_webp":"","wealthy_info":null,"new_style":0,"is_mystery":false,"uinfo":{"uid":297462432,"base":{"name":"梯度上升","face":"https://i0.hdslb.com/bfs/face/80b9e93131f7815ada27960858850391ccbf1097.jpg","name_color":0,"is_mystery":false,"risk_ctrl_info":{"name":"梯度上升","face":"https://i0.hdslb.com/bfs/face/80b9e93131f7815ada27960858850391ccbf1097.jpg"},"origin_info":{"name":"梯度上升","face":"https://i0.hdslb.com/bfs/face/80b9e93131f7815ada27960858850391ccbf1097.jpg"}}}}}
// INFO[0217] debug cmdINTERACT_WORD body:{"cmd":"INTERACT_WORD","data":{"contribution":{"grade":1},"contribution_v2":{"grade":1,"rank_type":"monthly_rank","text":"月榜前3用户"},"core_user_type":0,"dmscore":28,"fans_medal":{"anchor_roomid":310773,"guard_level":0,"icon_id":0,"is_lighted":0,"medal_color":6067854,"medal_color_border":12632256,"medal_color_end":12632256,"medal_color_start":12632256,"medal_level":3,"medal_name":"水心","score":939,"special":"","target_id":927587},"group_medal":null,"identities":[1],"is_mystery":false,"is_spread":0,"msg341568,"spread_desc":"","spread_info":"","tail_icon":0,"tail_text":"","timestamp":1704534341,"trigger_time":1704534340507624400,"uid":297462432,"uinfo":{"base":{"face":"https://i0.hdslb.com/bfs/face/80b9e93131f7815ada27960858850391ccbf1097.jpg","is_mystery":false,"name":"梯度上升","name_color":0,"origin_info":{"face":"https://i0.hdslb.com/bfs/face/80b9e93131f7815ada27960858850391ccbf1097.jpg","name":"梯度上升"}

// 欢迎进入直播间
func OnWelcome() {
	event.On(WELCOME, event.ListenerFunc(func(e event.Event) error {
		defer func() {
			if pan := recover(); pan != nil {
				log.Errorf("handler packet error: %v\n%s", pan, debug.Stack())
			}
		}()
		if s, ok := e.Data()["data"]; ok {
			data := gjson.Get(s.(string), "data").String()
			log.Infof("%s: %s\n", WELCOME, data)
			return nil
		} else {
			return fmt.Errorf("event:%s no data", WELCOME)
		}
	}))
}

func OnInterActWord() {
	eventName := INTERACT_WORD
	event.On(eventName, event.ListenerFunc(func(e event.Event) error {
		defer func() {
			if pan := recover(); pan != nil {
				log.Errorf("handler packet error: %v\n%s", pan, debug.Stack())
			}
		}()
		if s, ok := e.Data()["data"]; ok {
			data := gjson.Get(s.(string), "data").String()
			log.Infof("%s: %s\n", eventName, data)
			// dmsg := NewDanMuMsg(
			// 	"[欢迎]",
			// 	fmt.Sprintf(" %s", data),
			// 	"进入直播间",
			// 	1,
			// )
			// event.MustFire(ON_DAN_MU, event.M{"data": dmsg})
			return nil
		} else {
			return fmt.Errorf("event:%s no data", eventName)
		}
	}))
}

func OnEnterEffect() {
	eventName := ENTRY_EFFECT
	event.On(eventName, event.ListenerFunc(func(e event.Event) error {
		defer func() {
			if pan := recover(); pan != nil {
				log.Errorf("handler packet error: %v\n%s", pan, debug.Stack())
			}
		}()
		if s, ok := e.Data()["data"]; ok {
			data := gjson.Get(s.(string), "data").String()
			log.Infof("%s: %s\n", eventName, data)
			// dmsg := NewDanMuMsg(
			// 	"[欢迎]",
			// 	fmt.Sprintf(" %s", data),
			// 	"进入直播间",
			// 	1,
			// )
			// event.MustFire(ON_DAN_MU, event.M{"data": dmsg})
			return nil
		} else {
			return fmt.Errorf("event:%s no data", eventName)
		}
	}))
}

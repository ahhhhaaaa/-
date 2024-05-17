package main

import (
	"log"
	"os"
)

// handler.go:75 超过限长1420不能回复问题 ok
// 不能连网问题
// 程序不能后台静默运行
import (
	"io"
	"openai-wechat/bot"
	"openai-wechat/handler"
	"openai-wechat/utils"
)

func init() {
	// 1. log init
	f, _ := os.OpenFile("run.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	log.SetOutput(io.MultiWriter(os.Stdout, f))
	log.SetPrefix("openai-wechat")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	// 2. WeChat bot init
	if err := bot.Init(); err != nil {
		log.Fatalf("微信登录失败，错误信息为：%v", err)
	}
	log.Println("登录成功")
}

func main() {
	// 获取登录的用户
	self, err := bot.Bot.GetCurrentUser()
	if err != nil {
		log.Printf("%+v", err)
		return
	}
	log.Printf("self=%s", utils.MarshalAnyToString(self))
	bot.Bot.MessageHandler = handler.MessageHandler // 微信消息回调注册
	bot.Bot.Block()
}

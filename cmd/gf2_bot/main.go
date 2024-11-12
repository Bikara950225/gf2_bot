package main

import (
	"context"
	bch "gf2_bot/internal/bot_cmd_handler"
	"gf2_bot/internal/bot_cmd_handler/debug_handler"
	"log"
	"time"

	botctl "gf2_bot/internal/bot_controller"

	"github.com/gin-gonic/gin"
	"github.com/tencent-connect/botgo"
	"github.com/tencent-connect/botgo/event"
	"github.com/tencent-connect/botgo/interaction/webhook"
	"github.com/tencent-connect/botgo/token"
)

const (
	appId_     = "102351637"
	appSecret_ = "K5qbN9vhTF1naNAxkXK8wkYMAyncRG5u"
	appToken   = "9c2A8M2NSNuUJM8VV8FB3KY0QsiySgzK"
)

func main() {
	ctx := context.Background()

	credentials := &token.QQBotCredentials{
		AppID:     appId_,
		AppSecret: appSecret_,
	}
	tokenSource := token.NewQQBotTokenSource(credentials)

	if err := token.StartRefreshAccessToken(ctx, tokenSource); err != nil {
		log.Fatalln(err)
	}

	// 初始化 openapi，正式环境
	qqBotOpenapi := botgo.NewOpenAPI(credentials.AppID, tokenSource).
		WithTimeout(5 * time.Second).SetDebug(true)
	//
	botHandlerProxy := bch.NewBotCmdHandlerProxy()
	debugHandler := debug_handler.NewDebugHandler()
	botHandlerProxy.Register("角色详情", debugHandler)
	botHandlerProxy.Register("指令测试1", debugHandler)
	botHandlerProxy.Register("抽卡分析", debugHandler)
	// 机器人的控制器，将指令转发
	botController := botctl.NewBotController(qqBotOpenapi, botHandlerProxy)

	_ = event.RegisterHandlers(
		// 注册事件
		botController.MessageHandler(),
	)

	engine := gin.Default()
	engine.Any("/qqbot", func(ginc *gin.Context) {
		webhook.HTTPHandler(ginc.Writer, ginc.Request, credentials)
	})

	if err := engine.Run(":50008"); err != nil {
		log.Fatalln(err)
	}
}

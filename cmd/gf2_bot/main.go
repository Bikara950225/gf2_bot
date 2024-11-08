package main

import (
	"context"
	bch "gf2_bot/internal/bot_cmd_handler"
	"log"
	"time"

	botctl "gf2_bot/internal/bot_controller"

	"github.com/gin-gonic/gin"
	"github.com/tencent-connect/botgo"
	"github.com/tencent-connect/botgo/dto"
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

func generateDemoMessage(input string, data dto.Message) *dto.MessageToCreate {
	log.Printf("收到指令: %+v\n", input)
	msg := ""
	if len(input) > 0 {
		msg += "收到:" + input
	}
	for _, _v := range data.Attachments {
		msg += ",收到文件类型:" + _v.ContentType
	}
	return &dto.MessageToCreate{
		Timestamp: time.Now().UnixMilli(),
		Content:   msg,
		MessageReference: &dto.MessageReference{
			// 引用这条消息
			MessageID:             data.ID,
			IgnoreGetMessageError: true,
		},
		MsgID: data.ID,
	}
}

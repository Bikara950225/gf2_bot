package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tencent-connect/botgo"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/dto/message"
	"github.com/tencent-connect/botgo/event"
	"github.com/tencent-connect/botgo/interaction/webhook"
	"github.com/tencent-connect/botgo/openapi"
	"github.com/tencent-connect/botgo/token"
	"log"
	"strings"
	"time"
)

var api openapi.OpenAPI

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
	api = botgo.NewOpenAPI(credentials.AppID, tokenSource).WithTimeout(5 * time.Second).SetDebug(true)

	_ = event.RegisterHandlers(
		// 注册事件
		groupATMessageEventHandler(),
	)

	engine := gin.Default()
	engine.Any("/qqbot", func(ginc *gin.Context) {
		webhook.HTTPHandler(ginc.Writer, ginc.Request, credentials)
	})

	//if err := engine.Run(":50008"); err != nil {
	//	log.Fatalln(err)
	//}
	if err := engine.RunTLS(":50008", "", ""); err != nil {
		log.Fatalln(err)
	}
}

// GroupATMessageEventHandler 实现处理 at 消息的回调
func groupATMessageEventHandler() event.GroupATMessageEventHandler {
	return func(event *dto.WSPayload, data *dto.WSGroupATMessageData) error {
		input := strings.ToLower(message.ETLInput(data.Content))
		return processGroupMessage(input, data)

	}
}

func processGroupMessage(content string, data *dto.WSGroupATMessageData) error {
	msg := generateDemoMessage(content, dto.Message(*data))
	if err := sendGroupReply(context.Background(), data.GroupID, msg); err != nil {
		_ = sendGroupReply(context.Background(), data.GroupID, genErrMessage(dto.Message(*data), err))
	}
	return nil
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

func sendGroupReply(ctx context.Context, groupID string, toCreate dto.APIMessage) error {
	log.Printf("EVENT ID: %v\n", toCreate.GetEventID())
	if _, err := api.PostGroupMessage(ctx, groupID, toCreate); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func genErrMessage(data dto.Message, err error) *dto.MessageToCreate {
	return &dto.MessageToCreate{
		Timestamp: time.Now().UnixMilli(),
		Content:   fmt.Sprintf("处理异常:%v", err),
		MessageReference: &dto.MessageReference{
			// 引用这条消息
			MessageID:             data.ID,
			IgnoreGetMessageError: true,
		},
		MsgID: data.ID,
	}
}

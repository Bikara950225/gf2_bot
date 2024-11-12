package bot_controller

import (
	"context"
	"fmt"
	bch "gf2_bot/internal/bot_cmd_handler"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/dto/message"
	"github.com/tencent-connect/botgo/event"
	"github.com/tencent-connect/botgo/openapi"
	"log"
	"strings"
	"time"
)

type botController struct {
	qqBotOpenapi    openapi.OpenAPI
	botHandlerProxy bch.BotCmdHandlerProxy
}

func NewBotController(api openapi.OpenAPI, botHandlerProxy bch.BotCmdHandlerProxy) *botController {
	return &botController{
		qqBotOpenapi:    api,
		botHandlerProxy: botHandlerProxy,
	}
}

func (s *botController) MessageHandler() event.GroupATMessageEventHandler {
	return func(event *dto.WSPayload, data *dto.WSGroupATMessageData) error {
		input := strings.ToLower(message.ETLInput(data.Content))
		message := dto.Message(*data)

		// TODO 上下文放到哪里？比如当前通讯的qq用户的用户信息等等
		ctx := bch.WithDtoMessage(context.Background(), &message)
		// 解析input，派发到不同的handler处理
		msg, err := s.botHandlerProxy.Handler(ctx, input)
		if err != nil {
			msg = s.genErrMessage(dto.Message(*data), err)
		}
		return s.sendGroupReply(ctx, data.GroupID, msg)
	}
}

func (s *botController) genErrMessage(data dto.Message, err error) *dto.MessageToCreate {
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

func (s *botController) sendGroupReply(ctx context.Context, groupID string, toCreate dto.APIMessage) error {
	log.Printf("EVENT ID: %v\n", toCreate.GetEventID())
	if _, err := s.qqBotOpenapi.PostGroupMessage(ctx, groupID, toCreate); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

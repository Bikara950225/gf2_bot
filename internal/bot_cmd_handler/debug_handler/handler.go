package debug_handler

import (
	"context"
	"fmt"
	bch "gf2_bot/internal/bot_cmd_handler"
	"github.com/tencent-connect/botgo/dto"
	"log"
	"time"
)

type debugHandler struct{}

var _ bch.BotCmdHandlerPlugin = (*debugHandler)(nil)

func NewDebugHandler() bch.BotCmdHandlerPlugin {
	return &debugHandler{}
}

func (s *debugHandler) Handler(ctx context.Context, params ...string) (*dto.MessageToCreate, error) {
	command := bch.GetCommandFromContext(ctx)
	actCommand := fmt.Sprintf("%s %v", command, params)
	data := bch.GetDtoMessageFromContext(ctx)

	log.Printf("收到指令: %s\n", actCommand)
	msg := ""
	if len(actCommand) > 0 {
		msg += "收到:" + actCommand
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
	}, nil
}

package bot_cmd_handler

import (
	"context"
	"github.com/tencent-connect/botgo/dto"
)

type (
	commandKey    struct{}
	dtoMessageKey struct{}
)

func WithCommand(parent context.Context, command string) context.Context {
	return context.WithValue(parent, commandKey{}, command)
}

func GetCommandFromContext(parent context.Context) string {
	return parent.Value(commandKey{}).(string)
}

func WithDtoMessage(parent context.Context, data *dto.Message) context.Context {
	return context.WithValue(parent, dtoMessageKey{}, data)
}

func GetDtoMessageFromContext(parent context.Context) *dto.Message {
	return parent.Value(dtoMessageKey{}).(*dto.Message)
}

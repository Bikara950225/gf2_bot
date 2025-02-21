package bot_cmd_handler

import (
	"context"
	"fmt"
	"github.com/tencent-connect/botgo/dto"
	"strings"
)

type BotCmdHandlerProxy interface {
	Handler(ctx context.Context, input string) (ret *dto.MessageToCreate, err error)
	Register(command string, handlePlugin BotCmdHandlerPlugin)
}

type BotCmdHandlerPlugin interface {
	Handler(ctx context.Context, params ...string) (ret *dto.MessageToCreate, err error)
}

type botCmdHandlerProxyImpl struct {
	proxy map[string]BotCmdHandlerPlugin
}

func NewBotCmdHandlerProxy() BotCmdHandlerProxy {
	return &botCmdHandlerProxyImpl{
		proxy: map[string]BotCmdHandlerPlugin{},
	}
}

func (b *botCmdHandlerProxyImpl) Register(command string, handlePlugin BotCmdHandlerPlugin) {
	b.proxy[command] = handlePlugin
}

func (b *botCmdHandlerProxyImpl) Handler(ctx context.Context, input string) (ret *dto.MessageToCreate, err error) {
	command, params, err := parseInput(input)
	if err != nil {
		return nil, err
	}

	if plugin, ok := b.proxy[command]; ok {
		ctx = WithCommand(ctx, command)

		return plugin.Handler(ctx, params...)
	} else {
		return nil, fmt.Errorf("未知的指令：%s", command)
	}
}

func parseInput(input string) (string, []string, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "", nil, fmt.Errorf("输入指令为空，不合法")
	}

	inputSplits := strings.Split(input, " ")
	command := inputSplits[0]
	if !strings.HasPrefix(command, "/") {
		return "", nil, fmt.Errorf("非法指令，格式不正确: %s", command)
	}
	command = strings.TrimPrefix(command, "/")

	return command, inputSplits[1:], nil
}

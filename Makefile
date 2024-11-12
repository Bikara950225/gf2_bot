.PHONY: build-linux

build-linux:
	GOOS="linux" GOARCH="amd64" go build -o gf2_bot_linux gf2_bot/cmd/gf2_bot

program_async:
	scp ./gf2_bot_linux root@gf2.zhajibot.sale:/root

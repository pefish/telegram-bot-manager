package main

import (
	"github.com/pefish/go-commander"
	go_logger "github.com/pefish/go-logger"
	"github.com/pefish/telegram-bot-manager/cmd/telegram-bot-manager/command"
	"github.com/pefish/telegram-bot-manager/version"
)

func main() {
	commanderInstance := commander.NewCommander(version.AppName, version.Version, version.AppName+" is a robot manager for telegram, enjoy this!!。author：pefish")
	commanderInstance.RegisterDefaultSubcommand(&commander.SubcommandInfo{
		Desc:       "监控机器人收到的命令，然后给予回复",
		Args:       nil,
		Subcommand: command.NewDefaultCommand(),
	})
	commanderInstance.RegisterSubcommand("send", &commander.SubcommandInfo{
		Desc:       "向群里发消息",
		Args:       []string{"group_id", "msg"},
		Subcommand: command.NewSendCommand(),
	})
	err := commanderInstance.Run()
	if err != nil {
		go_logger.Logger.Error(err)
	}
}

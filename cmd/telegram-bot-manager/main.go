package main

import (
	"github.com/pefish/go-commander"
	go_logger "github.com/pefish/go-logger"
	"github.com/pefish/telegram-bot-manager/cmd/telegram-bot-manager/command"
	"github.com/pefish/telegram-bot-manager/version"
)

func main() {
	commanderInstance := commander.NewCommander(version.AppName, version.Version, version.AppName + " is a robot manager for telegram, enjoy this!!。author：pefish")
	commanderInstance.RegisterDefaultSubcommand(command.NewDefaultCommand())
	err := commanderInstance.Run()
	if err != nil {
		go_logger.Logger.Error(err)
	}
}

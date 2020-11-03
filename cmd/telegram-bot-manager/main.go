package main

import (
	"github.com/pefish/go-commander"
	"github.com/pefish/telegram-bot-manager/pkg/command"
	"github.com/pefish/telegram-bot-manager/version"
	"log"
)

func main() {
	commanderInstance := commander.NewCommander(version.AppName, version.Version, version.AppName + " is a robot manager for telegram, enjoy this!!。author：pefish")
	commanderInstance.RegisterDefaultSubcommand(command.NewDefaultCommand())
	err := commanderInstance.Run()
	if err != nil {
		log.Fatal(err)
	}
}

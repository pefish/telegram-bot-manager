package command

import (
	"time"

	"github.com/pefish/go-commander"
	tg_sender "github.com/pefish/tg-sender"
)

type Config struct {
	commander.BasicConfig
	Token string `json:"token" default:"" usage:"Bot token."`
}

var SendCommandConfig Config

type SendCommand struct {
}

func NewSendCommand() *SendCommand {
	return &SendCommand{}
}

func (dc *SendCommand) Config() interface{} {
	return &SendCommandConfig
}

func (dc *SendCommand) Data() interface{} {
	return nil
}

func (s *SendCommand) OnExited(commander *commander.Commander) error {
	return nil
}

func (s *SendCommand) Init(commander *commander.Commander) error {
	return nil
}

func (s *SendCommand) Start(commander *commander.Commander) error {
	err := tg_sender.NewTgSender(SendCommandConfig.Token).SendMsg(
		&tg_sender.MsgStruct{
			ChatId: commander.Args["group_id"],
			Msg:    commander.Args["msg"],
		},
		0,
	)
	if err != nil {
		return err
	}
	time.Sleep(time.Second)
	return nil
}

// TOKEN=XX go run ./cmd/telegram-bot-manager send -- "XX" test

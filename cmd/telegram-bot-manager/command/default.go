package command

import (
	"flag"
	"github.com/pefish/go-commander"
	go_config "github.com/pefish/go-config"
	go_error "github.com/pefish/go-error"
	telegram_robot "github.com/pefish/telegram-bot-manager/pkg/telegram-robot"
	"io/ioutil"
	"os"
)

type DefaultCommand struct {
	robot *telegram_robot.Robot
}

func NewDefaultCommand() *DefaultCommand {
	return &DefaultCommand{}
}

func (s *DefaultCommand) DecorateFlagSet(flagSet *flag.FlagSet) error {
	return nil
}

func (s *DefaultCommand) OnExited() error {
	err := s.robot.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *DefaultCommand) Start(data commander.StartData) error {
	commandsJsFile, err := go_config.Config.GetString("commandsJsFile")
	if err != nil {
		return go_error.WithStack(err)
	}
	fs, err := os.Open(commandsJsFile)
	if err != nil {
		return go_error.WithStack(err)
	}
	defer fs.Close()
	scriptBytes, err := ioutil.ReadAll(fs)
	if err != nil {
		return go_error.WithStack(err)
	}

	token, err := go_config.Config.GetString("token")
	if err != nil {
		return go_error.WithStack(err)
	}
	s.robot = telegram_robot.NewRobot(string(scriptBytes), token)
	err = s.robot.Start(data.DataDir)
	if err != nil {
		return go_error.WithStack(err)
	}
	return nil
}

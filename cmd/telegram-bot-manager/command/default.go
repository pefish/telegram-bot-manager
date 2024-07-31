package command

import (
	"io"
	"os"
	"strings"
	"time"

	"github.com/pefish/go-commander"
	go_error "github.com/pefish/go-error"
	vm3 "github.com/pefish/go-jsvm"
	go_logger "github.com/pefish/go-logger"
	"github.com/pefish/telegram-bot-manager/pkg/global"
	telegram_robot "github.com/pefish/telegram-bot-manager/pkg/telegram-robot"
)

type DefaultCommand struct {
	robot *telegram_robot.Robot
}

func NewDefaultCommand() *DefaultCommand {
	return &DefaultCommand{}
}

func (dc *DefaultCommand) Config() interface{} {
	return &global.GlobalConfig
}

func (dc *DefaultCommand) Data() interface{} {
	return nil
}

func (s *DefaultCommand) OnExited(commander *commander.Commander) error {
	if s.robot != nil {
		err := s.robot.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *DefaultCommand) Init(commander *commander.Commander) error {
	return nil
}

func (s *DefaultCommand) Start(commander *commander.Commander) error {
	fs, err := os.Open(global.GlobalConfig.CommandsJsFile)
	if err != nil {
		return go_error.WithStack(err)
	}
	defer fs.Close()
	scriptBytes, err := io.ReadAll(fs)
	if err != nil {
		return go_error.WithStack(err)
	}

	s.robot = telegram_robot.NewRobot(global.GlobalConfig.Token, 2*time.Second)
	s.robot.SetLogger(go_logger.Logger)

	vm := vm3.NewVm(string(scriptBytes) + "\n" + `
commands["/help"] = {
    func: function (args) {
        var result = ""
        for (var k of Object.keys(commands)) {
            if (k === "/help") {
                continue
            }
            result += k + "  " + (commands[k].desc || "") + "\n"
        }
        if (result === "") {
            return "No useful commands!!!"
        }
        result = "You can use commands：\n\n" + result
        return result
    }
}

function execute(command, args) {
    if (!commands[command]) {
        return "Sorry, I don't understand."
    }
    if (!commands[command].func) {
        return "Internal Error!!! func param not be set, contact admin please!"
    }
    return commands[command].func(args)
}
`)
	if err != nil {
		return go_error.WithStack(err)
	}

	err = s.robot.Start(commander.Ctx, commander.DataDir, func(command string, data string) (string, error) {
		commandTextArr := strings.Split(data, " ")
		result, err := vm.RunFunc("execute", []interface{}{
			command, commandTextArr,
		})
		if err != nil {
			return "", err
		}
		return result.(string), nil
	})
	if err != nil {
		return go_error.WithStack(err)
	}
	return nil
}

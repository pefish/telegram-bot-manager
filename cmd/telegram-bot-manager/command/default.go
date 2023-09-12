package command

import (
	"flag"
	"github.com/pefish/go-commander"
	go_config "github.com/pefish/go-config"
	go_error "github.com/pefish/go-error"
	vm3 "github.com/pefish/go-jsvm/pkg/vm"
	go_logger "github.com/pefish/go-logger"
	telegram_robot "github.com/pefish/telegram-bot-manager/pkg/telegram-robot"
	"io/ioutil"
	"os"
	"strings"
	"time"
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

func (s *DefaultCommand) OnExited(data *commander.StartData) error {
	if s.robot != nil {
		err := s.robot.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *DefaultCommand) Init(data *commander.StartData) error {
	return nil
}

func (s *DefaultCommand) Start(data *commander.StartData) error {
	commandsJsFile, err := go_config.ConfigManagerInstance.GetString("commandsJsFile")
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

	token, err := go_config.ConfigManagerInstance.GetString("token")
	if err != nil {
		return go_error.WithStack(err)
	}
	s.robot = telegram_robot.NewRobot(token, 2*time.Second)
	s.robot.SetLogger(go_logger.Logger)

	vm, err := vm3.NewVmAndLoad(string(scriptBytes) + "\n" + `
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

	err = s.robot.Start(data.ExitCancelCtx, data.DataDir, func(command string, data string) (string, error) {
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

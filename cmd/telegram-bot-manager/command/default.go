package command

import (
	"flag"
	"github.com/pefish/go-commander"
	go_config "github.com/pefish/go-config"
	go_error "github.com/pefish/go-error"
	go_logger "github.com/pefish/go-logger"
	telegram_robot "github.com/pefish/telegram-bot-manager/pkg/telegram-robot"
	vm2 "github.com/pefish/telegram-bot-manager/pkg/vm"
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
	s.robot = telegram_robot.NewRobot(token, 2 * time.Second)
	s.robot.SetLogger(go_logger.Logger)

	vm := vm2.NewVm()
	_, err = vm.RunString(string(scriptBytes) + "\n" + `
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

	var fn func(string, []string) string
	err = vm.ExportTo(vm.Get("execute"), &fn)
	if err != nil {
		return go_error.WithStack(err)
	}

	err = s.robot.Start(data.ExitCancelCtx, data.DataDir, func(command string, data string) string {
		commandTextArr := strings.Split(data, " ")
		return fn(command, commandTextArr)
	})
	if err != nil {
		return go_error.WithStack(err)
	}
	return nil
}

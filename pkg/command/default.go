package command

import (
	"flag"
	"fmt"
	go_config "github.com/pefish/go-config"
	go_decimal "github.com/pefish/go-decimal"
	"github.com/pefish/go-http"
	go_logger "github.com/pefish/go-logger"
	vm2 "github.com/pefish/telegram-bot-manager/pkg/vm"
	"github.com/pefish/telegram-bot-manager/version"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
)

type DefaultCommand struct {
	offsetFileFs *os.File
}

func NewDefaultCommand() *DefaultCommand {
	return &DefaultCommand{}
}

func (s *DefaultCommand) DecorateFlagSet(flagSet *flag.FlagSet) error {
	flagSet.String("data-dir", os.ExpandEnv("$HOME/.") + version.AppName, "set data dictionary")
	return nil
}

func (s *DefaultCommand) OnExited() error {
	err := s.offsetFileFs.Sync()
	if err != nil {
		return err
	}
	err = s.offsetFileFs.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *DefaultCommand) Start() error {
	commandsJsFile := go_config.Config.MustGetString("commandsJsFile")
	fs, err := os.Open(commandsJsFile)
	if err != nil {
		return err
	}
	defer fs.Close()
	scriptBytes, err := ioutil.ReadAll(fs)
	if err != nil {
		return err
	}

	vm := vm2.NewVm()
	_, err = vm.RunString(string(scriptBytes) + "\n" + `
function execute(command, args) {
	if (!commands[command]) {
		return "Sorry, I don't understand."
	}
    return commands[command](args)
}
`)
	if err != nil {
		return err
	}

	var fn func(string, []string) string
	err = vm.ExportTo(vm.Get("execute"), &fn)
	if err != nil {
		panic(err)
	}

	token := go_config.Config.MustGetString("token")
	timer := time.NewTimer(0)
	// load offset
	var offsetStr string = "0"
	dataDir := go_config.Config.MustGetString("data-dir")
	// create it if dataDir not exist
	info, err := os.Stat(dataDir)
	if err != nil || !info.IsDir() {
		err = os.MkdirAll(dataDir, 0755)
		if err != nil {
			return err
		}
	}
	offsetFilename := path.Join(dataDir, "./offset")
	s.offsetFileFs, err = os.OpenFile(offsetFilename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	offsetBytes, err := ioutil.ReadAll(s.offsetFileFs)
	if err != nil {
		return err
	}
	if len(offsetBytes) != 0 {
		offsetStr = string(offsetBytes)
	}

	go_logger.Logger.InfoF("current offset: %s", offsetStr)
	type GetUpdatesResult struct {
		Ok     bool `json:"ok"`
		Result []struct {
			UpdateId uint64 `json:"update_id"`
			Message  struct {
				MessageId uint64 `json:"message_id"`
				From      struct {
					Id        uint64 `json:"id"`
					IsBot     bool   `json:"is_bot"`
					FirstName string `json:"first_name"`
					Username  string `json:"username"`
				} `json:"from"`
				Chat struct {
					Id        int64 `json:"id"`
					FirstName string `json:"first_name"`
					Username  string `json:"username"`
					Type      string `json:"type"`
				} `json:"chat"`
				Date     uint64 `json:"date"`
				Text     string `json:"text"`
				Entities []struct {
					Offset uint64 `json:"offset"`
					Length uint64 `json:"length"`
					Type   string `json:"type"`
				} `json:"entities"`
			} `json:"message"`
		} `json:"result"`
	}
	for range timer.C {
		var getUpdatesResult GetUpdatesResult
		_, err := go_http.NewHttpRequester(go_http.WithLogger(go_logger.Logger)).GetForStruct(go_http.RequestParam{
			Url: fmt.Sprintf("https://api.telegram.org/bot%s/getUpdates?offset=%s&limit=10", token, offsetStr),
		}, &getUpdatesResult)
		if err != nil {
			go_logger.Logger.Error(err)
			timer.Reset(2 * time.Second)
			continue
		}
		go_logger.Logger.Debug(getUpdatesResult)
		if !getUpdatesResult.Ok {
			go_logger.Logger.Error("getUpdatesResult.Ok not true")
			timer.Reset(2 * time.Second)
			continue
		}
		if len(getUpdatesResult.Result) == 0 {
			go_logger.Logger.Info("no updates")
			timer.Reset(2 * time.Second)
			continue
		}
		go_logger.Logger.InfoF("-- start to process %d updates", len(getUpdatesResult.Result))
		for _, result := range getUpdatesResult.Result {
			// change offset
			offsetStr = go_decimal.Decimal.Start(result.UpdateId).AddForString(1)
			_, err = s.offsetFileFs.WriteAt([]byte(offsetStr), 0)
			if err != nil {
				go_logger.Logger.Error(err)
				continue
			}
			// decode command
			commandText := result.Message.Text
			commandTextArr := strings.Split(commandText, " ")
			// execute command
			executeResult := fn(commandTextArr[0], commandTextArr[1:])
			// ack
			go_logger.Logger.InfoF("---- process command: %s", commandText)
			go_logger.Logger.InfoF("---- update_id: %d", result.UpdateId)
			var sendMessageResult struct{
				Ok bool `json:"ok"`
				ErrorCode uint64 `json:"error_code"`
				Description string `json:"description"`
			}
			_, err := go_http.NewHttpRequester(go_http.WithLogger(go_logger.Logger)).GetForStruct(go_http.RequestParam{
				Url:       fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%d&text=%s", token, result.Message.Chat.Id, executeResult),
			}, &sendMessageResult)
			if err != nil {
				go_logger.Logger.Error(err)
				continue
			}
		}
		timer.Reset(time.Second)
	}
	return nil
}

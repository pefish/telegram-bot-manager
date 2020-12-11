package telegram_robot

import (
	"errors"
	"fmt"
	go_decimal "github.com/pefish/go-decimal"
	go_error "github.com/pefish/go-error"
	go_http "github.com/pefish/go-http"
	go_logger "github.com/pefish/go-logger"
	telegram_sender "github.com/pefish/telegram-bot-manager/pkg/telegram-sender"
	vm2 "github.com/pefish/telegram-bot-manager/pkg/vm"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

type Robot struct {
	commandsStr string
	token string
	offsetFileFs *os.File
	telegramSender *telegram_sender.TelegramSender
}

func (r *Robot) TelegramSender() *telegram_sender.TelegramSender {
	return r.telegramSender
}

/**
**commandsStr**
var commands = {
    "/test": {
        desc: "测试命令",
        func: function (args) {
            // console.log(args)
            return "test: " + JSON.stringify(args)
        }
    },
    "/haha": {
        desc: "有点意思",
        func: function (args) {
            return "xixi"
        }
    },
}
*/
func NewRobot(commandsStr, token string) *Robot {
	return &Robot{
		commandsStr: commandsStr,
		token: token,
	}
}

func (r *Robot) Close() error {
	err := r.offsetFileFs.Sync()
	if err != nil {
		return err
	}
	err = r.offsetFileFs.Close()
	if err != nil {
		return err
	}
	return nil
}


func (r *Robot) Start(dataDir string) error {
	vm := vm2.NewVm()
	_, err := vm.RunString(r.commandsStr + "\n" + `
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

	timer := time.NewTimer(0)
	// load offset
	var offsetStr string = "0"
	// create it if dataDir not exist
	info, err := os.Stat(dataDir)
	if err != nil || !info.IsDir() {
		err = os.MkdirAll(dataDir, 0755)
		if err != nil {
			return go_error.WithStack(err)
		}
	}
	offsetFilename := path.Join(dataDir, "./offset")
	r.offsetFileFs, err = os.OpenFile(offsetFilename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return go_error.WithStack(err)
	}
	offsetBytes, err := ioutil.ReadAll(r.offsetFileFs)
	if err != nil {
		return go_error.WithStack(err)
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

	r.telegramSender = telegram_sender.NewTelegramSender(r.token)
	r.telegramSender.SetLogger(go_logger.Logger)

	for range timer.C {
		var getUpdatesResult GetUpdatesResult
		_, err := go_http.NewHttpRequester(go_http.WithLogger(go_logger.Logger)).GetForStruct(go_http.RequestParam{
			Url: fmt.Sprintf("https://api.telegram.org/bot%s/getUpdates?offset=%s&limit=10", r.token, offsetStr),
		}, &getUpdatesResult)
		if err != nil {
			go_logger.Logger.Error(go_error.WithStack(err))
			timer.Reset(2 * time.Second)
			continue
		}
		go_logger.Logger.Debug(getUpdatesResult)
		if !getUpdatesResult.Ok {
			go_logger.Logger.Error(errors.New("getUpdatesResult.Ok not true"))
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
			_, err = r.offsetFileFs.WriteAt([]byte(offsetStr), 0)
			if err != nil {
				go_logger.Logger.Error(go_error.WithStack(err))
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
			r.telegramSender.SendMsg(telegram_sender.MsgStruct{
				ChatId: result.Message.Chat.Id,
				Msg:    []byte(url.QueryEscape(executeResult)),
			}, 0)
		}
		timer.Reset(time.Second)
	}
	return nil
}

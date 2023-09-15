package telegram_robot

import (
	"context"
	"errors"
	"fmt"
	go_decimal "github.com/pefish/go-decimal"
	go_error "github.com/pefish/go-error"
	go_format "github.com/pefish/go-format"
	go_http "github.com/pefish/go-http"
	go_logger "github.com/pefish/go-logger"
	telegram_sender "github.com/pefish/telegram-bot-manager/pkg/telegram-sender"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
)

type Robot struct {
	token          string
	loopInterval   time.Duration
	offsetFileFs   *os.File
	telegramSender *telegram_sender.TelegramSender
	logger         go_logger.InterfaceLogger
}

func (r *Robot) TelegramSender() *telegram_sender.TelegramSender {
	return r.telegramSender
}

func NewRobot(token string, loopInterval time.Duration) *Robot {
	telegramSender := telegram_sender.NewTelegramSender(token)
	return &Robot{
		token:          token,
		telegramSender: telegramSender,
		loopInterval:   loopInterval,
	}
}

func (r *Robot) SetLogger(logger go_logger.InterfaceLogger) {
	r.logger = logger
	r.telegramSender.SetLogger(logger)
}

func (r *Robot) Close() error {
	if r.offsetFileFs != nil {
		err := r.offsetFileFs.Sync()
		if err != nil {
			return err
		}
		err = r.offsetFileFs.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Robot) Start(ctx context.Context, dataDir string, processCmdFn func(string, string) (string, error)) error {
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
	r.logger.Debug(offsetFilename)
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

	r.logger.DebugF("current offset: %s", offsetStr)
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
					Id        int64  `json:"id"`
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
over:
	for {
		select {
		case <-timer.C:
			var getUpdatesResult GetUpdatesResult
			_, _, err := go_http.NewHttpRequester(
				go_http.WithLogger(r.logger),
				go_http.WithTimeout(20*time.Second),
			).GetForStruct(go_http.RequestParam{
				Url: fmt.Sprintf("https://api.telegram.org/bot%s/getUpdates?offset=%s&limit=10", r.token, offsetStr),
			}, &getUpdatesResult)
			if err != nil {
				r.logger.Error(go_error.WithStack(err))
				timer.Reset(r.loopInterval)
				continue
			}
			r.logger.Debug(getUpdatesResult)
			if !getUpdatesResult.Ok {
				r.logger.Error(errors.New("getUpdatesResult.Ok not true"))
				timer.Reset(r.loopInterval)
				continue
			}
			if len(getUpdatesResult.Result) == 0 {
				r.logger.Debug("no updates")
				timer.Reset(r.loopInterval)
				continue
			}
			r.logger.DebugF("-- start to process %d updates", len(getUpdatesResult.Result))
			for _, result := range getUpdatesResult.Result {
				// change offset
				offsetStr = go_decimal.Decimal.Start(result.UpdateId).AddForString(1)
				_, err = r.offsetFileFs.WriteAt([]byte(offsetStr), 0)
				if err != nil {
					r.logger.Error(go_error.WithStack(err))
					continue
				}
				r.logger.DebugF("---- process msg: %s", result.Message.Text) // 是整个消息，如 /test hhh
				r.logger.DebugF("---- update_id: %d", result.UpdateId)
				commandTextArr := strings.Split(result.Message.Text, " ")
				processResult, err := processCmdFn(commandTextArr[0], strings.Join(commandTextArr[1:], ""))
				if err != nil {
					r.logger.Error(go_error.WithStack(err))
					continue
				}
				if processResult == "" {
					continue
				}
				// ack
				r.telegramSender.SendMsg(telegram_sender.MsgStruct{
					ChatId: go_format.FormatInstance.ToString(result.Message.Chat.Id),
					Msg:    processResult,
					Ats:    []string{result.Message.From.Username},
				}, 0)
			}
			timer.Reset(r.loopInterval)
		case <-ctx.Done():
			break over
		}
	}
	return nil
}

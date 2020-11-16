package telegram_sender

import (
	"fmt"
	go_http "github.com/pefish/go-http"
	"github.com/pefish/go-interface-logger"
	go_logger "github.com/pefish/go-logger"
	"github.com/pkg/errors"
	"sync"
)

type MsgStruct struct {
	ChatId int64
	Msg    []byte
}

type TelegramSender struct {
	msgs        []MsgStruct
	msgLock     sync.Mutex
	msgReceived chan bool
	token       string
	logger      go_interface_logger.InterfaceLogger
}

func NewTelegramSender(token string) *TelegramSender {
	ts := &TelegramSender{
		msgs:        make([]MsgStruct, 0, 10),
		token:       token,
		logger:      go_interface_logger.DefaultLogger,
		msgReceived: make(chan bool),
	}

	go func() {
		for {
			for _, msg := range ts.msgs {
				go func(msg MsgStruct) {
					err := ts.send(msg.ChatId, string(msg.Msg))
					if err != nil {
						ts.logger.Error(err)
					}
				}(msg)
			}
			ts.msgLock.Lock()
			ts.msgs = make([]MsgStruct, 0, 10)
			ts.msgLock.Unlock()
			select {
			case <-ts.msgReceived:
				ts.logger.Debug("notify received")
			}
			ts.logger.Debug("to send...")
		}
	}()

	return ts
}

func (ts *TelegramSender) SetLogger(logger go_interface_logger.InterfaceLogger) {
	ts.logger = logger
}

func (ts *TelegramSender) SendMsg(msg MsgStruct) {
	ts.msgLock.Lock()
	ts.msgs = append(ts.msgs, msg)
	ts.msgLock.Unlock()
	// try to notify
	select {
	case ts.msgReceived <- true:
		ts.logger.Debug("notify succeed")
	default:
		ts.logger.Debug("no need to notify")
	}
}

func (ts *TelegramSender) send(chatId int64, text string) error {
	var sendMessageResult struct {
		Ok          bool   `json:"ok"`
		ErrorCode   uint64 `json:"error_code"`
		Description string `json:"description"`
	}
	_, err := go_http.NewHttpRequester(go_http.WithLogger(go_logger.Logger)).GetForStruct(go_http.RequestParam{
		Url: fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%d&text=%s", ts.token, chatId, text),
	}, &sendMessageResult)
	if err != nil {
		return err
	}
	if !sendMessageResult.Ok {
		return errors.New(sendMessageResult.Description)
	}
	return nil
}

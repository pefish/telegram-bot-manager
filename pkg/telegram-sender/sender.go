package telegram_sender

import (
	"encoding/json"
	"fmt"
	go_error "github.com/pefish/go-error"
	go_http "github.com/pefish/go-http"
	go_logger "github.com/pefish/go-logger"
	"github.com/pkg/errors"
	"net/url"
	"sync"
	"time"
)

type MsgStruct struct {
	ChatId string   `json:"chat_id"`
	Msg    string   `json:"msg"`
	Ats    []string `json:"ats"`
}

type TelegramSender struct {
	msgs        []MsgStruct
	msgLock     sync.Mutex
	msgReceived chan bool
	token       string
	logger      go_logger.InterfaceLogger

	lastSend      map[string]time.Time
	httpRequester go_http.IHttp
}

func NewTelegramSender(token string) *TelegramSender {
	ts := &TelegramSender{
		msgs:          make([]MsgStruct, 0, 10),
		token:         token,
		logger:        go_logger.Logger,
		msgReceived:   make(chan bool),
		lastSend:      make(map[string]time.Time, 10),
		httpRequester: go_http.NewHttpRequester(go_http.WithLogger(go_logger.Logger), go_http.WithTimeout(20*time.Second)),
	}

	go func() {
		for {
			for _, msg := range ts.msgs {
				go func(msg MsgStruct) {
					if msg.Ats != nil && len(msg.Ats) > 0 {
						for _, at := range msg.Ats {
							msg.Msg += " @" + at
						}
					}
					err := ts.send(msg.ChatId, url.QueryEscape(msg.Msg))
					if err != nil {
						ts.logger.Error(go_error.WithStack(err))
						return
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

func (ts *TelegramSender) SetLogger(logger go_logger.InterfaceLogger) *TelegramSender {
	ts.logger = logger
	return ts
}

// interval: interval间隔内不发送
func (ts *TelegramSender) SendMsg(msg MsgStruct, interval time.Duration) error {
	mar, err := json.Marshal(msg)
	if err != nil {
		return go_error.WithStack(err)
	}
	if lastTime, ok := ts.lastSend[string(mar)]; ok && time.Now().Sub(lastTime) < interval {
		return errors.New("trigger interval")
	}
	ts.lastSend[string(mar)] = time.Now()

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
	return nil
}

func (ts *TelegramSender) send(chatId string, text string) error {
	var sendMessageResult struct {
		Ok          bool   `json:"ok"`
		ErrorCode   uint64 `json:"error_code"`
		Description string `json:"description"`
	}
	_, _, err := ts.httpRequester.GetForStruct(go_http.RequestParam{
		Url: fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s", ts.token, chatId, text),
	}, &sendMessageResult)
	if err != nil {
		return go_error.WithStack(err)
	}
	if !sendMessageResult.Ok {
		return errors.New(sendMessageResult.Description)
	}
	return nil
}

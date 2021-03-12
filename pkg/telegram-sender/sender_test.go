package telegram_sender

import (
	"encoding/json"
	"github.com/golang/mock/gomock"
	go_http "github.com/pefish/go-http"
	"github.com/pefish/go-test-assert"
	mock_go_http "github.com/pefish/telegram-bot-manager/mock/mock-go-http"
	"net/http"
	"testing"
	"time"
)

func TestTelegramSender_SendMsg(t *testing.T) {
	successCount := 0

	ctrl := gomock.NewController(t)
	instance := mock_go_http.NewMockIHttp(ctrl)
	instance.EXPECT().GetForStruct(gomock.Any(), gomock.Any()).DoAndReturn(func(param go_http.RequestParam, struct_ interface{}) (*http.Response, error) {
		test.Equal(t, "https://api.telegram.org/bot/sendMessage?chat_id=111&text=haha", param.Url)
		successCount++
		err := json.Unmarshal([]byte(`{"ok": true}`), struct_)
		test.Equal(t, nil, err)
		return nil, nil
	}).AnyTimes()

	sender := NewTelegramSender("")
	sender.httpRequester = instance
	err := sender.SendMsg(MsgStruct{
		ChatId: 111,
		Msg:    []byte("haha"),
	}, 10 * time.Second)
	test.Equal(t, nil, err)
	err = sender.SendMsg(MsgStruct{
		ChatId: 111,
		Msg:    []byte("haha"),
	}, 10 * time.Second)
	test.Equal(t,  true, err != nil)
	test.Equal(t, "trigger interval", err.Error())
	err = sender.SendMsg(MsgStruct{
		ChatId: 111,
		Msg:    []byte("haha"),
	}, 0)
	test.Equal(t, nil, err)

	time.Sleep(1 * time.Second)
	test.Equal(t, 2, successCount)
}

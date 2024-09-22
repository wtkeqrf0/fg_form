package api

import (
	"context"
	"fmt"
	tgBotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/wtkeqrf0/tg_form/internal/appeal"
	"github.com/wtkeqrf0/tg_form/pkg/util"
	"log"
	"strings"
	"time"
)

//go:generate ifacemaker -f *.go -o api_if.go -i Api -s tgBot -p api
type tgBot struct {
	api    *tgBotAPI.BotAPI
	appeal appeal.Method
}

func New(ctx context.Context, token string) Api {
	bot, err := tgBotAPI.NewBotAPI(token)
	if err != nil {
		panic(err)
	}
	context.AfterFunc(ctx, bot.StopReceivingUpdates)
	return &tgBot{api: bot}
}

func (t *tgBot) Start(appeal appeal.Method) {
	t.appeal = appeal
	updChan := t.api.GetUpdatesChan(tgBotAPI.UpdateConfig{
		Timeout: 30,
	})

	read := func(updChan <-chan tgBotAPI.Update) {
		for update := range updChan {
			if chatID, err := t.process(update); err != nil {
				log.Println(err.Error())
				_, _ = t.api.Send(tgBotAPI.NewMessage(chatID, err.Error()))
			}
		}
	}

	for i := 0; i < 5; i++ {
		go read(updChan)
	}
	read(updChan)
}

func (t *tgBot) process(update tgBotAPI.Update) (chatID int64, err error) {
	defer func() {
		if r := recover(); r != nil {
			if err, _ = r.(error); err == nil {
				err = fmt.Errorf("%v", r)
			}
		}
	}()

	var (
		msgs        = make([]tgBotAPI.Chattable, 1)
		tgMsg       *tgBotAPI.Message
		ctx, cancel = context.WithTimeout(context.Background(), time.Minute*20)
	)
	defer cancel()

	switch {
	case update.Message != nil:
		tgMsg = update.Message
		chatID = tgMsg.Chat.ID

		switch update.Message.Command() {
		case "start":
			msgs[0], err = t.SendMe(chatID)
		case "new":
			msgs[0], err = t.appeal.ChooseDivisions(ctx, chatID)
		default:
			text := strings.TrimSpace(tgMsg.Text)
			if len(text) == 0 {
				return
			}
			msgs[0], err = t.appeal.StrForm(ctx, chatID, text)
		}
	case update.CallbackQuery != nil:
		tgMsg = update.CallbackQuery.Message
		chatID = tgMsg.Chat.ID

		key, value := util.ParseCallBack(update.CallbackQuery.Data)
		switch key {
		case util.DivisionCallBackKey:
			msgs[0], err = t.appeal.SelectDivision(ctx, chatID, tgMsg.MessageID, value)
		case util.SubmitCallBackKey:
			msgs, err = t.appeal.SendAppeal(ctx, chatID, tgMsg.MessageID, tgMsg.Chat.UserName, value)
		case util.FastReplyCallBackKey:
			msgs, err = t.appeal.SendFastReply(ctx, chatID, tgMsg.MessageID, value)
		}
	}

	if err != nil {
		return
	}

	for _, msg := range msgs {
		if msg != nil {
			if _, err = t.api.Send(msg); err != nil {
				return
			}
		}
	}
	return
}

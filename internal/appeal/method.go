package appeal

import (
	"context"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/redis/go-redis/v9"
	"github.com/wtkeqrf0/tg_form/internal/session"
	"github.com/wtkeqrf0/tg_form/pkg/util"
	"strconv"
	"strings"
)

type method struct {
	repo    Repository
	session session.Session
}

//go:generate ifacemaker -f method.go -o method_if.go -i Method -s method -p appeal
func NewMethod(repo Repository, session session.Session) Method {
	return &method{
		repo:    repo,
		session: session,
	}
}

func (m *method) ChooseDivisions(ctx context.Context, chatID int64) (tgbotapi.Chattable, error) {
	if err := m.session.Delete(ctx, fmt.Sprint(chatID)); err != nil {
		return nil, err
	}

	row := make([]tgbotapi.InlineKeyboardButton, len(util.Divisions))
	for i := range util.Divisions {
		row[i] = tgbotapi.NewInlineKeyboardButtonData(util.Divisions[i], fmt.Sprintf("%s/%d", util.DivisionCallBackKey, i+1))
	}

	msg := tgbotapi.NewMessage(chatID, "Выберете подразделение")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(row)
	return msg, nil
}

func (m *method) SelectDivision(ctx context.Context, chatID int64, messageID int, division string) (tgbotapi.Chattable, error) {
	err := m.session.SetValue(ctx, fmt.Sprint(chatID), util.DivisionField, division)
	if err != nil {
		return nil, err
	}
	return tgbotapi.NewEditMessageText(chatID, messageID, "Введите заголовок (тему) обращения"), nil
}

func (m *method) StrForm(ctx context.Context, chatID int64, text string) (tgbotapi.Chattable, error) {
	payload, err := m.session.Get(ctx, fmt.Sprint(chatID))
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	if payload.Division == 0 {
		return m.ChooseDivisions(ctx, chatID)
	} else if len(payload.Subject) == 0 {
		return m.setSubject(ctx, chatID, text)
	} else if len(payload.Text) == 0 {
		return m.setText(ctx, chatID, text)
	}
	return nil, nil
}

func (m *method) setSubject(ctx context.Context, chatID int64, text string) (tgbotapi.Chattable, error) {
	chStr := fmt.Sprint(chatID)

	err := m.session.SetValue(ctx, chStr, util.SubjectField, text)
	if err != nil {
		return nil, err
	}
	return tgbotapi.NewMessage(chatID, "Введите текст обращения"), nil
}

func (m *method) setText(ctx context.Context, chatID int64, text string) (tgbotapi.Chattable, error) {
	chStr := fmt.Sprint(chatID)
	payload, err := m.session.Get(ctx, chStr)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return m.ChooseDivisions(ctx, chatID)
		}
		return nil, err
	}
	payload.Text = text

	msgID := util.GenerateString(5)
	if err = m.session.Set(ctx, fmt.Sprintf("%s:%s", chStr, msgID), payload); err != nil {
		return nil, err
	}

	if err = m.session.Delete(ctx, chStr); err != nil {
		return nil, err
	}

	resp := new(strings.Builder)
	resp.WriteString("Вот, что у вас вышло:\n\n")
	payload.ToString(resp)

	msg := tgbotapi.NewMessage(chatID, resp.String())
	msg.ReplyMarkup = tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
			{
				tgbotapi.NewInlineKeyboardButtonData("Отправить",
					fmt.Sprintf("%s/%s", util.SubmitCallBackKey, msgID)),
			},
		},
	}
	msg.ParseMode = tgbotapi.ModeHTML
	return msg, nil
}

func (m *method) SendAppeal(ctx context.Context, chatID int64, messageID int, userName, value string) ([]tgbotapi.Chattable, error) {
	key := fmt.Sprintf("%d:%s", chatID, value)
	payload, err := m.session.Get(ctx, key)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return []tgbotapi.Chattable{
				tgbotapi.NewEditMessageText(chatID, messageID, "Сожаеем!\nОбращение устарело."),
			}, nil
		}
		return nil, err
	}

	saved, err := m.repo.SaveAppeal(ctx, &Appeal{
		Division: util.Divisions[payload.Division-1],
		Subject:  payload.Subject,
		Text:     payload.Text,
		ChatID:   chatID,
		Username: userName,
	})
	if err != nil {
		return nil, err
	}

	if err = m.session.Delete(ctx, key); err != nil {
		return nil, err
	}

	resp := new(strings.Builder)
	resp.WriteString(fmt.Sprintf("Вам пришло новое обращение от @%s:\n\n", userName))
	payload.ToString(resp)

	adminMsg := tgbotapi.NewMessage(saved.AdminChatID, resp.String())
	adminMsg.ParseMode = tgbotapi.ModeHTML
	adminMsg.ReplyMarkup = tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
			{
				tgbotapi.NewInlineKeyboardButtonData("Быстрый ответ", fmt.Sprintf("%s/%d", util.FastReplyCallBackKey, saved.AppealID)),
			},
		},
	}

	return []tgbotapi.Chattable{
		adminMsg, tgbotapi.NewEditMessageText(chatID, messageID,
			"Обращение успешно отправлено!\nВ ближайшее время вам ответит наш администратор"),
	}, nil
}

func (m *method) SendFastReply(ctx context.Context, chatID int64, messageID int, appealID string) ([]tgbotapi.Chattable, error) {
	id, err := strconv.ParseInt(appealID, 10, 64)
	if err != nil {
		return nil, err
	}

	appeal, err := m.repo.GetAppeal(ctx, id)
	if err != nil {
		return nil, err
	}

	if err = m.repo.AnswerAppeal(ctx, id); err != nil {
		return nil, err
	}

	return []tgbotapi.Chattable{
		tgbotapi.NewMessage(appeal.ChatID, "Hello world!"),
		tgbotapi.NewEditMessageText(chatID, messageID,
			fmt.Sprintf("Ответ отправлен пользователю %s", appeal.Username)),
	}, nil
}

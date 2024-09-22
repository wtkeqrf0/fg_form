package api

import (
	"fmt"
	tgBotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

func (t *tgBot) SendMe(chatID int64) (tgBotAPI.Chattable, error) {
	cmds, err := t.api.GetMyCommands()
	if err != nil {
		return nil, err
	}
	txt := new(strings.Builder)
	for _, cmd := range cmds {
		txt.WriteString(fmt.Sprintf("\n/%s - %s", cmd.Command, cmd.Description))
	}

	return tgBotAPI.NewMessage(chatID, txt.String()), nil
}

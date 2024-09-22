package session

import (
	"fmt"
	"github.com/wtkeqrf0/tg_form/pkg/util"
	"strings"
)

type Payload struct {
	Division util.Division `redis:"division"`
	Subject  string        `redis:"subject"`
	Text     string        `redis:"text"`
}

func (p *Payload) ToString(dst *strings.Builder) {
	const tmpl = "<strong>Подразделение:</strong> %s" +
		"\n<strong>Тема:</strong> %s" +
		"\n<strong>Текст:</strong> %s"
	dst.WriteString(fmt.Sprintf(tmpl, util.Divisions[p.Division-1], p.Subject, p.Text))
}

package LocalMail

import (
	"fmt"
	"strings"
)

type Mail struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func (m *Mail) String() string {
	builder := strings.Builder{}
	fmt.Fprintf(&builder, "From: %s\n", m.From)
	fmt.Fprintf(&builder, "To: %s\n", m.To)
	fmt.Fprintf(&builder, "Subject: %s\n", m.Subject)
	fmt.Fprint(&builder, "\n\n")
	fmt.Fprintf(&builder, "%s\n", m.Body)
	return builder.String()
}

func GetTestMessage(body string) *Mail {
	return &Mail{
		From:    "Jannes Morning Briefing <morning@jskweb.de>",
		To:      "jannes@jskweb.de",
		Subject: "Morning Briefing Teste",
		Body:    body,
	}
}

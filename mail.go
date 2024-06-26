package LocalMail

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/mail"
	"strings"
)

func MessageToString(msg *mail.Message) (string, error) {
	var buf bytes.Buffer

	// Write the headers
	for header, values := range msg.Header {
		for _, value := range values {
			_, err := fmt.Fprintf(&buf, "%s: %s\r\n", header, value)
			if err != nil {
				return "", err
			}
		}
	}

	// Write a blank line to separate headers from the body
	buf.WriteString("\r\n")

	// Write the body
	body, err := io.ReadAll(msg.Body)
	if err != nil {
		return "", err
	}
	buf.Write(body)
	return buf.String(), nil
}

func GetTestMessage(body string) *mail.Message {
	message := mail.Message{}

	message.Header = mail.Header{
		"From":       []string{"Jannes' Morning Briefing <morning@jskweb.de>"},
		"To":         []string{"jannes@jskweb.de"},
		"Subject":    []string{"Morning Briefing Test"},
		"Message-ID": []string{uuid.New().String()},
	}

	message.Body = strings.NewReader(body)
	return &message
}

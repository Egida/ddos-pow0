package challenge_resp

import (
	"fmt"
	"strings"
)

const (
	QUIT               = "quit"
	REQUEST_CHALLENGE  = "request-challenge"
	RESPONSE_CHALLENGE = "response-challenge"
	REQUEST_RESOURCE   = "request-resource"
	RESPONSE_RESOURCE  = "response_resource"
)

type Message struct {
	Header  string
	Payload string
}

func (m *Message) Stringify() string {
	return fmt.Sprintf("%s|%s", m.Header, m.Payload)
}

func ParseMessage(str string) (*Message, error) {
	str = strings.TrimSpace(str)

	parts := strings.Split(str, "|")
	if len(parts) < 1 || len(parts) > 2 {
		return nil, fmt.Errorf("message doesn't match protocol")
	}

	msg := Message{
		Header: parts[0],
	}
	if len(parts) == 2 {
		msg.Payload = parts[1]
	}
	return &msg, nil
}

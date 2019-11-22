package game

import "github.com/qystishere/survemu/pkg/bitbuf"

type Message struct {
	ID byte
}

type ReceivedMessage struct {
	*Message
	*bitbuf.Reader
}

type SentMessage struct {
	*Message
	*bitbuf.Writer
}

func NewSentMessage(id byte, length int) *SentMessage {
	return &SentMessage{
		Message: &Message{
			ID: id,
		},
		Writer: bitbuf.NewWriter(length),
	}
}

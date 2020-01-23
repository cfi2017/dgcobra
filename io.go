package dgcobra

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

type MessageWriter struct {
	Session *discordgo.Session
	Channel string
}

func NewMessageWriter(session *discordgo.Session, channel string) *MessageWriter {
	return &MessageWriter{
		Session: session,
		Channel: channel,
	}
}

func (w *MessageWriter) Write(b []byte) (n int, err error) {
	m := strings.TrimSpace(string(b))
	if m == "" {
		return 0, nil
	}
	msg, err := w.Session.ChannelMessageSend(w.Channel, string(b))
	if err != nil {
		return 0, err
	}
	return len(msg.Content), err
}

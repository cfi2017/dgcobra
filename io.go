package dgcobra

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

// MessageWriter is a writer implementation for a discord channel. This is used for example by cobra to output usage and help.
type MessageWriter struct {
	Session *discordgo.Session
	Channel string
}

// NewMessageWriter creates a message writer based on a given session and channel. Each write calls ChannelMessageSend on the session.
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
	msg, err := w.Session.ChannelMessageSend(w.Channel, m)
	if err != nil {
		return 0, err
	}
	return len(msg.Content), err
}

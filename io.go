package dgcobra

import (
	"strings"
	"sync"

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

type Buffer struct {
	sync.RWMutex
	buffer string
}

func (b *Buffer) Append(s string) {
	b.Lock()
	defer b.Unlock()
	b.buffer += s
}

func (b *Buffer) Clear() {
	b.Lock()
	defer b.Unlock()
	b.buffer = ""
}

func (b *Buffer) Flush() string {
	b.Lock()
	defer b.Unlock()
	m := b.buffer
	b.buffer = ""
	return m
}

type BufferedMessageWriter struct {
	sync.RWMutex
	Session *discordgo.Session
	Channel string
	buffer  *Buffer
}

func NewBufferedMessageWriter(session *discordgo.Session, channel string) *BufferedMessageWriter {
	return &BufferedMessageWriter{
		Session: session,
		Channel: channel,
		buffer:  &Buffer{},
	}
}

func (w *BufferedMessageWriter) Write(b []byte) (n int, err error) {
	m := strings.TrimSpace(string(b))
	if m == "" {
		return 0, nil
	}
	w.buffer.Append(m + "\n")
	return len(m) + 1, nil
}

func (w *BufferedMessageWriter) Flush() (n int, err error) {
	m := w.buffer.Flush()
	if m == "" {
		return 0, nil
	}
	msg, err := w.Session.ChannelMessageSend(w.Channel, m)
	if err != nil {
		return 0, err
	}
	return len(msg.Content), err
}

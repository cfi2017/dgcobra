package dgcobra

import (
	"encoding/csv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/cobra"
)

type ErrorInvalidArgs struct {
	Message string
	Err     error
}

func (err ErrorInvalidArgs) Error() string {
	return err.Message
}

func (err ErrorInvalidArgs) Unwrap() error {
	return err.Err
}

type Handler struct {
	RootFactory func(session *discordgo.Session, event *discordgo.MessageCreate) *cobra.Command
	session     *discordgo.Session
	Prefixes    []string
	PrefixFunc  func(session *discordgo.Session, event *discordgo.MessageCreate) []string
	ErrFunc     func(err error)
}

func NewHandler(session *discordgo.Session) *Handler {
	return &Handler{
		session: session,
	}
}

func (h *Handler) AddPrefix(prefix string) {
	h.Prefixes = append(h.Prefixes, prefix)
}

func (h *Handler) Start() {
	h.session.AddHandler(func(_ *discordgo.Session, event *discordgo.MessageCreate) {
		prefixes := h.Prefixes
		if h.PrefixFunc != nil {
			prefixes = append(prefixes, h.PrefixFunc(h.session, event)...)
		}
		for _, prefix := range prefixes {
			if strings.HasPrefix(event.Content, prefix) {
				args, err := parseArgs(strings.TrimPrefix(event.Content, prefix))
				if err != nil && h.ErrFunc != nil {
					h.ErrFunc(ErrorInvalidArgs{Err: err, Message: "couldn't parse args"})
					return
				}

				w := NewMessageWriter(h.session, event.ChannelID)
				// get commands
				root := h.RootFactory(h.session, event)
				root.SetArgs(args)
				root.SetOut(w)
				err = root.Execute()
				if err != nil && h.ErrFunc != nil {
					h.ErrFunc(ErrorInvalidArgs{Err: err, Message: "couldn't execute command"})
					return
				}
				return
			}
		}
	})
}

func parseArgs(argString string) ([]string, error) {
	r := csv.NewReader(strings.NewReader(argString))
	r.Comma = ' ' // space
	return r.Read()
}

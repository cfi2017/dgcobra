package dgcobra

import (
	"encoding/csv"
	"io"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/cobra"
)

// ErrorInvalidArgs indicates invalid arguments were passed in a command. You can call Unwrap() to get the underlying error.
type ErrorInvalidArgs struct {
	Session *discordgo.Session
	Event   *discordgo.MessageCreate
	Message string
	Err     error
}

func (err ErrorInvalidArgs) Error() string {
	return err.Message
}

// Unwrap returns the underlying error behind ErrorInvalidArgs.
func (err ErrorInvalidArgs) Unwrap() error {
	return err.Err
}

/*
Handler represents a dgcobra Command Handler. This builds upon a RootCommandFactory.

To use this, use NewHandler() and add a prefix and a root command factory. Then call Handler.Start().
Examples in examples folder of this repository.
*/
type Handler struct {
	// Root command factory for the bot. This needs to be set.
	RootFactory func(session *discordgo.Session, event *discordgo.MessageCreate) *cobra.Command
	session     *discordgo.Session
	// List of global prefixes for the bot.
	Prefixes []string
	// Function to load prefixes for a specific message. Use this to allow guild-specific prefixes.
	PrefixFunc                 func(session *discordgo.Session, event *discordgo.MessageCreate) []string
	PrefixFuncOverridesGlobals bool
	// Function that is called when the message event errors for some reason.
	ErrFunc          func(err error)
	OutWriterFactory func(s *discordgo.Session, channel string) io.Writer
}

// NewHandler creates a new handler with a given session.
func NewHandler(session *discordgo.Session) *Handler {
	return &Handler{
		session: session,
		OutWriterFactory: func(s *discordgo.Session, channel string) io.Writer {
			return NewBufferedMessageWriter(s, channel)
		},
		PrefixFuncOverridesGlobals: true,
	}
}

// AddPrefix registers a new global Prefix
func (h *Handler) AddPrefix(prefix string) {
	h.Prefixes = append(h.Prefixes, prefix)
}

// Start registers a new handler with discordgo and starts receiving commands. This function is non-blocking.
func (h *Handler) Start() {
	h.session.AddHandler(func(_ *discordgo.Session, event *discordgo.MessageCreate) {
		prefixes := h.Prefixes
		if h.PrefixFunc != nil {
			additionals := h.PrefixFunc(h.session, event)
			if len(additionals) > 0 && h.PrefixFuncOverridesGlobals {
				prefixes = h.PrefixFunc(h.session, event)
			} else {
				prefixes = append(prefixes, h.PrefixFunc(h.session, event)...)
			}
		}
		for _, prefix := range prefixes {
			if strings.HasPrefix(event.Content, prefix) {
				a := strings.TrimSpace(strings.TrimPrefix(event.Content, prefix))
				if a == "" {
					return
				}
				args, err := parseArgs(a)
				if err != nil && h.ErrFunc != nil {
					h.ErrFunc(ErrorInvalidArgs{Err: err, Message: "couldn't parse args"})
					return
				}

				// get commands
				root := h.RootFactory(h.session, event)
				root.SetArgs(args)
				var w io.Writer
				if root.OutOrStdout() == os.Stdout { // only override out if it's not already overriden on the command
					w = h.OutWriterFactory(h.session, event.ChannelID)
					root.SetOut(w)
				}
				root.Use = prefix
				err = root.Execute()
				// flush writer before calling errfunc
				if b, ok := w.(*BufferedMessageWriter); ok {
					_, flushErr := b.Flush()
					if flushErr != nil && h.ErrFunc != nil {
						go h.ErrFunc(ErrorInvalidArgs{
							Session: h.session,
							Event:   event,
							Message: "couldn't flush output",
							Err:     flushErr,
						})
					}
				}
				if err != nil && h.ErrFunc != nil {
					go h.ErrFunc(ErrorInvalidArgs{
						Event:   event,
						Session: h.session,
						Err:     err,
						Message: "couldn't execute command",
					})
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

package dgcobra

import (
	"errors"

	"github.com/bwmarrin/discordgo"
)

// ErrorInvalidArgs indicates invalid arguments were passed in a command. You can call Unwrap() to get the underlying error.
type ErrorInvalidArgs struct {
	Session *discordgo.Session
	Event   *discordgo.MessageCreate
	Message string
	Err     error
}

var (
	ErrInvalidCommand = errors.New("error: invalid command")
)

func (err ErrorInvalidArgs) Error() string {
	return err.Message
}

// Unwrap returns the underlying error behind ErrorInvalidArgs.
func (err ErrorInvalidArgs) Unwrap() error {
	return err.Err
}

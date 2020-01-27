package cmd

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/cobra"
)

func EchoCmdFactory(_ *discordgo.Session, _ *discordgo.MessageCreate) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "echo",
		Aliases: []string{"e"},
		Short:   "Echo the given message",
		Long:    "Echoes the given message. Requires at least one argument.",
		Args:    cobra.MinimumNArgs(1),
	}
	caps := cmd.Flags().BoolP("caps", "c", false, "full caps message")
	blacklist := cmd.Flags().StringSliceP("blacklist", "b", []string{}, "blacklist words")
	cmd.Run = func(cmd *cobra.Command, args []string) {
		args = filterArray(args, *blacklist)
		msg := strings.Join(args, " ")
		if *caps {
			msg = strings.ToUpper(msg)
		}
		cmd.Println(msg)
	}
	return cmd
}

func filterArray(args []string, blacklist []string) []string {
	msgArgs := make([]string, 0)
	for _, arg := range args {
		blacklisted := false
		for _, s := range blacklist {
			if arg == s {
				blacklisted = true
				break
			}
		}
		if !blacklisted {
			msgArgs = append(msgArgs, arg)
		}
	}
	return msgArgs
}

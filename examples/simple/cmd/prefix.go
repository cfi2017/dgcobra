package cmd

import (
	"github.com/bwmarrin/discordgo"
	"github.com/spf13/cobra"
)

var (
	Prefixes = make(map[string]string)
)

func prefixCmdFactory(_ *discordgo.Session, e *discordgo.MessageCreate) *cobra.Command {
	return &cobra.Command{
		Use:   "prefix",
		Short: "Set the prefix for this guild",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			Prefixes[e.GuildID] = args[0]
			cmd.Println("Successfully set prefix.")
		},
	}
}

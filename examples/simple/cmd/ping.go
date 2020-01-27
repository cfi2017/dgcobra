package cmd

import (
	"github.com/bwmarrin/discordgo"
	"github.com/spf13/cobra"
)

func pingCmdFactory(s *discordgo.Session, _ *discordgo.MessageCreate) *cobra.Command {
	return &cobra.Command{
		Use:     "ping",
		Aliases: []string{"p"},
		Short:   "Pong.",
		Run: func(cmd *cobra.Command, args []string) {
			// cmd.Printf prints to a dgcobra.MessageWriter (and the underlying channel, usually the channel of the event)
			cmd.Printf("Pong. (%s)\n", s.LastHeartbeatAck.Sub(s.LastHeartbeatSent).String())
		},
	}
}

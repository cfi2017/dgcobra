package cmd

import (
	"github.com/bwmarrin/discordgo"
	"github.com/spf13/cobra"
)

func RootCmdFactory(session *discordgo.Session, event *discordgo.MessageCreate) *cobra.Command {
	cmd := &cobra.Command{
		Short: "Simple example bot",
		Example: `
!ping - pong.
!echo <message> - echoes the given message.
!help - get help
`,
	}
	cmd.AddCommand(
		PingCmdFactory(session, event),
		EchoCmdFactory(session, event),
	)
	return cmd
}

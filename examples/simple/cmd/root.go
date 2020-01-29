package cmd

import (
	"github.com/bwmarrin/discordgo"
	"github.com/spf13/cobra"
)

// RootCmdFactory is an example root command factory for this package. Add subcommands with cmd.AddCommand.
func RootCmdFactory(session *discordgo.Session, event *discordgo.MessageCreate) *cobra.Command {
	cmd := &cobra.Command{
		// short usage
		Short: "Simple example bot",
		// example is shown when help is called
		Example: `
!ping - pong.
!echo <message> - echoes the given message.
!help - get help
`,
	}
	// register child commands
	cmd.AddCommand(
		pingCmdFactory(session, event),
		echoCmdFactory(session, event),
		prefixCmdFactory(session, event),
	)
	return cmd
}

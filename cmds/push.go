package cmds

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var pushCommand = &cobra.Command{
	Use:   "push",
	Short: "push flist",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Msg("pushing...")
	},
}

func init() {
	RootCommand.AddCommand(pushCommand)
}

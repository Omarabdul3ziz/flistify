package cmds

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var runCommand = &cobra.Command{
	Use:   "run",
	Short: "run flist",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Msg("running...")
	},
}

func init() {
	RootCommand.AddCommand(runCommand)
}

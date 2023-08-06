package cmds

import (
	"github.com/omarabdul3ziz/flistify/internal/hub"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var hubCommand = &cobra.Command{
	Use:   "hub",
	Short: "push/pull flist to/from 0Hub",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Msg("hub...")
	},
}

var pushCommand = &cobra.Command{
	Use:   "push",
	Short: "push flist to 0Hub. export HUB_JWT",
	Run: func(cmd *cobra.Command, args []string) {
		source, _ := cmd.Flags().GetString("source")

		log.Info().Msgf("pushing %v...", source)

		if err := hub.Push(source); err != nil {
			log.Error().Err(err).Msgf("failed pushing flist from source: %v", source)
		}
	},
}

var pullCommand = &cobra.Command{
	Use:   "pull",
	Short: "pull flist to 0Hub.",
	Run: func(cmd *cobra.Command, args []string) {
		remote, _ := cmd.Flags().GetString("remote")

		log.Info().Msgf("pulling %v...", remote)
	},
}

func init() {
	pullCommand.PersistentFlags().String("remote", "", "flist name on 0Hub. repo/flist")
	pushCommand.PersistentFlags().String("source", "", "local path either for rootfs or .tar.gz file")

	hubCommand.AddCommand(pullCommand)
	hubCommand.AddCommand(pushCommand)

	RootCommand.AddCommand(hubCommand)
}

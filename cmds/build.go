package cmds

import (
	"github.com/omarabdul3ziz/flistify/internal/builder"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var buildCommand = &cobra.Command{
	Use:   "build",
	Short: "Build flist from Zerofile",
	Run: func(cmd *cobra.Command, args []string) {
		zerofilePath, _ := cmd.Flags().GetString("from")
		name, _ := cmd.Flags().GetString("name")

		log.Info().Msgf("building %v from %v...", name, zerofilePath)

		bl, err := builder.NewBuilder(name)
		if err != nil {
			log.Error().Err(err).Msgf("couldn't get new builder for flist: %v", name)
		}

		if err := bl.Build(zerofilePath); err != nil {
			log.Error().Err(err).Msgf("couldn't build this Zerofile: %v", zerofilePath)
		}
	},
}

func init() {
	buildCommand.PersistentFlags().String("from", "", "Path to zerofile to build image from")
	buildCommand.PersistentFlags().String("name", "", "Name of generated flist")

	RootCommand.AddCommand(buildCommand)
}

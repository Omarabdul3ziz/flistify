package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/omarabdul3ziz/flistify/cmds"
)

func initLogger() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
}

func main() {
	initLogger()

	if err := cmds.RootCommand.Execute(); err != nil {
		log.Error().Err(err).Msgf("failed to execute root command: %v", err)
	}
}

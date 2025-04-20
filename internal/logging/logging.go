package logging

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func ConfigureLogging() {
	// Pretty logs in dev
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	debugMode := os.Getenv("DEBUG_LOGGING") == "1"

	if debugMode {
		zerolog.SetGlobalLevel(zerolog.InfoLevel) // includes info, warn, error, fatal
	} else {
		zerolog.SetGlobalLevel(zerolog.ErrorLevel) // only error and fatal
	}

	// log.Debug().Msg("This is a debug log")  // hidden in both cases
	// log.Info().Msg("This is an info log")   // only in debug
	// log.Warn().Msg("This is a warning log") // only in debug
	// log.Error().Msg("This is an error log") // always shown
	// log.Fatal().Msg("This is a fatal log")  // always shown, exits immediately
}

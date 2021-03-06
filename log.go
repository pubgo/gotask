package gotask

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

var logger zerolog.Logger

func init() {
	logger = log.With().Str("pkg", "gotask").Logger()
}

func _InitDebugLog() {
	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Timestamp().Logger()
	logger = log.Logger
}

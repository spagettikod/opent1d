package envctx

import (
	"os"

	"github.com/rs/zerolog"
)

const (
	LOG_LEVEL = "OPENT1D_LOGLEVEL"
)

func EnvToLogLevel() zerolog.Level {
	logLevel, found := os.LookupEnv(LOG_LEVEL)
	if !found {
		return zerolog.ErrorLevel
	}
	switch logLevel {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	default:
		return zerolog.ErrorLevel
	}
}

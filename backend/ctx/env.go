package ctx

import "os"

const (
	DEBUG_ENV = "OPENT1D_DEBUG"
)

func IsDebug() bool {
	_, found := os.LookupEnv(DEBUG_ENV)
	return found
}

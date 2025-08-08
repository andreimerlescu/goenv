package env

import (
	"fmt"
	"log"
	"os"
)

func init() {
	// Initialize Data
	ZeroList = make([]string, 0)
	ZeroMap = make(map[string]string)
	ShowVerbose = UseLogger && EnableVerboseLogging

	// Create loggers
	OutLogger = log.New(os.Stdout, "INFO:"+packageName, LogFlags)
	ErrLogger = log.New(os.Stderr, " ERR:"+packageName, LogFlags)

	// Intentionally Use Import Side-Effects
	if UseMagic {
		Magic()
	}
}

// logError is an internal helper to centralize error logging logic.
func logError(envName string, err error) {
	if !PrintErrors {
		return
	}
	msg := fmt.Sprintf("Error processing env var '%s': %v", envName, err)
	if UseLogger {
		ErrLogger.Println(msg)
	} else {
		_, _ = fmt.Fprintln(os.Stderr, msg)
	}
}

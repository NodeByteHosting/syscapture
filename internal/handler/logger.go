package handler

import (
	"os"

	"github.com/sirupsen/logrus"
)

// NewLogger creates and configures a new logrus Logger instance.
func NewLogger(debug bool) *logrus.Logger {
	logger := logrus.New()

	// Set the output to stdout
	logger.SetOutput(os.Stdout)

	// Set the log level based on the debug setting
	if debug {
		logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}

	// Set the log format (e.g., JSONFormatter, TextFormatter)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	return logger
}

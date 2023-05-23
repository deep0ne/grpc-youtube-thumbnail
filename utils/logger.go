package utils

import (
	"os"

	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

func NewLogger() *logrus.Logger {
	logger := logrus.New()
	logger.Out = os.Stdout
	logger.Level = logrus.DebugLevel
	formatter := &prefixed.TextFormatter{
		ForceColors:     true,
		ForceFormatting: true,
	}

	formatter.SetColorScheme(&prefixed.ColorScheme{
		InfoLevelStyle:  "green+b",
		WarnLevelStyle:  "red+b",
		DebugLevelStyle: "blue+h",
	})

	logger.Formatter = formatter
	return logger
}

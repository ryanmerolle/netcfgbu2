package utils

import (
	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

func InitLogger() {
	// Set the log format, level, etc.
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	Log.SetLevel(logrus.DebugLevel)
}

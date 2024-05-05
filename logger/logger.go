package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

var logger = logrus.New()

// get logger returns an instance of logger
func GetLogger() *logrus.Logger {
	return logger
}

// set logger set an instance of logger, used when you do not want to rewrite logfile (i.e. when testing)
func SetLogger(l *logrus.Logger) {
	logger = l
}

func init() {

	logger.SetLevel(logrus.DebugLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logger.SetOutput(os.Stdout)
	file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		logger.Fatal("Failed to open log file:")
	}
	logger.SetOutput(file)

}

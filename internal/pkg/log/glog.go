package log

import (
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type (
	glog struct {
		logger *logrus.Entry
	}
)

func newGlog() *glog {
	return &glog{
		logger: newLogrusEntry(),
	}
}

func newGlogWithField(key, value string) *glog {
	return &glog{
		logger: newLogrusEntry().WithField(key, value),
	}
}
func newGlogWithFields(fields Fields) *glog {
	return &glog{
		logger: newLogrusEntry().WithFields(fields),
	}
}

func newLogrusEntry() *logrus.Entry {
	logger := logrus.New()
	logger.SetFormatter(getFormaterFromEnv())
	logger.SetLevel(getLevelFromEnv())
	logger.SetOutput(getOutputFromEnv())
	return logrus.NewEntry(logger)
}

func getFormaterFromEnv() logrus.Formatter {
	if strings.ToLower(os.Getenv("LOG_FORMAT")) == "json" {
		return &logrus.JSONFormatter{
			TimestampFormat: time.RFC1123,
		}
	}
	return &logrus.TextFormatter{
		TimestampFormat: time.RFC1123,
		FullTimestamp:   true,
	}
}

func getLevelFromEnv() logrus.Level {
	lvl, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		lvl = logrus.DebugLevel
	}
	return lvl
}

func getOutputFromEnv() io.WriteCloser {
	out := os.Getenv("LOG_OUTPUT")
	if strings.HasPrefix(out, filePrefix) {
		name := out[len(filePrefix):]
		file, err := os.Create(name)
		if err != nil {
			log.Printf("failed to create log file: %s, err; %v\n", name, err)
		}
		return file
	}
	return os.Stdout
}

// Info print info
func (g *glog) Info(args ...interface{}) {
	g.logger.Infoln(args...)
}

// Debug print debug
func (g *glog) Debug(args ...interface{}) {
	g.logger.Debugln(args...)
}

// Warn print warn
func (g *glog) Warn(args ...interface{}) {
	g.logger.Warnln(args...)
}

// Error print error
func (g *glog) Error(args ...interface{}) {
	g.logger.Errorln(args...)
}

// Panic panic
func (g *glog) Panic(args ...interface{}) {
	g.logger.Panicln(args...)
}

// Infof print info with format
func (g *glog) Infof(format string, args ...interface{}) {
	g.logger.Infof(format, args...)
}

// Debugf print debug with format
func (g *glog) Debugf(format string, args ...interface{}) {
	g.logger.Debugf(format, args...)
}

// Warnf print warn with format
func (g *glog) Warnf(format string, args ...interface{}) {
	g.logger.Warnf(format, args...)
}

// Errorf print error with format
func (g *glog) Errorf(format string, args ...interface{}) {
	g.logger.Errorf(format, args...)
}

// Panicf panic with format
func (g *glog) Panicf(format string, args ...interface{}) {
	g.logger.Panicf(format, args...)
}

// WithField return a new logger with fields
// func (g *glog) WithField(key string, val interface{}) Logger {
// 	return &glog{
// 		logger: g.logger.WithField(key, val),
// 	}
// }

// WithFields return a new logger with fields
// func (g *glog) WithFields(fields Fields) Logger {
// 	return &glog{
// 		logger: g.logger.WithFields(logrus.Fields(fields)),
// 	}
// }

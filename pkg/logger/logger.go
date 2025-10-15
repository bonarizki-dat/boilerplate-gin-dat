package logger

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"strings"
	"time"
)

var logger = logrus.New()

func init() {
	logger.Level = logrus.InfoLevel
	logger.Formatter = &formatter{}
	logger.SetReportCaller(true)

	// Open log file
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		// Fallback to stdout if opening file fails
		logger.SetOutput(os.Stdout)
	} else {
		// Write logs to both stdout and file
		mw := io.MultiWriter(os.Stdout, file)
		logger.SetOutput(mw)
	}
}

// SetLogLevel sets the log level for the logger
func SetLogLevel(level logrus.Level) {
	logger.Level = level
}

// Fields type alias for logrus.Fields
type Fields logrus.Fields

// Debugf logs a message at level Debug.
func Debugf(format string, args ...interface{}) {
	if logger.Level >= logrus.DebugLevel {
		// Create a log entry with empty fields
		entry := logger.WithFields(logrus.Fields{})
		entry.Debugf(format, args...)
	}
}

// Infof logs a message at level Info.
func Infof(format string, args ...interface{}) {
	if logger.Level >= logrus.InfoLevel {
		// Create a log entry with empty fields
		entry := logger.WithFields(logrus.Fields{})
		entry.Infof(format, args...)
	}
}

// Warnf logs a message at level Warn.
func Warnf(format string, args ...interface{}) {
	if logger.Level >= logrus.WarnLevel {
		// Create a log entry with empty fields
		entry := logger.WithFields(logrus.Fields{})
		entry.Warnf(format, args...)
	}
}

// Errorf logs a message at level Error.
func Errorf(format string, args ...interface{}) {
	if logger.Level >= logrus.ErrorLevel {
		// Create a log entry with empty fields
		entry := logger.WithFields(logrus.Fields{})
		entry.Errorf(format, args...)
	}
}

// Fatalf logs a message at level Fatal and exits the program.
func Fatalf(format string, args ...interface{}) {
	if logger.Level >= logrus.FatalLevel {
		// Create a log entry with empty fields
		entry := logger.WithFields(logrus.Fields{})
		entry.Fatalf(format, args...)
	}
}

// Formatter implements logrus.Formatter
type formatter struct {
	prefix string
}

// Format builds the log message according to custom format
func (f *formatter) Format(entry *logrus.Entry) ([]byte, error) {
	var sb bytes.Buffer

	// Write log level in uppercase
	sb.WriteString(strings.ToUpper(entry.Level.String()))
	sb.WriteString(" ")
	// Write timestamp in RFC3339 format
	sb.WriteString(entry.Time.Format(time.RFC3339))
	sb.WriteString(" ")
	// Write prefix (if any)
	sb.WriteString(f.prefix)
	// Write the actual log message
	sb.WriteString(entry.Message)
	// Ensure each log entry ends with a newline
	sb.WriteByte('\n')

	return sb.Bytes(), nil
}
package logger

import (
	"bytes"
	"context"
	"io"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
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

// contextKey type for context values to avoid collisions.
type contextKey string

// RequestIDContextKey is the key used to store request_id in context.Context.
// Middleware should set it via context.WithValue(c.Request.Context(), logger.RequestIDContextKey, requestID).
const RequestIDContextKey contextKey = "request_id"

// Formatter implements logrus.Formatter
type formatter struct {
	prefix string
}

// Format builds the log message according to custom format.
// If entry has "request_id" in Data, prepends [request_id] to the message.
func (f *formatter) Format(entry *logrus.Entry) ([]byte, error) {
	var sb bytes.Buffer

	sb.WriteString(strings.ToUpper(entry.Level.String()))
	sb.WriteString(" ")
	sb.WriteString(entry.Time.Format(time.RFC3339))
	sb.WriteString(" ")
	sb.WriteString(f.prefix)
	if rid, ok := entry.Data["request_id"]; ok {
		if ridStr, ok := rid.(string); ok && ridStr != "" {
			sb.WriteString("[")
			sb.WriteString(ridStr)
			sb.WriteString("] ")
		}
	}
	sb.WriteString(entry.Message)
	sb.WriteByte('\n')

	return sb.Bytes(), nil
}

// WithRequestID returns a log entry that includes request_id in every log line.
func WithRequestID(requestID string) *logrus.Entry {
	return logger.WithField("request_id", requestID)
}

// FromContext returns a log entry with request_id from ctx if present; otherwise without request_id.
func FromContext(ctx context.Context) *logrus.Entry {
	if ctx == nil {
		return logger.WithFields(logrus.Fields{})
	}
	if rid, ok := ctx.Value(RequestIDContextKey).(string); ok && rid != "" {
		return logger.WithField("request_id", rid)
	}
	return logger.WithFields(logrus.Fields{})
}

// Span holds start time and component/method for START/FINISH logging.
type Span struct {
	requestID string
	component string
	method    string
	start     time.Time
	entry     *logrus.Entry
}

// StartWithRequestID begins a span and logs START. Use for controllers (have request_id from gin).
func StartWithRequestID(requestID, component, method string) *Span {
	s := &Span{
		requestID: requestID,
		component: component,
		method:    method,
		start:     time.Now(),
		entry:     WithRequestID(requestID),
	}
	s.entry.Infof("START %s.%s", component, method)
	return s
}

// Start begins a span using request_id from ctx and logs START. Use for services (receive ctx).
func Start(ctx context.Context, component, method string) *Span {
	entry := FromContext(ctx)
	requestID := ""
	if ctx != nil {
		if rid, ok := ctx.Value(RequestIDContextKey).(string); ok {
			requestID = rid
		}
	}
	s := &Span{
		requestID: requestID,
		component: component,
		method:    method,
		start:     time.Now(),
		entry:     entry,
	}
	s.entry.Infof("START %s.%s", component, method)
	return s
}

// Finish logs FINISH with SUCCESS or FAIL and duration. Call with err from handler return.
func (s *Span) Finish(err error) {
	dur := time.Since(s.start).Milliseconds()
	status := "SUCCESS"
	if err != nil {
		status = "FAIL"
	}
	s.entry.Infof("FINISH %s.%s (%s) duration=%dms", s.component, s.method, status, dur)
}

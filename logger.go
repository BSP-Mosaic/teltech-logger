package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

type severity int

const (
	DEBUG severity = iota
	INFO
	WARN
	ERROR
	CRITICAL

	defaultCallerSkip = 2
)

func (s severity) String() string {
	return logLevelName[s]
}

var logLevelName = [...]string{
	"DEBUG",
	"INFO",
	"WARN",
	"ERROR",
	"CRITICAL",
}

var logLevelValue = map[string]severity{
	"DEBUG":    DEBUG,
	"INFO":     INFO,
	"WARN":     WARN,
	"ERROR":    ERROR,
	"CRITICAL": CRITICAL,
}

// Fields is used to wrap the log entries payload
type Fields map[string]interface{}

// ServiceContext is required by the Stackdriver Error format
type ServiceContext struct {
	Service string `json:"service,omitempty"`
	Version string `json:"version,omitempty"`
}

// ReportLocation is required by the Stackdriver Error format
type ReportLocation struct {
	FilePath     string `json:"filePath"`
	FunctionName string `json:"functionName"`
	LineNumber   int    `json:"lineNumber"`
}

// Context is required by the Stackdriver Error format
type Context struct {
	Data           Fields          `json:"data,omitempty"`
	ReportLocation *ReportLocation `json:"reportLocation,omitempty"`
}

// Payload groups all the data for a log entry
type Payload struct {
	Severity       string          `json:"severity"`
	EventTime      string          `json:"eventTime"`
	Caller         string          `json:"caller,omitempty"`
	Message        string          `json:"message"`
	ServiceContext *ServiceContext `json:"serviceContext,omitempty"`
	Context        *Context        `json:"context,omitempty"`
	Stacktrace     string          `json:"stacktrace,omitempty"`
}

// Log is the main type for the logger package
type Log struct {
	level          severity
	mux            sync.RWMutex
	fields         Fields
	serviceContext *ServiceContext
	writer         io.Writer
	callerSkip     int
}

var (
	defaultLogLevel severity
	service         string
	version         string
)

func init() {
	logLevel, ok := logLevelValue[strings.ToUpper(os.Getenv("LOG_LEVEL"))]
	if !ok {
		fmt.Println("logger WARN: LOG_LEVEL is not valid or not set, defaulting to INFO")
		logLevel = logLevelValue[INFO.String()]
	}

	if os.Getenv("SERVICE") == "" || os.Getenv("VERSION") == "" {
		fmt.Println("logger ERROR: cannot instantiate the logger, make sure the SERVICE and VERSION environment vars are set correctly")
	}

	initConfig(logLevel, os.Getenv("SERVICE"), os.Getenv("VERSION"))
}

func initConfig(lvl severity, svc, ver string) {
	defaultLogLevel = lvl
	service = svc
	version = ver
}

// New instantiates and returns a Log object
func New() *Log {
	l := &Log{
		fields:     Fields{},
		writer:     os.Stdout,
		level:      defaultLogLevel,
		callerSkip: defaultCallerSkip,
	}

	if service != "" && version != "" {
		l.serviceContext = &ServiceContext{
			Service: service,
			Version: version,
		}
	}

	return l
}

// WithOutput creates a copy of a Log with a different output.
func (l *Log) WithOutput(w io.Writer) *Log {
	n := l.With(Fields{})
	n.writer = w
	return n
}

// WithLevel creates a copy of a Log with a different log level
func (l *Log) WithLevel(logLevel severity) *Log {
	n := l.With(Fields{})
	n.level = logLevel
	return n
}

// AddCallerSkip increases the number of callers skipped by caller annotation.
// When building wrappers around the Logger, supplying this value prevents logger
// from always reporting the wrapper code as the caller.
func (l *Log) AddCallerSkip(skip int) {
	l.mux.Lock()
	defer l.mux.Unlock()

	l.callerSkip += skip
}

func (l *Log) log(severity, message, stacktrace string, reportLocation *ReportLocation) {
	l.mux.Lock()
	defer l.mux.Unlock()

	// Do not persist the payload here, just format it, marshal it and return it
	payload := &Payload{
		Severity:       severity,
		EventTime:      time.Now().Format(time.RFC3339),
		Message:        message,
		ServiceContext: l.serviceContext,
		Context: &Context{
			Data:           l.fields,
			ReportLocation: reportLocation,
		},
		Stacktrace: stacktrace,
	}

	b, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("logger ERROR: cannot marshal payload: %s", err)
		return
	}

	l.writer.Write(b)
	l.writer.Write([]byte{'\n'})
}

// Checks whether the specified log level is valid
func (l *Log) isValidLogLevel(s severity) bool {
	l.mux.Lock()
	defer l.mux.Unlock()

	return s >= l.level
}

// fields returns a valid Fields whether or not one exists in the *Log.
func (l *Log) getFields() Fields {
	f := Fields{}

	for k, v := range l.fields {
		f[k] = v
	}

	return f
}

// With is used as a chained method to specify which values go in the log entry's context
func (l *Log) With(fields Fields) *Log {
	l.mux.RLock()
	defer l.mux.RUnlock()

	f := l.getFields()

	for k, v := range fields {
		f[k] = v
	}

	return &Log{
		serviceContext: l.serviceContext,
		fields:         f,
		writer:         l.writer,
		level:          l.level,
		callerSkip:     l.callerSkip,
	}
}

// Debug prints out a message with DEBUG severity level
func (l *Log) Debug(message string) {
	if !l.isValidLogLevel(DEBUG) {
		return
	}

	l.log(DEBUG.String(), message, "", nil)
}

// Debugf prints out a message with DEBUG severity level
func (l *Log) Debugf(message string, args ...interface{}) {
	l.Debug(fmt.Sprintf(message, args...))
}

// Info prints out a message with INFO severity level
func (l *Log) Info(message string) {
	if !l.isValidLogLevel(INFO) {
		return
	}

	l.log(INFO.String(), message, "", nil)
}

// Infof prints out a message with INFO severity level
func (l *Log) Infof(message string, args ...interface{}) {
	l.Info(fmt.Sprintf(message, args...))
}

// Warn prints out a message with WARN severity level
func (l *Log) Warn(message string) {
	if !l.isValidLogLevel(WARN) {
		return
	}

	l.log(WARN.String(), message, "", nil)
}

// Warnf prints out a message with WARN severity level
func (l *Log) Warnf(message string, args ...interface{}) {
	l.Warn(fmt.Sprintf(message, args...))
}

// Error prints out a message with ERROR severity level
func (l *Log) Error(message string) {
	l.error(ERROR.String(), message)
}

// Errorf prints out a message with ERROR severity level
func (l *Log) Errorf(message string, args ...interface{}) {
	l.error(ERROR.String(), fmt.Sprintf(message, args...))
}

// Fatal is equivalent to Error() followed by a call to os.Exit(1).
// It prints out a message with CRITICAL severity level
func (l *Log) Fatal(message string) {
	l.error(CRITICAL.String(), message)
	os.Exit(1)
}

// Fatalf is equivalent to Errorf() followed by a call to os.Exit(1).
// It prints out a message with CRITICAL severity level
func (l *Log) Fatalf(message string, args ...interface{}) {
	l.error(CRITICAL.String(), fmt.Sprintf(message, args...))
	os.Exit(1)
}

// ERROR prints out a message with the passed severity level (ERROR or CRITICAL)
func (l *Log) error(severity, message string) {
	buffer := make([]byte, 1024)
	buffer = buffer[:runtime.Stack(buffer, false)]
	fpc, file, line, _ := runtime.Caller(l.callerSkip)

	funcName := "unknown"
	fun := runtime.FuncForPC(fpc)
	if fun != nil {
		_, funcName = filepath.Split(fun.Name())
	}

	l.log(severity, message, string(buffer), &ReportLocation{
		FilePath:     file,
		FunctionName: funcName,
		LineNumber:   line,
	})
}

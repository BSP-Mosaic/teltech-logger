package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type severity int

const (
	DEBUG severity = iota
	INFO
	WARN
	ERROR
	CRITICAL
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
	payload *Payload
	writer  io.Writer
}

var (
	logLevel severity
	service  string
	version  string
)

func init() {
	ll, ok := logLevelValue[strings.ToUpper(os.Getenv("LOG_LEVEL"))]
	if !ok {
		fmt.Println("logger WARN: LOG_LEVEL is not valid or not set, defaulting to INFO")
		logLevel = logLevelValue[INFO.String()]
	} else {
		logLevel = ll
	}

	if os.Getenv("SERVICE") == "" || os.Getenv("VERSION") == "" {
		fmt.Println("logger ERROR: cannot instantiate the logger, make sure the SERVICE and VERSION environment vars are set correctly")
	}

	initConfig(logLevel, os.Getenv("SERVICE"), os.Getenv("VERSION"))
}

func initConfig(lvl severity, svc, ver string) {
	logLevel = lvl
	service = svc
	version = ver
}

// New instantiates and returns a Log object
func New() *Log {
	// Set the ServiceContext only within a GCP context
	p := &Payload{}
	if service != "" && version != "" {
		p = &Payload{
			ServiceContext: &ServiceContext{
				Service: service,
				Version: version,
			},
		}
	}

	return &Log{
		payload: p,
		writer:  os.Stdout,
	}
}

// WithOutput creates a copy of a Log with a different output.
func (l *Log) WithOutput(w io.Writer) *Log {
	n := l.With(Fields{})
	n.writer = w
	return n
}

func (l *Log) log(severity, message string) {
	// Do not persist the payload here, just format it, marshal it and return it
	l.payload = &Payload{
		Severity:       severity,
		EventTime:      time.Now().Format(time.RFC3339),
		Message:        message,
		ServiceContext: l.payload.ServiceContext,
		Context:        l.payload.Context,
		Stacktrace:     l.payload.Stacktrace,
	}

	payload, ok := json.Marshal(l.payload)
	if ok != nil {
		fmt.Printf("logger ERROR: cannot marshal payload: %s", ok.Error())
	}

	fmt.Fprintln(l.writer, string(payload))
}

// Checks whether the specified log level is valid in the current environment
func isValidLogLevel(s severity) bool {
	return s >= logLevel
}

// fields returns a valid Fields whether or not one exists in the *Log.
func (l *Log) fields() Fields {
	f := make(Fields)
	if l.payload == nil {
		return f
	}
	if l.payload.Context == nil {
		return f
	}

	if l.payload.Context.Data == nil {
		return f
	}

	for k, v := range l.payload.Context.Data {
		f[k] = v
	}
	return f
}

// With is used as a chained method to specify which values go in the log entry's context
func (l *Log) With(fields Fields) *Log {
	f := l.fields()
	for k, v := range fields {
		f[k] = v
	}

	return &Log{
		payload: &Payload{
			ServiceContext: l.payload.ServiceContext,
			Context: &Context{
				Data: f,
			},
			Stacktrace: "",
		},
		writer: os.Stdout,
	}
}

// Debug prints out a message with DEBUG severity level
func (l Log) Debug(message string) {
	if !isValidLogLevel(DEBUG) {
		return
	}

	l.log(DEBUG.String(), message)
}

// Debugf prints out a message with DEBUG severity level
func (l Log) Debugf(message string, args ...interface{}) {
	l.Debug(fmt.Sprintf(message, args...))
}

// Info prints out a message with INFO severity level
func (l Log) Info(message string) {
	if !isValidLogLevel(INFO) {
		return
	}

	l.log(INFO.String(), message)
}

// Infof prints out a message with INFO severity level
func (l Log) Infof(message string, args ...interface{}) {
	l.Info(fmt.Sprintf(message, args...))
}

// Warn prints out a message with WARN severity level
func (l Log) Warn(message string) {
	if !isValidLogLevel(WARN) {
		return
	}

	l.log(WARN.String(), message)
}

// Warnf prints out a message with WARN severity level
func (l Log) Warnf(message string, args ...interface{}) {
	l.Warn(fmt.Sprintf(message, args...))
}

// Error prints out a message with ERROR severity level
func (l Log) Error(message string) {
	l.error(ERROR.String(), message)
}

// Errorf prints out a message with ERROR severity level
func (l Log) Errorf(message string, args ...interface{}) {
	l.error(ERROR.String(), fmt.Sprintf(message, args...))
}

// Fatal is equivalent to Error() followed by a call to os.Exit(1).
// It prints out a message with CRITICAL severity level
func (l Log) Fatal(message string) {
	l.error(CRITICAL.String(), message)
	os.Exit(1)
}

// Fatalf is equivalent to Errorf() followed by a call to os.Exit(1).
// It prints out a message with CRITICAL severity level
func (l Log) Fatalf(message string, args ...interface{}) {
	l.error(CRITICAL.String(), fmt.Sprintf(message, args...))
	os.Exit(1)
}

// ERROR prints out a message with the passed severity level (ERROR or CRITICAL)
func (l Log) error(severity, message string) {
	buffer := make([]byte, 1024)
	buffer = buffer[:runtime.Stack(buffer, false)]
	fpc, file, line, _ := runtime.Caller(2)

	funcName := "unknown"
	fun := runtime.FuncForPC(fpc)
	if fun != nil {
		_, funcName = filepath.Split(fun.Name())
	}

	// Set the data when the context is empty
	if l.payload.Context == nil {
		l.payload.Context = &Context{
			Data: Fields{},
		}
	}

	l.payload = &Payload{
		ServiceContext: l.payload.ServiceContext,
		Context: &Context{
			Data: l.payload.Context.Data,
			ReportLocation: &ReportLocation{
				FilePath:     file,
				FunctionName: funcName,
				LineNumber:   line,
			},
		},
		Stacktrace: string(buffer),
	}

	l.log(severity, message)
}

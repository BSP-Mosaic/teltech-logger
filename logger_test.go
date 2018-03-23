package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestLoggerInfoWithOneTimeContext(t *testing.T) {
	initConfig(DEBUG, "my-app", "1.0")

	buf := new(bytes.Buffer)

	log := New().With(Fields{
		"key":      "value",
		"function": "TestLoggerDebug",
	}).SetWriter(buf)

	log.Info("INFO message")
	expected := fmt.Sprintf(`{"severity":"INFO","eventTime":"%s","message":"INFO message","serviceContext":{"service":"my-app","version":"1.0"},"context":{"data":{"function":"TestLoggerDebug","key":"value"}}}`, time.Now().Format(time.RFC3339))
	got := strings.TrimRight(buf.String(), "\n")
	if expected != got {
		t.Errorf("output %s does not match expected string %s", got, expected)
	}

	// Clean-up the buffer in preparation for new assertions
	buf.Reset()

	log.With(Fields{"foo": "bar"}).SetWriter(buf).Info("unique INFO message")
	expected = fmt.Sprintf(`{"severity":"INFO","eventTime":"%s","message":"unique INFO message","serviceContext":{"service":"my-app","version":"1.0"},"context":{"data":{"foo":"bar"}}}`, time.Now().Format(time.RFC3339))
	got = strings.TrimRight(buf.String(), "\n")
	if expected != got {
		t.Errorf("output file %s does not match expected string %s", got, expected)
	}

	// Clean-up the buffer in preparation for new assertions
	buf.Reset()

	log.SetWriter(buf).Info("unique INFO message")
	expected = fmt.Sprintf(`{"severity":"INFO","eventTime":"%s","message":"unique INFO message","serviceContext":{"service":"my-app","version":"1.0"},"context":{"data":{"function":"TestLoggerDebug","key":"value"}}}`, time.Now().Format(time.RFC3339))
	got = strings.TrimRight(buf.String(), "\n")
	if expected != got {
		t.Errorf("output %s does not match expected string %s", got, expected)
	}
}

func TestLoggerErrorWithOneTimeContext(t *testing.T) {
	initConfig(DEBUG, "my-app", "1.0")

	buf := new(bytes.Buffer)

	log := New().With(Fields{
		"key":      "value",
		"function": "TestLoggerError",
	}).SetWriter(buf)

	log.Error("ERROR message")
	expected := fmt.Sprintf(`{"severity":"ERROR","eventTime":"%s","message":"ERROR message","serviceContext":{"service":"my-app","version":"1.0"},"context":{"data":{"function":"TestLoggerError","key":"value"},"reportLocation"`, time.Now().Format(time.RFC3339))
	got := strings.TrimRight(buf.String(), "\n")
	if !strings.Contains(got, expected) {
		t.Errorf("output %s does not contain substring %s", got, expected)
	}

	// Check that the ERROR entry contains the context
	if !strings.Contains(got, `"context":{"data":{"function":"TestLoggerError","key":"value"}`) {
		t.Errorf("output %s does not contain the context", got)
	}

	// Check that the ERROR entry has an stacktrace key
	if !strings.Contains(got, "stacktrace") {
		t.Errorf("output %s does not contain a stacktrace key", got)
	}

	// Clean-up the buffer in preparation for new assertions
	buf.Reset()

	log.With(Fields{"foo": "bar"}).SetWriter(buf).Error("unique ERROR message")
	expected = fmt.Sprintf(`{"severity":"ERROR","eventTime":"%s","message":"unique ERROR message","serviceContext":{"service":"my-app","version":"1.0"},"context":{"data":{"foo":"bar"},"reportLocation"`, time.Now().Format(time.RFC3339))
	got = strings.TrimRight(buf.String(), "\n")
	if !strings.Contains(got, expected) {
		t.Errorf("output %s does not contain substring %s", got, expected)
	}

	// Check that the ERROR entry contains the context
	if !strings.Contains(got, `"context":{"data":{"foo":"bar"}`) {
		t.Errorf("output %s does not contain the context", got)
	}

	// Check that the ERROR entry has an stacktrace key
	if !strings.Contains(got, "stacktrace") {
		t.Errorf("output %s does not contain a stacktrace key", got)
	}

	// Clean-up the buffer in preparation for new assertions
	buf.Reset()

	log.SetWriter(buf).Error("unique ERROR message")
	expected = fmt.Sprintf(`{"severity":"ERROR","eventTime":"%s","message":"unique ERROR message","serviceContext":{"service":"my-app","version":"1.0"},"context":{"data":{"function":"TestLoggerError","key":"value"},"reportLocation"`, time.Now().Format(time.RFC3339))
	got = strings.TrimRight(buf.String(), "\n")
	if !strings.Contains(got, expected) {
		t.Errorf("output %s does not contain substring %s", got, expected)
	}

	// Check that the ERROR entry contains the context
	if !strings.Contains(got, `"context":{"data":{"function":"TestLoggerError","key":"value"}`) {
		t.Errorf("output %s does not contain the context", got)
	}

	// Check that the ERROR entry has an stacktrace key
	if !strings.Contains(got, "stacktrace") {
		t.Errorf("output %s does not contain a stacktrace key", got)
	}
}

func TestLoggerWithDifferentLogLevels(t *testing.T) {
	initConfig(WARN, "my-app", "1.0")

	buf := new(bytes.Buffer)

	log := New().With(Fields{
		"key": "value",
	}).SetWriter(buf)

	// LogLevel set to WARN, DEBUG messages should not be output
	log.Debug("DEBUG message")
	got := strings.TrimRight(buf.String(), "\n")

	if got != "" {
		t.Errorf("output %s does not match empty string", got)
	}

	// LogLevel set to WARN, INFO messages should not be output
	log.Info("INFO message")
	got = strings.TrimRight(buf.String(), "\n")

	if got != "" {
		t.Errorf("output %s does not match empty string", got)
	}

	log.Warn("WARN message")
	expected := fmt.Sprintf(`{"severity":"WARN","eventTime":"%s","message":"WARN message","serviceContext":{"service":"my-app","version":"1.0"},"context":{"data":{"key":"value"}}}`, time.Now().Format(time.RFC3339))
	got = strings.TrimRight(buf.String(), "\n")
	if expected != got {
		t.Errorf("output %s does not match expected string %s", got, expected)
	}

	// Clean-up the buffer in preparation for new assertions
	buf.Reset()

	log.Error("ERROR message")
	expected = fmt.Sprintf(`{"severity":"ERROR","eventTime":"%s","message":"ERROR message","serviceContext":{"service":"my-app","version":"1.0"},"context":{"data":{"function":"TestLoggerError","key":"value"},"reportLocation"`, time.Now().Format(time.RFC3339))
	got = strings.TrimRight(buf.String(), "\n")
	if strings.Contains(got, expected) {
		t.Errorf("expecting %s; got %s", expected, got)
	}
}

func TestLoggerDebugWithImplicitContext(t *testing.T) {
	initConfig(DEBUG, "my-app", "1.0")

	buf := new(bytes.Buffer)

	log := New().With(Fields{
		"key":      "value",
		"function": "TestLoggerDebug",
	}).SetWriter(buf)

	log.Debug("DEBUG message")

	expected := fmt.Sprintf(`{"severity":"DEBUG","eventTime":"%s","message":"DEBUG message","serviceContext":{"service":"my-app","version":"1.0"},"context":{"data":{"function":"TestLoggerDebug","key":"value"}}}`, time.Now().Format(time.RFC3339))
	got := strings.TrimRight(buf.String(), "\n")
	if expected != got {
		t.Errorf("output %s does not match expected string %s", got, expected)
	}
}

func TestLoggerDebugWithoutContext(t *testing.T) {
	initConfig(DEBUG, "my-app", "1.0")

	buf := new(bytes.Buffer)
	log := New().SetWriter(buf)

	log.Debug("DEBUG message")
	expected := fmt.Sprintf(`{"severity":"DEBUG","eventTime":"%s","message":"DEBUG message","serviceContext":{"service":"my-app","version":"1.0"}}`, time.Now().Format(time.RFC3339))
	got := strings.TrimRight(buf.String(), "\n")
	if expected != got {
		t.Errorf("output %s does not match expected string %s", got, expected)
	}
}

func TestLoggerDebugfWithoutContext(t *testing.T) {
	initConfig(DEBUG, "my-app", "1.0")

	buf := new(bytes.Buffer)

	log := New().SetWriter(buf)

	param := "with param"
	log.Debugf("DEBUG message %s", param)
	expected := fmt.Sprintf(`{"severity":"DEBUG","eventTime":"%s","message":"DEBUG message with param","serviceContext":{"service":"my-app","version":"1.0"}}`, time.Now().Format(time.RFC3339))
	got := strings.TrimRight(buf.String(), "\n")
	if expected != got {
		t.Errorf("output %s does not match expected string %s", got, expected)
	}
}

func TestLoggerInfo(t *testing.T) {
	initConfig(DEBUG, "my-app", "1.0")

	buf := new(bytes.Buffer)

	log := New().With(Fields{
		"key":      "value",
		"function": "TestLoggerInfo",
	}).SetWriter(buf)

	log.Info("INFO message")
	expected := fmt.Sprintf(`{"severity":"INFO","eventTime":"%s","message":"INFO message","serviceContext":{"service":"my-app","version":"1.0"},"context":{"data":{"function":"TestLoggerInfo","key":"value"}}}`, time.Now().Format(time.RFC3339))
	got := strings.TrimRight(buf.String(), "\n")
	if expected != got {
		t.Errorf("output %s does not match expected string %s", got, expected)
	}
}

func TestLoggerInfof(t *testing.T) {
	initConfig(DEBUG, "my-app", "1.0")

	buf := new(bytes.Buffer)

	log := New().With(Fields{
		"key":      "value",
		"function": "TestLoggerInfo",
	}).SetWriter(buf)

	param := "with param"
	log.Infof("INFO message %s", param)
	expected := fmt.Sprintf(`{"severity":"INFO","eventTime":"%s","message":"INFO message with param","serviceContext":{"service":"my-app","version":"1.0"},"context":{"data":{"function":"TestLoggerInfo","key":"value"}}}`, time.Now().Format(time.RFC3339))
	got := strings.TrimRight(buf.String(), "\n")
	if expected != got {
		t.Errorf("output %s does not match expected string %s", got, expected)
	}
}

func TestResponseIsValidJson(t *testing.T) {
	initConfig(DEBUG, "my-app", "1.0")

	buf := new(bytes.Buffer)
	log := New().With(Fields{"key": "value"}).SetWriter(buf)

	log.Error("ERROR message")
	got := strings.TrimRight(buf.String(), "\n")

	p := Payload{}
	err := json.Unmarshal([]byte(got), &p)
	if err != nil {
		t.Errorf("response cannot be unmarshalled: %s", err.Error())
	}
}

func TestGetCallerFunctionName(t *testing.T) {
	initConfig(DEBUG, "my-app", "1.0")

	buf := new(bytes.Buffer)
	log := New().With(Fields{"key": "value"}).SetWriter(buf)

	log.Error("ERROR message")
	got := strings.TrimRight(buf.String(), "\n")

	// Encode the returned error and check the "functionName" key value
	p := Payload{}
	err := json.Unmarshal([]byte(got), &p)
	if err != nil {
		t.Errorf("failed to unmarshal payload: %s", err.Error())
	}

	expected := "logger.TestGetCallerFunctionName"
	if p.Context.ReportLocation.FunctionName != expected {
		t.Errorf("output %s does not containsubstring %s", p.Context.ReportLocation.FunctionName, expected)
	}
}

func TestLoggerError(t *testing.T) {
	initConfig(DEBUG, "my-app", "1.0")

	buf := new(bytes.Buffer)

	log := New().With(Fields{
		"key":      "value",
		"function": "TestLoggerError",
	}).SetWriter(buf)

	log.Error("ERROR message")
	expected := fmt.Sprintf(`{"severity":"ERROR","eventTime":"%s","message":"ERROR message","serviceContext":{"service":"my-app","version":"1.0"},"context":{"data":{"function":"TestLoggerError","key":"value"},"reportLocation"`, time.Now().Format(time.RFC3339))
	got := strings.TrimRight(buf.String(), "\n")
	if !strings.Contains(got, expected) {
		t.Errorf("output %s does not containsubstring %s", got, expected)
	}

	// Check that the ERROR entry contains the context
	if !strings.Contains(got, `"context":{"data":{"function":"TestLoggerError","key":"value"}`) {
		t.Errorf("output %s does not contain the context", got)
	}

	// Check that the ERROR entry has an stacktrace key
	if !strings.Contains(got, "stacktrace") {
		t.Errorf("output %s does not contain a stacktrace key", got)
	}
}

func TestLoggerErrorWithoutContext(t *testing.T) {
	initConfig(DEBUG, "my-app", "1.0")

	buf := new(bytes.Buffer)

	log := New().SetWriter(buf)

	log.Error("ERROR message")
	expected := fmt.Sprintf(`{"severity":"ERROR","eventTime":"%s","message":"ERROR message","serviceContext":{"service":"my-app","version":"1.0"},"context":{"reportLocation"`, time.Now().Format(time.RFC3339))
	got := strings.TrimRight(buf.String(), "\n")
	if !strings.Contains(got, expected) {
		t.Errorf("output %s does not containsubstring %s", got, expected)
	}

	// Check that the ERROR entry contains the context
	if strings.Contains(got, `"context":{"data":`) {
		t.Errorf("output %s has a context and it wasn't supposed to", got)
	}

	// Check that the ERROR entry has an stacktrace key
	if !strings.Contains(got, "stacktrace") {
		t.Errorf("output file %s does not contain a stacktrace key", got)
	}
}

func TestLoggerErrorf(t *testing.T) {
	initConfig(DEBUG, "my-app", "1.0")

	buf := new(bytes.Buffer)

	log := New().With(Fields{
		"key":      "value",
		"function": "TestLoggerError",
	}).SetWriter(buf)

	param := "with param"
	log.Errorf("ERROR message %s", param)
	expected := fmt.Sprintf(`{"severity":"ERROR","eventTime":"%s","message":"ERROR message with param","serviceContext":{"service":"my-app","version":"1.0"},"context":{"data":{"function":"TestLoggerError","key":"value"},"reportLocation"`, time.Now().Format(time.RFC3339))
	got := strings.TrimRight(buf.String(), "\n")
	if !strings.Contains(got, expected) {
		t.Errorf("output %s does not containsubstring %s", got, expected)
	}
}

func TestLoggerInfoWithSeveralContextEntries(t *testing.T) {
	initConfig(DEBUG, "my-app", "1.0")

	buf := new(bytes.Buffer)

	log := New().With(Fields{
		"function": "TestLoggerInfo",
		"key":      "value",
		"package":  "logger",
	}).SetWriter(buf)

	log.Info("INFO message")
	expected := fmt.Sprintf(`{"severity":"INFO","eventTime":"%s","message":"INFO message","serviceContext":{"service":"my-app","version":"1.0"},"context":{"data":{"function":"TestLoggerInfo","key":"value","package":"logger"}}}`, time.Now().Format(time.RFC3339))
	got := strings.TrimRight(buf.String(), "\n")
	if expected != got {
		t.Errorf("output file %s does not match expected string %s", got, expected)
	}
}

func TestLoggerErrorWithSeveralContextEntries(t *testing.T) {
	initConfig(DEBUG, "my-app", "1.0")

	buf := new(bytes.Buffer)

	log := New().With(Fields{
		"function": "TestLoggerError",
		"key":      "value",
		"package":  "logger",
	}).SetWriter(buf)

	log.Error("ERROR message")
	expected := fmt.Sprintf(`{"severity":"ERROR","eventTime":"%s","message":"ERROR message","serviceContext":{"service":"my-app","version":"1.0"}`, time.Now().Format(time.RFC3339))
	got := strings.TrimRight(buf.String(), "\n")
	if !strings.Contains(got, expected) {
		t.Errorf("output %s does not containsubstring %s", got, expected)
	}

	// Check that the ERROR entry contains the context
	if !strings.Contains(got, `"context":{"data":{"function":"TestLoggerError","key":"value","package":"logger"}`) {
		t.Errorf("output file %s does not contain the context", got)
	}

	// Check that the ERROR entry has an stacktrace key
	if !strings.Contains(got, "stacktrace") {
		t.Errorf("output file %s does not contain a stacktrace key", got)
	}
}

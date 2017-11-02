package logger

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestLoggerInfoWithOneTimeContext(t *testing.T) {
	initConfig(debug, "robokiller-ivr", "1.0")

	buf := new(bytes.Buffer)

	log := New().With(Fields{
		"key":      "value",
		"function": "TestLoggerDebug",
	}).SetWriter(buf)

	log.Info("info message")
	expected := fmt.Sprintf("{\"severity\":\"INFO\",\"eventTime\":\"%s\",\"message\":\"info message\",\"serviceContext\":{\"service\":\"robokiller-ivr\",\"version\":\"1.0\"},\"context\":{\"data\":{\"function\":\"TestLoggerDebug\",\"key\":\"value\"}}}", time.Now().Format(time.RFC3339))
	got := strings.TrimRight(buf.String(), "\n")
	if expected != got {
		t.Errorf("output %s does not match expected string %s", got, expected)
	}

	// Clean-up the buffer in preparation for new assertions
	buf.Reset()

	log.With(Fields{"foo": "bar"}).SetWriter(buf).Info("unique info message")
	expected = fmt.Sprintf("{\"severity\":\"INFO\",\"eventTime\":\"%s\",\"message\":\"unique info message\",\"serviceContext\":{\"service\":\"robokiller-ivr\",\"version\":\"1.0\"},\"context\":{\"data\":{\"foo\":\"bar\"}}}", time.Now().Format(time.RFC3339))
	got = strings.TrimRight(buf.String(), "\n")
	if expected != got {
		t.Errorf("output file %s does not match expected string %s", got, expected)
	}

	// Clean-up the buffer in preparation for new assertions
	buf.Reset()

	log.SetWriter(buf).Info("unique info message")
	expected = fmt.Sprintf("{\"severity\":\"INFO\",\"eventTime\":\"%s\",\"message\":\"unique info message\",\"serviceContext\":{\"service\":\"robokiller-ivr\",\"version\":\"1.0\"},\"context\":{\"data\":{\"function\":\"TestLoggerDebug\",\"key\":\"value\"}}}", time.Now().Format(time.RFC3339))
	got = strings.TrimRight(buf.String(), "\n")
	if expected != got {
		t.Errorf("output %s does not match expected string %s", got, expected)
	}
}

func TestLoggerErrorWithOneTimeContext(t *testing.T) {
	initConfig(debug, "robokiller-ivr", "1.0")

	buf := new(bytes.Buffer)

	log := New().With(Fields{
		"key":      "value",
		"function": "TestLoggerError",
	}).SetWriter(buf)

	log.Error("error message")
	expected := fmt.Sprintf("{\"severity\":\"ERROR\",\"eventTime\":\"%s\",\"message\":\"error message\",\"serviceContext\":{\"service\":\"robokiller-ivr\",\"version\":\"1.0\"},\"context\":{\"data\":{\"function\":\"TestLoggerError\",\"key\":\"value\"},\"reportLocation\"", time.Now().Format(time.RFC3339))
	got := strings.TrimRight(buf.String(), "\n")
	if !strings.Contains(got, expected) {
		t.Errorf("output %s does not contain substring %s", got, expected)
	}

	// Check that the error entry contains the context
	if !strings.Contains(got, "\"context\":{\"data\":{\"function\":\"TestLoggerError\",\"key\":\"value\"}") {
		t.Errorf("output %s does not contain the context", got)
	}

	// Check that the error entry has an stacktrace key
	if !strings.Contains(got, "stacktrace") {
		t.Errorf("output %s does not contain a stacktrace key", got)
	}

	// Clean-up the buffer in preparation for new assertions
	buf.Reset()

	log.With(Fields{"foo": "bar"}).SetWriter(buf).Error("unique error message")
	expected = fmt.Sprintf("{\"severity\":\"ERROR\",\"eventTime\":\"%s\",\"message\":\"unique error message\",\"serviceContext\":{\"service\":\"robokiller-ivr\",\"version\":\"1.0\"},\"context\":{\"data\":{\"foo\":\"bar\"},\"reportLocation\"", time.Now().Format(time.RFC3339))
	got = strings.TrimRight(buf.String(), "\n")
	if !strings.Contains(got, expected) {
		t.Errorf("output %s does not contain substring %s", got, expected)
	}

	// Check that the error entry contains the context
	if !strings.Contains(got, "\"context\":{\"data\":{\"foo\":\"bar\"}") {
		t.Errorf("output %s does not contain the context", got)
	}

	// Check that the error entry has an stacktrace key
	if !strings.Contains(got, "stacktrace") {
		t.Errorf("output %s does not contain a stacktrace key", got)
	}

	// Clean-up the buffer in preparation for new assertions
	buf.Reset()

	log.SetWriter(buf).Error("unique error message")
	expected = fmt.Sprintf("{\"severity\":\"ERROR\",\"eventTime\":\"%s\",\"message\":\"unique error message\",\"serviceContext\":{\"service\":\"robokiller-ivr\",\"version\":\"1.0\"},\"context\":{\"data\":{\"function\":\"TestLoggerError\",\"key\":\"value\"},\"reportLocation\"", time.Now().Format(time.RFC3339))
	got = strings.TrimRight(buf.String(), "\n")
	if !strings.Contains(got, expected) {
		t.Errorf("output %s does not contain substring %s", got, expected)
	}

	// Check that the error entry contains the context
	if !strings.Contains(got, "\"context\":{\"data\":{\"function\":\"TestLoggerError\",\"key\":\"value\"}") {
		t.Errorf("output %s does not contain the context", got)
	}

	// Check that the error entry has an stacktrace key
	if !strings.Contains(got, "stacktrace") {
		t.Errorf("output %s does not contain a stacktrace key", got)
	}
}

func TestLoggerWithDifferentLogLevels(t *testing.T) {
	initConfig(warn, "robokiller-ivr", "1.0")

	buf := new(bytes.Buffer)

	log := New().With(Fields{
		"key": "value",
	}).SetWriter(buf)

	// LogLevel set to warn, DEBUG messages should not be output
	log.Debug("debug message")
	got := strings.TrimRight(buf.String(), "\n")

	if got != "" {
		t.Errorf("output %s does not match empty string", got)
	}

	// LogLevel set to warn, INFO messages should not be output
	log.Info("info message")
	got = strings.TrimRight(buf.String(), "\n")

	if got != "" {
		t.Errorf("output %s does not match empty string", got)
	}

	log.Warn("warn message")
	expected := fmt.Sprintf("{\"severity\":\"WARN\",\"eventTime\":\"%s\",\"message\":\"warn message\",\"serviceContext\":{\"service\":\"robokiller-ivr\",\"version\":\"1.0\"},\"context\":{\"data\":{\"key\":\"value\"}}}", time.Now().Format(time.RFC3339))
	got = strings.TrimRight(buf.String(), "\n")
	if expected != got {
		t.Errorf("output %s does not match expected string %s", got, expected)
	}

	// Clean-up the buffer in preparation for new assertions
	buf.Reset()

	// should print error as well
	log.Error("error message")
	expected = fmt.Sprintf("{\"severity\":\"ERROR\",\"eventTime\":\"%s\",\"message\":\"error message\",\"serviceContext\":{\"service\":\"robokiller-ivr\",\"version\":\"1.0\"},\"context\":{\"data\":{\"function\":\"TestLoggerError\",\"key\":\"value\"},\"reportLocation\"", time.Now().Format(time.RFC3339))
	got = strings.TrimRight(buf.String(), "\n")
	if strings.Contains(got, expected) {
		t.Errorf("expecting %s; got %s", expected, got)
	}
}

func TestLoggerDebugWithImplicitContext(t *testing.T) {
	initConfig(debug, "robokiller-ivr", "1.0")

	buf := new(bytes.Buffer)

	log := New().With(Fields{
		"key":      "value",
		"function": "TestLoggerDebug",
	}).SetWriter(buf)

	log.Debug("debug message")

	expected := fmt.Sprintf("{\"severity\":\"DEBUG\",\"eventTime\":\"%s\",\"message\":\"debug message\",\"serviceContext\":{\"service\":\"robokiller-ivr\",\"version\":\"1.0\"},\"context\":{\"data\":{\"function\":\"TestLoggerDebug\",\"key\":\"value\"}}}", time.Now().Format(time.RFC3339))
	got := strings.TrimRight(buf.String(), "\n")
	if expected != got {
		t.Errorf("output %s does not match expected string %s", got, expected)
	}
}

func TestLoggerDebugWithoutContext(t *testing.T) {
	initConfig(debug, "robokiller-ivr", "1.0")

	buf := new(bytes.Buffer)
	log := New().SetWriter(buf)

	log.Debug("debug message")
	expected := fmt.Sprintf("{\"severity\":\"DEBUG\",\"eventTime\":\"%s\",\"message\":\"debug message\",\"serviceContext\":{\"service\":\"robokiller-ivr\",\"version\":\"1.0\"}}", time.Now().Format(time.RFC3339))
	got := strings.TrimRight(buf.String(), "\n")
	if expected != got {
		t.Errorf("output %s does not match expected string %s", got, expected)
	}
}

func TestLoggerDebugfWithoutContext(t *testing.T) {
	initConfig(debug, "robokiller-ivr", "1.0")

	buf := new(bytes.Buffer)

	log := New().SetWriter(buf)

	param := "with param"
	log.Debugf("debug message %s", param)
	expected := fmt.Sprintf("{\"severity\":\"DEBUG\",\"eventTime\":\"%s\",\"message\":\"debug message with param\",\"serviceContext\":{\"service\":\"robokiller-ivr\",\"version\":\"1.0\"}}", time.Now().Format(time.RFC3339))
	got := strings.TrimRight(buf.String(), "\n")
	if expected != got {
		t.Errorf("output %s does not match expected string %s", got, expected)
	}
}

func TestLoggerMetric(t *testing.T) {
	initConfig(debug, "robokiller-ivr", "1.0")

	buf := new(bytes.Buffer)

	log := New().SetWriter(buf)

	log.Metric("custom_metric")
	expected := fmt.Sprintf("{\"severity\":\"INFO\",\"eventTime\":\"%s\",\"message\":\"custom_metric\",\"serviceContext\":{\"service\":\"robokiller-ivr\",\"version\":\"1.0\"}}", time.Now().Format(time.RFC3339))
	got := strings.TrimRight(buf.String(), "\n")
	if expected != got {
		t.Errorf("output %s does not match expected string %s", got, expected)
	}
}

func TestLoggerInfo(t *testing.T) {
	initConfig(debug, "robokiller-ivr", "1.0")

	buf := new(bytes.Buffer)

	log := New().With(Fields{
		"key":      "value",
		"function": "TestLoggerInfo",
	}).SetWriter(buf)

	log.Info("info message")
	expected := fmt.Sprintf("{\"severity\":\"INFO\",\"eventTime\":\"%s\",\"message\":\"info message\",\"serviceContext\":{\"service\":\"robokiller-ivr\",\"version\":\"1.0\"},\"context\":{\"data\":{\"function\":\"TestLoggerInfo\",\"key\":\"value\"}}}", time.Now().Format(time.RFC3339))
	got := strings.TrimRight(buf.String(), "\n")
	if expected != got {
		t.Errorf("output %s does not match expected string %s", got, expected)
	}
}

func TestLoggerInfof(t *testing.T) {
	initConfig(debug, "robokiller-ivr", "1.0")

	buf := new(bytes.Buffer)

	log := New().With(Fields{
		"key":      "value",
		"function": "TestLoggerInfo",
	}).SetWriter(buf)

	param := "with param"
	log.Infof("info message %s", param)
	expected := fmt.Sprintf("{\"severity\":\"INFO\",\"eventTime\":\"%s\",\"message\":\"info message with param\",\"serviceContext\":{\"service\":\"robokiller-ivr\",\"version\":\"1.0\"},\"context\":{\"data\":{\"function\":\"TestLoggerInfo\",\"key\":\"value\"}}}", time.Now().Format(time.RFC3339))
	got := strings.TrimRight(buf.String(), "\n")
	if expected != got {
		t.Errorf("output %s does not match expected string %s", got, expected)
	}
}

func TestLoggerError(t *testing.T) {
	initConfig(debug, "robokiller-ivr", "1.0")

	buf := new(bytes.Buffer)

	log := New().With(Fields{
		"key":      "value",
		"function": "TestLoggerError",
	}).SetWriter(buf)

	log.Error("error message")
	expected := fmt.Sprintf("{\"severity\":\"ERROR\",\"eventTime\":\"%s\",\"message\":\"error message\",\"serviceContext\":{\"service\":\"robokiller-ivr\",\"version\":\"1.0\"},\"context\":{\"data\":{\"function\":\"TestLoggerError\",\"key\":\"value\"},\"reportLocation\"", time.Now().Format(time.RFC3339))
	got := strings.TrimRight(buf.String(), "\n")
	if !strings.Contains(got, expected) {
		t.Errorf("output %s does not containsubstring %s", got, expected)
	}

	// Check that the error entry contains the context
	if !strings.Contains(got, "\"context\":{\"data\":{\"function\":\"TestLoggerError\",\"key\":\"value\"}") {
		t.Errorf("output %s does not contain the context", got)
	}

	// Check that the error entry has an stacktrace key
	if !strings.Contains(got, "stacktrace") {
		t.Errorf("output %s does not contain a stacktrace key", got)
	}
}

func TestLoggerErrorWithoutContext(t *testing.T) {
	initConfig(debug, "robokiller-ivr", "1.0")

	buf := new(bytes.Buffer)

	log := New().SetWriter(buf)

	log.Error("error message")
	expected := fmt.Sprintf("{\"severity\":\"ERROR\",\"eventTime\":\"%s\",\"message\":\"error message\",\"serviceContext\":{\"service\":\"robokiller-ivr\",\"version\":\"1.0\"},\"context\":{\"reportLocation\"", time.Now().Format(time.RFC3339))
	got := strings.TrimRight(buf.String(), "\n")
	if !strings.Contains(got, expected) {
		t.Errorf("output %s does not containsubstring %s", got, expected)
	}

	// Check that the error entry contains the context
	if strings.Contains(got, "\"context\":{\"data\":") {
		t.Errorf("output %s has a context and it wasn't supposed to", got)
	}

	// Check that the error entry has an stacktrace key
	if !strings.Contains(got, "stacktrace") {
		t.Errorf("output file %s does not contain a stacktrace key", got)
	}
}

func TestLoggerErrorf(t *testing.T) {
	initConfig(debug, "robokiller-ivr", "1.0")

	buf := new(bytes.Buffer)

	log := New().With(Fields{
		"key":      "value",
		"function": "TestLoggerError",
	}).SetWriter(buf)

	param := "with param"
	log.Errorf("error message %s", param)
	expected := fmt.Sprintf("{\"severity\":\"ERROR\",\"eventTime\":\"%s\",\"message\":\"error message with param\",\"serviceContext\":{\"service\":\"robokiller-ivr\",\"version\":\"1.0\"},\"context\":{\"data\":{\"function\":\"TestLoggerError\",\"key\":\"value\"},\"reportLocation\"", time.Now().Format(time.RFC3339))
	got := strings.TrimRight(buf.String(), "\n")
	if !strings.Contains(got, expected) {
		t.Errorf("output %s does not containsubstring %s", got, expected)
	}
}

func TestLoggerInfoWithSeveralContextEntries(t *testing.T) {
	initConfig(debug, "robokiller-ivr", "1.0")

	buf := new(bytes.Buffer)

	log := New().With(Fields{
		"function": "TestLoggerInfo",
		"key":      "value",
		"package":  "logger",
	}).SetWriter(buf)

	log.Info("info message")
	expected := fmt.Sprintf("{\"severity\":\"INFO\",\"eventTime\":\"%s\",\"message\":\"info message\",\"serviceContext\":{\"service\":\"robokiller-ivr\",\"version\":\"1.0\"},\"context\":{\"data\":{\"function\":\"TestLoggerInfo\",\"key\":\"value\",\"package\":\"logger\"}}}", time.Now().Format(time.RFC3339))
	got := strings.TrimRight(buf.String(), "\n")
	if expected != got {
		t.Errorf("output file %s does not match expected string %s", got, expected)
	}
}

func TestLoggerErrorWithSeveralContextEntries(t *testing.T) {
	initConfig(debug, "robokiller-ivr", "1.0")

	buf := new(bytes.Buffer)

	log := New().With(Fields{
		"function": "TestLoggerError",
		"key":      "value",
		"package":  "logger",
	}).SetWriter(buf)

	log.Error("error message")
	expected := fmt.Sprintf("{\"severity\":\"ERROR\",\"eventTime\":\"%s\",\"message\":\"error message\",\"serviceContext\":{\"service\":\"robokiller-ivr\",\"version\":\"1.0\"}", time.Now().Format(time.RFC3339))
	got := strings.TrimRight(buf.String(), "\n")
	if !strings.Contains(got, expected) {
		t.Errorf("output %s does not containsubstring %s", got, expected)
	}

	// Check that the error entry contains the context
	if !strings.Contains(got, "\"context\":{\"data\":{\"function\":\"TestLoggerError\",\"key\":\"value\",\"package\":\"logger\"}") {
		t.Errorf("output file %s does not contain the context", got)
	}

	// Check that the error entry has an stacktrace key
	if !strings.Contains(got, "stacktrace") {
		t.Errorf("output file %s does not contain a stacktrace key", got)
	}
}

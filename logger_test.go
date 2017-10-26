package logger

import (
	"testing"
	"os"
	"io/ioutil"
	"fmt"
	"time"
	"strings"
)

type LogFunc func(message string)

const OUTFILE = "out.json"

func setEnv() {
	os.Setenv("SERVICE", "robokiller-ivr")
	os.Setenv("VERSION", "1.0")
}

func createOutFile() *os.File {
	// Delete file first if exists
	os.Remove(OUTFILE)

	file, err := os.Create(OUTFILE)
	if err != nil {
		panic("Unable to create test file")
	}

	return file
}

func compareWithOutFile(expected string) bool {
	data, err := ioutil.ReadFile(OUTFILE)
	if err != nil {
		panic("Unable to read test file")
	}

	return strings.TrimRight(string(data), "\n") == expected
}

func outFileContains(substring string) bool {
	data, err := ioutil.ReadFile(OUTFILE)
	if err != nil {
		panic("Unable to read test file")
	}

	fileData := strings.TrimRight(string(data), "\n")
	return strings.Contains(fileData, substring)
}

func TestLoggerDebug(t *testing.T) {
	file := createOutFile()
	defer file.Close()

	setEnv()
	log := New()
	log.SetWriter(file)

	log.Set("key", "value")
	log.Debug("debug message", Fields{"function": "TestLoggerDebug"})
	expected := fmt.Sprintf("{\"severity\":\"DEBUG\",\"eventTime\":\"%s\",\"message\":\"debug message\",\"data\":{\"function\":\"TestLoggerDebug\"},\"serviceContext\":{\"service\":\"robokiller-ivr\",\"version\":\"1.0\"},\"context\":{\"data\":{\"key\":\"value\"}}}", time.Now().Format(time.RFC3339))
	if !compareWithOutFile(expected) {
		t.Errorf("output file %s does not match expected string %s", OUTFILE, expected)
	}
}

func TestLoggerMetric(t *testing.T) {
	file := createOutFile()
	defer file.Close()

	setEnv()
	log := New()
	log.SetWriter(file)

	log.Metric("custom_metric")
	expected := fmt.Sprintf("{\"severity\":\"INFO\",\"eventTime\":\"%s\",\"message\":\"custom_metric\",\"serviceContext\":{\"service\":\"robokiller-ivr\",\"version\":\"1.0\"}}", time.Now().Format(time.RFC3339))
	if !compareWithOutFile(expected) {
		t.Errorf("output file %s does not match expected string %s", OUTFILE, expected)
	}
}

func TestLoggerInfo(t *testing.T) {
	file := createOutFile()
	defer file.Close()

	setEnv()
	log := New()
	log.SetWriter(file)

	log.Set("key", "value")
	log.Info("info message", Fields{"function": "TestLoggerInfo"})
	expected := fmt.Sprintf("{\"severity\":\"INFO\",\"eventTime\":\"%s\",\"message\":\"info message\",\"data\":{\"function\":\"TestLoggerInfo\"},\"serviceContext\":{\"service\":\"robokiller-ivr\",\"version\":\"1.0\"},\"context\":{\"data\":{\"key\":\"value\"}}}", time.Now().Format(time.RFC3339))
	if !compareWithOutFile(expected) {
		t.Errorf("output file %s does not match expected string %s", OUTFILE, expected)
	}
}

func TestLoggerError(t *testing.T) {
	file := createOutFile()
	defer file.Close()

	setEnv()
	log := New()
	log.SetWriter(file)

	log.Set("key", "value")
	log.Error("error message", Fields{"function": "TestLoggerError"})
	expected := fmt.Sprintf("{\"severity\":\"ERROR\",\"eventTime\":\"%s\",\"message\":\"error message\",\"data\":{\"function\":\"TestLoggerError\"},\"serviceContext\":{\"service\":\"robokiller-ivr\",\"version\":\"1.0\"}", time.Now().Format(time.RFC3339))
	if !outFileContains(expected) {
		t.Errorf("output file %s does not containsubstring %s", OUTFILE, expected)
	}

	// Check that the error entry contains the payload
	if !outFileContains("\"data\":{\"function\":\"TestLoggerError\"}") {
		t.Errorf("output file %s does not contain a data key", OUTFILE)
	}

	// Check that the error entry contains the context
	if !outFileContains("\"context\":{\"data\":{\"key\":\"value\"}") {
		t.Errorf("output file %s does not contain the context", OUTFILE)
	}

	// Check that the error entry has an stacktrace key
	if !outFileContains("stacktrace") {
		t.Errorf("output file %s does not contain a stacktrace key", OUTFILE)
	}
}

func TestLoggerInfoWithSeveralPayloadEntries(t *testing.T) {
	file := createOutFile()
	defer file.Close()

	setEnv()
	log := New()
	log.SetWriter(file)

	log.Set("key", "value")
	log.Info("info message", Fields{"function": "TestLoggerInfo", "package": "logger"})
	expected := fmt.Sprintf("{\"severity\":\"INFO\",\"eventTime\":\"%s\",\"message\":\"info message\",\"data\":{\"function\":\"TestLoggerInfo\",\"package\":\"logger\"},\"serviceContext\":{\"service\":\"robokiller-ivr\",\"version\":\"1.0\"},\"context\":{\"data\":{\"key\":\"value\"}}}", time.Now().Format(time.RFC3339))
	if !compareWithOutFile(expected) {
		t.Errorf("output file %s does not match expected string %s", OUTFILE, expected)
	}
}

func TestLoggerErrorWithSeveralPayloadEntries(t *testing.T) {
	file := createOutFile()
	defer file.Close()

	setEnv()
	log := New()
	log.SetWriter(file)

	log.Set("key", "value")
	log.Error("error message", Fields{"function": "TestLoggerError", "package": "logger"})
	expected := fmt.Sprintf("{\"severity\":\"ERROR\",\"eventTime\":\"%s\",\"message\":\"error message\",\"data\":{\"function\":\"TestLoggerError\",\"package\":\"logger\"},\"serviceContext\":{\"service\":\"robokiller-ivr\",\"version\":\"1.0\"}", time.Now().Format(time.RFC3339))
	if !outFileContains(expected) {
		t.Errorf("output file %s does not containsubstring %s", OUTFILE, expected)
	}

	// Check that the error entry contains the payload
	if !outFileContains("\"data\":{\"function\":\"TestLoggerError\",\"package\":\"logger\"}") {
		t.Errorf("output file %s does not contain a data key", OUTFILE)
	}

	// Check that the error entry contains the context
	if !outFileContains("\"context\":{\"data\":{\"key\":\"value\"}") {
		t.Errorf("output file %s does not contain the context", OUTFILE)
	}

	// Check that the error entry has an stacktrace key
	if !outFileContains("stacktrace") {
		t.Errorf("output file %s does not contain a stacktrace key", OUTFILE)
	}
}

func TestLoggerInfoWithSeveralContextEntries(t *testing.T) {
	file := createOutFile()
	defer file.Close()

	setEnv()
	log := New()
	log.SetWriter(file)

	log.Set("key", "value")
	log.Set("extraKey", "extraValue")
	log.Info("info message", Fields{"function": "TestLoggerInfo", "package": "logger"})
	expected := fmt.Sprintf("{\"severity\":\"INFO\",\"eventTime\":\"%s\",\"message\":\"info message\",\"data\":{\"function\":\"TestLoggerInfo\",\"package\":\"logger\"},\"serviceContext\":{\"service\":\"robokiller-ivr\",\"version\":\"1.0\"},\"context\":{\"data\":{\"extraKey\":\"extraValue\",\"key\":\"value\"}}}", time.Now().Format(time.RFC3339))
	if !compareWithOutFile(expected) {
		t.Errorf("output file %s does not match expected string %s", OUTFILE, expected)
	}
}

func TestLoggerErrorWithSeveralContextEntries(t *testing.T) {
	file := createOutFile()
	defer file.Close()

	setEnv()
	log := New()
	log.SetWriter(file)

	log.Set("key", "value")
	log.Set("extraKey", "extraValue")
	log.Error("error message", Fields{"function": "TestLoggerError", "package": "logger"})
	expected := fmt.Sprintf("{\"severity\":\"ERROR\",\"eventTime\":\"%s\",\"message\":\"error message\",\"data\":{\"function\":\"TestLoggerError\",\"package\":\"logger\"},\"serviceContext\":{\"service\":\"robokiller-ivr\",\"version\":\"1.0\"}", time.Now().Format(time.RFC3339))
	if !outFileContains(expected) {
		t.Errorf("output file %s does not containsubstring %s", OUTFILE, expected)
	}

	// Check that the error entry contains the payload
	if !outFileContains("\"data\":{\"function\":\"TestLoggerError\",\"package\":\"logger\"}") {
		t.Errorf("output file %s does not contain a data key", OUTFILE)
	}

	// Check that the error entry contains the context
	if !outFileContains("\"context\":{\"data\":{\"extraKey\":\"extraValue\",\"key\":\"value\"}") {
		t.Errorf("output file %s does not contain the context", OUTFILE)
	}

	// Check that the error entry has an stacktrace key
	if !outFileContains("stacktrace") {
		t.Errorf("output file %s does not contain a stacktrace key", OUTFILE)
	}
}

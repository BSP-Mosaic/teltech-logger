# teltech/logger

Super simple structured logging mechanism for Go projects with [Stackdriver format](https://cloud.google.com/error-reporting/docs/formatting-error-messages) compatibility

## Installation

``` sh
go get -u github.com/teltech/logger
```

## Usage
``` go
package main

import (
    "github.com/teltech/logger"
)

// There should be a LOG_LEVEL environment variable set, which is read by the library
// If no value is set, the default LOG_LEVEL will be INFO

func main() {
    // Stackdriver requires a project name and version to be set. Use your environment for these values.
    // SERVICE should be your GCP project-id, e.g. robokiller-146813
    // VERSION is an arbitrary value
    log, err := log.New()
    if err != nil {
        fmt.Errorf("cannot initiate logger")
    }

    // A metric is an INFO log entry without a payload
    log.Metric("CUSTOM_METRIC_ENTRY")

    // Add context values to all subsequent log entries using Set(), the values will persisted for the scope of the logger instance
    log.Set("user", "+1234567890")
    log.Set("action", "create-account")

    // Log a DEBUG message, only visible in when LOG_LEVEL is set to DEBUG
    log.Debug("debug message goes here", log.Fields{"key":"val"})

    // Log an INFO message
    log.Info("info message goes here", log.Fields{"key":"val"})

    // Log a WARN message
    log.Warn("warn message goes here", log.Fields{"key":"val"})

    // Error() prints the stacktrace as part of the payload for each entry and sends the
    // data to Stackdriver Error Reporting service
    log.Error("error message goes here", log.Fields{"key":"val"})
}
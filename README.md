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
    log := logger.New()

    // You can also initialize the logger with a context, the values will persisted throughout the scope of the logger instance
    log = logger.New().With(logger.Fields{
        "user":   "+1234567890",
        "action": "create-account",
    })

    param := "something useful here"

    // Log a DEBUG message, only visible in when LOG_LEVEL is set to DEBUG
    log.With(logger.Fields{"key": "val", "something": true}).Debug("debug message goes here")
    log.With(logger.Fields{"key": "val"}).Debugf("debug message with %s", param)

    // Log an INFO message, should be used for metrics as well
    log.Info("CUSTOM_METRIC")
    log.With(logger.Fields{"key": "val", "names": []string{"Mauricio", "Manuel"}}).Info("info message goes here")
    log.With(logger.Fields{"key": "val"}).Infof("info message with %s", param)

    // Log a WARN message
    log.With(logger.Fields{"key": "val"}).Warn("warn message goes here")
    log.With(logger.Fields{"key": "val"}).Warnf("warn message with %s", param)

    // Error() prints the stacktrace as part of the payload for each entry and sends the
    // data to Stackdriver Error Reporting service
    log.With(logger.Fields{"key": "val"}).Error("error message goes here")
    log.With(logger.Fields{"key": "val"}).Errorf("error message with %s", param)
}
```

## Output

The errors require a specific JSON format for them to be ingested and processed by Google Cloud Platform Stackdriver Logging and Error Reporting. See: [https://cloud.google.com/error-reporting/docs/formatting-error-messages](https://cloud.google.com/error-reporting/docs/formatting-error-messages). The resulting output has the following format, optional fields are... well, optional:
```json
       {
          "severity": "ERROR",
          "eventTime": "2017-04-26T02:29:33-04:00",
          "message": "An error just happened!",
          "serviceContext": {
             "service": "robokiller-ivr",
             "version": "1.0"
          },
          "context": {
            "data": {
              "clientIP": "127.0.0.1"
              "userAgent": "Mosaic 1.0"
            },
            "reportLocation": {
              "filePath": "\/Users\/mc\/Documents\/src\/github.com\/macuenca\/apex\/mauricio.go",
              "functionName": "unknown",
              "lineNumber": 15
            }
          },
         "stacktrace": "goroutine 1 [running]:main.main()\n\t\/github.com\/macuenca\/mauricio.go:15 +0x1a9\n"
       }
```

## License

This package is licensed under the [BSD 3-clause](https://opensource.org/licenses/BSD-3-Clause) license.

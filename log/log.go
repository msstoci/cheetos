package log

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

const logFormat = `date=%s, method=%s, url=%s,  response_time=%s`

var logLevel = flag.String("log_level", "info", "set log level")
var version = flag.Bool("v", false, "for version")

var stdoutLog, stderrLog string

func init() {
	flag.StringVar(&stdoutLog, "loginfo", "", "log file for stdout")
	flag.StringVar(&stderrLog, "logerror", "", "log file for stderr")

	if *version {
		os.Exit(0)
		return
	}
	SetLevel(*logLevel)
}

// New create logger instance with request_id as a field
// also set the default format to JSON.
func New() *logrus.Entry {
	baseLogger := logrus.WithField("request_id", uuid.New().String())
	baseLogger.Logger.SetFormatter(&logrus.JSONFormatter{})
	return baseLogger
}

func LogInit() {

	if stdoutLog != stderrLog && stdoutLog != "" {
		logrus.Println("Log Init: using ", stdoutLog, stderrLog)
	}

	reopen(1, stdoutLog)
	reopen(2, stderrLog)

	setupLogs()

}

// WriterHook is a hook that writes logs of specified LogLevels to specified Writer
type WriterHook struct {
	Writer    io.Writer
	LogLevels []logrus.Level
}

// Fire will be called when some logging function is called with current hook
// It will format log entry to string and write it to appropriate writer
func (hook *WriterHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}
	_, err = hook.Writer.Write([]byte(line))
	return err
}

// Levels define on which log levels this hook would trigger
func (hook *WriterHook) Levels() []logrus.Level {
	return hook.LogLevels
}

// setupLogs adds hooks to send logs to different destinations depending on level
func setupLogs() {
	logrus.SetOutput(ioutil.Discard) // Send all logs to nowhere by default

	logrus.AddHook(&WriterHook{ // Send logs with level higher than warning to stderr
		Writer: os.Stderr,
		LogLevels: []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel,
		},
	})
	logrus.AddHook(&WriterHook{ // Send info and debug logs to stdout
		Writer: os.Stdout,
		LogLevels: []logrus.Level{
			logrus.InfoLevel,
			logrus.DebugLevel,
		},
	})
}

func reopen(fd int, filename string) {

	if filename == "" {
		return
	}

	logFile, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)

	if err != nil {
		logrus.Println("Error in opening ", filename, err)
		os.Exit(2)
	}

	if err = syscall.Dup2(int(logFile.Fd()), fd); err != nil {
		logrus.Println("Failed to dup", filename)
	}
}

type Fields logrus.Fields

func SetLevel(level string) {
	switch level {
	case "panic":
		logrus.SetLevel(logrus.PanicLevel)
	case "fatal":
		logrus.SetLevel(logrus.FatalLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	case "warning":
		logrus.SetLevel(logrus.WarnLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}
}

func GetLevel() string {
	return strings.ToUpper(logrus.GetLevel().String())
}

func Request(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		t := time.Now()
		next(w, r, ps)
		Infof(logFormat, t, r.Method, r.RequestURI, time.Since(t))
	}
}

func Info(args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		file = file[slash+1:]
	}
	logrus.WithField("source", fmt.Sprintf("%s:%d", file, line)).Info(args...)
}

func Infoln(args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		file = file[slash+1:]
	}
	logrus.WithField("source", fmt.Sprintf("%s:%d", file, line)).Infoln(args...)
}

func Infof(format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		file = file[slash+1:]
	}
	logrus.WithField("source", fmt.Sprintf("%s:%d", file, line)).Infof(format, args...)
}

func Print(args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		file = file[slash+1:]
	}
	logrus.WithField("source", fmt.Sprintf("%s:%d", file, line)).Info(args...)
}

func Println(args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		file = file[slash+1:]
	}
	logrus.WithField("source", fmt.Sprintf("%s:%d", file, line)).Infoln(args...)
}

func Printf(format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		file = file[slash+1:]
	}
	logrus.WithField("source", fmt.Sprintf("%s:%d", file, line)).Infof(format, args...)
}

func Debug(args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		file = file[slash+1:]
	}
	logrus.WithField("source", fmt.Sprintf("%s:%d", file, line)).Debug(args...)
}

func Debugln(args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		file = file[slash+1:]
	}
	logrus.WithField("source", fmt.Sprintf("%s:%d", file, line)).Debugln(args...)
}

func Debugf(format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		file = file[slash+1:]
	}
	logrus.WithField("source", fmt.Sprintf("%s:%d", file, line)).Debugf(format, args...)
}

func Warn(args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		file = file[slash+1:]
	}
	logrus.WithField("source", fmt.Sprintf("%s:%d", file, line)).Warn(args...)
}

func Warnln(args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		file = file[slash+1:]
	}
	logrus.WithField("source", fmt.Sprintf("%s:%d", file, line)).Warnln(args...)
}

func Warnf(format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		file = file[slash+1:]
	}
	logrus.WithField("source", fmt.Sprintf("%s:%d", file, line)).Warnf(format, args...)
}

func Error(args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		file = file[slash+1:]
	}
	logrus.WithField("source", fmt.Sprintf("%s:%d", file, line)).Error(args...)
}

func Errorln(args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		file = file[slash+1:]
	}
	logrus.WithField("source", fmt.Sprintf("%s:%d", file, line)).Errorln(args...)
}

func Errorf(format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		file = file[slash+1:]
	}
	logrus.WithField("source", fmt.Sprintf("%s:%d", file, line)).Errorf(format, args...)
}

func Fatal(args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		file = file[slash+1:]
	}
	logrus.WithField("source", fmt.Sprintf("%s:%d", file, line)).Fatal(args...)
}

func Fatalln(args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		file = file[slash+1:]
	}
	logrus.WithField("source", fmt.Sprintf("%s:%d", file, line)).Fatalln(args...)
}

func Fatalf(format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		file = file[slash+1:]
	}
	logrus.WithField("source", fmt.Sprintf("%s:%d", file, line)).Fatalf(format, args...)
}

func WithFields(fields Fields) *logrus.Entry {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		file = file[slash+1:]
	}

	fields["source"] = fmt.Sprintf("%s:%d", file, line)

	logrusFields := logrus.Fields{}

	for key, value := range fields {
		logrusFields[key] = value
	}

	return logrus.WithFields(logrusFields)
}

func WithError(err error) *logrus.Entry {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		file = file[slash+1:]
	}

	fields := logrus.Fields{
		"source": fmt.Sprintf("%s:%d", file, line),
		"error":  err,
	}

	return logrus.WithFields(fields)
}
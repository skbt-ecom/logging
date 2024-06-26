package logging

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"runtime"
)

var e *logrus.Entry

type Logger struct {
	*logrus.Entry
}

func InitLogger() *Logger {
	return &Logger{e}
}

func init() {
	//Logger instance creating
	l := logrus.New()

	//Logger instance settings
	l.SetLevel(logrus.TraceLevel)
	l.SetOutput(io.Discard)
	l.SetReportCaller(true)

	l.Formatter = &logrus.TextFormatter{
		DisableColors:   true,
		FullTimestamp:   true,
		TimestampFormat: "03:04:05 02/01/2006",
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			pc := make([]uintptr, 4)
			count := runtime.Callers(10, pc)

			frames := runtime.CallersFrames(pc[:count])

			frames.Next()
			frame, _ := frames.Next()
			return fmt.Sprintf("%s()", frame.Function), fmt.Sprintf("%s:%d", path.Base(frame.File), frame.Line)
		},
	}

	l.AddHook(&writerHook{
		Writer:    []io.Writer{os.Stdout},
		LogLevels: logrus.AllLevels,
	})

	e = logrus.NewEntry(l)
}

func (l *Logger) SetLevel(logLevel string) {
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		l.WithError(err).Fatalf("failed to set logger level: %s", logLevel)
	}

	l.Logger.SetLevel(level)
}

func (l *Logger) WithExtraField(key string, value any) *Logger {
	return &Logger{l.WithField(key, value)}
}

type Fields logrus.Fields

func (l *Logger) WithExtraFields(fields Fields) *Logger {
	return &Logger{l.WithFields(logrus.Fields(fields))}
}

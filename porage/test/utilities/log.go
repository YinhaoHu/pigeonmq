package utilities

import (
	"fmt"
	"sync"
	"time"

	"github.com/fatih/color"
)

type TestLogger struct {
	mutex sync.Mutex
}

var Logger *TestLogger = makeTestLogger()

func makeTestLogger() *TestLogger {
	return &TestLogger{
		mutex: sync.Mutex{},
	}
}

func (l *TestLogger) Logf(format string, args ...interface{}) {
	timeString := time.Now().Format("2006-01-02 15:04:05")
	sectionString := color.HiYellowString("TEST")
	content := fmt.Sprintf(format, args...)
	fmt.Printf("[%s %s] %s\n", timeString, sectionString, content)
}

func (l *TestLogger) FatalIfErr(err error, format string, args ...interface{}) {
	if err != nil {
		l.Logf(format, args...)
		panic(err)
	}
}

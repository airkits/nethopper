package log

import (
	"fmt"
	"runtime"
	"time"

	"github.com/airkits/nethopper/base"
)

const BATCH_LOG_SIZE = 1

func Run(s *LogModule) {
	ctxDone := false
	exitFlag := false
	start := time.Now()
	Info("Module [Logger] starting")
	for {
		s.OnRun(time.Since(start))
		start = time.Now()
		if s.MQ().Length() == 0 {
			t := time.Duration(s.IdleTimes()) * time.Nanosecond
			time.Sleep(t)
			s.IdleTimesAdd()

		}
		if ctxDone, exitFlag = s.CanExit(ctxDone); exitFlag {
			fmt.Printf("module log exit")
			return
		}
		runtime.Gosched()
	}
}

// GLoggerModule global log module
var GLoggerModule *LogModule

// SetGLogger set logger module instance
func InitLogger(conf *Config) {
	GLoggerModule = &LogModule{}
	GLoggerModule.Setup(conf)
	base.GO(Run, GLoggerModule)
}

//WriteLog send log to queue
func WriteLog(level int32, v ...interface{}) error {
	// UserData return logger level
	if GLoggerModule == nil || level > GLoggerModule.UserData() {
		return nil
	}
	msg := FormatLog(level, v...)
	if err := GLoggerModule.PushBytes(level, []byte(msg)); err != nil {
		return err
	}
	return nil
}

// Fatal system is unusable
func Fatal(v ...interface{}) error {
	return WriteLog(FATAL, v...)
}

// Error error conditions
func Error(v ...interface{}) error {
	return WriteLog(ERROR, v...)
}

// Trace error conditions
func Trace(v ...interface{}) error {
	return WriteLog(TRACE, v...)
}

// Warning warning conditions
func Warning(v ...interface{}) error {
	return WriteLog(WARNING, v...)
}

// Info informational messages
func Info(v ...interface{}) error {
	return WriteLog(INFO, v...)
}

// Debug debug-level messages
func Debug(v ...interface{}) error {
	return WriteLog(DEBUG, v...)
}

// TraceCost calc the api cost time
// usage: defer TraceCose("func")()
func TraceCost(msg string) func() {
	start := time.Now()
	return func() {
		Trace("%s [TraceCost] cost (%s)\n", msg, time.Since(start))
	}
}

//PrintStack print current stack
func PrintStack(all bool) {
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, all)

	Fatal("[FATAL] catch a panic,stack is: %s", string(buf[:n]))
}

// GetStack get current stack
func GetStack(all bool) string {
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, all)
	return string(buf[:n])
}

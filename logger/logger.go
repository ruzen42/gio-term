package logger

import (
	"fmt"

	"github.com/go-color-term/go-color-term/coloring"
	"io"
	"os"
	"sync"
	"time"
)

var (
	mu     sync.Mutex
	out    io.Writer = os.Stdout 
	colors           = true     
)

func getTime() string {
	return time.Now().Format("15:04:05")
}

func writeLog(ktype string, colorFn func(string) string, format string, args ...any) {
	mu.Lock()
	defer mu.Unlock()

	msg := fmt.Sprintf(format, args...)
	tag := ktype
	if colors && colorFn != nil {
		tag = colorFn(ktype)
	}

	fmt.Fprintf(out, "%s %s %s\n", getTime(), tag, msg)
}

func Debug(format string, args ...any) { writeLog("DEBUG", coloring.Faint, format, args...) }
func Info(format string, args ...any)  { writeLog("INFO", coloring.Green, format, args...) }
func Warn(format string, args ...any)  { writeLog("WARN", coloring.Yellow, format, args...) }
func Error(format string, args ...any) { writeLog("ERROR", coloring.Red, format, args...) }

func Fatal(format string, args ...any) {
	writeLog("FATAL", coloring.Red, format, args...)
	os.Exit(1)
}

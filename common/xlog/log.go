package xlog

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

func Fatalf(format string, args ...interface{}) {
	caller := getCaller(2)
	m := map[string]interface{}{
		"caller":     caller,
		"@timestamp": time.Now().Format("2006-01-02T15:04:05.000Z07:00"),
		"content":    fmt.Sprintf(format, args...),
		"level":      "fatal",
	}
	msg, _ := json.Marshal(m)
	fmt.Println(string(msg))
	os.Exit(1)
}

func getCaller(callDepth int) string {
	_, file, line, ok := runtime.Caller(callDepth)
	if !ok {
		return ""
	}

	return prettyCaller(file, line)
}

func prettyCaller(file string, line int) string {
	idx := strings.LastIndexByte(file, '/')
	if idx < 0 {
		return fmt.Sprintf("%s:%d", file, line)
	}

	idx = strings.LastIndexByte(file[:idx], '/')
	if idx < 0 {
		return fmt.Sprintf("%s:%d", file, line)
	}

	return fmt.Sprintf("%s:%d", file[idx+1:], line)
}

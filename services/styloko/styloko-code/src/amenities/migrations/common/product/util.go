package product

import (
	"common/notification"
	"common/notification/datadog"
	"fmt"
	"regexp"
	"runtime"
	"strconv"

	"github.com/jabong/floRest/src/common/utils/logger"
)

type generalMap map[string]interface{}

func SanitizeImageName(input string) string {
	illegalChars := regexp.MustCompile(`[^[a-zA-Z0-9\s]]*`)
	output := illegalChars.ReplaceAllString(input, "-")
	return output
}

func SqlSafe(input string) string {
	return "`" + input + "`"
}

func dbToString(x interface{}) (string, bool) {

	out, ok := x.([]uint8)
	if !ok {
		return "", false
	}
	return string(out), true
}

func dbToInt(x interface{}) (int, bool) {

	out, ok := x.([]uint8)
	if !ok {
		return 0, false
	}
	outI, err := strconv.Atoi(string(out))
	if err != nil {
		return 0, false
	}
	return outI, true
}

func printErr(e error) {
	logger.Error(e.Error())
	fmt.Println(e.Error())
}

func printInfo(s string) {
	logger.Info(s)
}

func recoverHandler(handler string) {
	if rec := recover(); rec != nil {
		logger.Error(fmt.Sprintf("[PANIC] occured with %s", handler))
		trace := make([]byte, 4096)
		count := runtime.Stack(trace, true)
		logger.Error(fmt.Sprintf("Stack of %d bytes, trace : \n%s", count, trace))
		logger.Error(fmt.Sprintf("Reason for panic: %s", rec))

		// Notification data below
		// Only pushing 1024 bytes of stack trace to datadog.
		trace = make([]byte, 1024)
		count = runtime.Stack(trace, true)
		stackTrace := fmt.Sprintf("Stack of %d bytes, trace : \n%s", count, trace)
		title := fmt.Sprintf("Panic occured")
		text := fmt.Sprintf("Panic reason %s\n\nStack Trace: %s", rec, stackTrace)
		tags := []string{"pool-error", "panic"}
		notification.SendNotification(title, text, tags, datadog.ERROR)
	}
}

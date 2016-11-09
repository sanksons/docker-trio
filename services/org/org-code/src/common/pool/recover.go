package pool

import (
	"common/notification"
	"common/notification/datadog"
	"fmt"
	"runtime"

	"github.com/jabong/floRest/src/common/utils/logger"
)

func RecoverHandler(handler string) {
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

package utils

import (
	"bytes"
	"common/notification"
	"common/notification/datadog"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jabong/floRest/src/common/utils/logger"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// GetMockPath -> Returns Mock Path
func GetMockPath() string {
	dir, _ := os.Getwd()
	return fmt.Sprintf("%s/../mocks/", strings.Replace(dir, " ", "\\ ", -1))
}

// ToString -> Returns string from interface
func ToString(x interface{}) string {

	switch i := x.(type) {
	case nil:
		return ""
	case float64:
		return strconv.FormatFloat(i, 'f', 2, 64)
	case *float64:
		var tmp float64 = *i
		return strconv.FormatFloat(tmp, 'f', 2, 64)
	case int:
		return strconv.Itoa(i)
	case int64:
		return strconv.Itoa(int(i))
	case string:
		return i
	case []interface{}:
		var sarr []string
		for _, v := range i {
			sarr = append(sarr, ToString(v))
		}
		return strings.Join(sarr, "|")
	default:
		return ""
	}
}

// GetFloat -> Float type conversion
func GetFloat(x interface{}) (float64, error) {
	switch i := x.(type) {
	case nil:
		return 0.00, nil
	case float64:
		return i, nil
	case int:
		return float64(i), nil
	case string:
		return strconv.ParseFloat(i, 64)
	default:
		return 0.00, errors.New("Cannot convert to float")
	}
}

// GetInt -> Int type conversion

func GetInt(x interface{}) (int, error) {
	switch i := x.(type) {
	case nil:
		return 0, nil
	case int:
		return i, nil
	case string:
		return strconv.Atoi(i)
	case float64:
		return int(i), nil
	case int64:
		return int(i), nil
	default:
		return 0, errors.New("Cannot convert to int")
	}
}

// Convert time to Mysql format
func ToMySqlTime(t *time.Time) (formatted string) {
	if t == nil {
		return
	}
	return t.Format("2006-01-02 15:04:05")
}

func SqlSafe(input string) string {
	return "`" + input + "`"
}

func Random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

func InArrayString(arr []string, needle string) bool {
	for _, s := range arr {
		if strings.ToLower(s) == strings.ToLower(needle) {
			return true
		}
	}
	return false
}

func InArrayInt(arr []int, needle int) bool {
	for _, s := range arr {
		if s == needle {
			return true
		}
	}
	return false
}

func JSONMarshal(v interface{}, safeEncoding bool) ([]byte, error) {
	b, err := json.Marshal(v)

	if safeEncoding {
		b = bytes.Replace(b, []byte("\\u003c"), []byte("<"), -1)
		b = bytes.Replace(b, []byte("\\u003e"), []byte(">"), -1)
		b = bytes.Replace(b, []byte("\\u0026"), []byte("&"), -1)
	}
	return b, err
}

// RecoverHandler -> recovers a panic
func RecoverHandler(handler string) {
	if rec := recover(); rec != nil {
		logger.Error(fmt.Sprintf("[PANIC] occured with %s", handler))
		trace := make([]byte, 4096)
		count := runtime.Stack(trace, true)
		logger.Error(fmt.Sprintf("Stack of %d bytes, trace : \n%s", count, trace))
		logger.Error(fmt.Sprintf("Reason for panic: %s", rec))

		trace = make([]byte, 1024)
		count = runtime.Stack(trace, true)
		stackTrace := fmt.Sprintf("Stack of %d bytes, trace : \n%s", count, trace)
		title := fmt.Sprintf("Panic occured")
		text := fmt.Sprintf("Panic reason %s\n\nStack Trace: %s", rec, stackTrace)
		tags := []string{"error", "panic"}
		notification.SendNotification(title, text, tags, datadog.ERROR)
	}
}

func GetBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

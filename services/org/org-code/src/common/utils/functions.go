package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
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
	case int:
		return strconv.Itoa(i)
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
	default:
		return 0, errors.New("Cannot convert to int")
	}
}

func ConvertStructArrToMapArr(strctArr interface{}) ([]map[string]interface{}, error) {
	m := make([]map[string]interface{}, 1)
	temp, err := json.Marshal(strctArr)
	if err != nil {
		return nil, errors.New("Error while Marshaling Struct")
	}
	err = json.Unmarshal(temp, &m)
	if err != nil {
		return nil, errors.New("Error while Unmarshaling Struct")
	}
	return m, nil
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

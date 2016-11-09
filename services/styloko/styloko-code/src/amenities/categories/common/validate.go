package common

import (
	"common/appconstant"
	"reflect"
	"strconv"
	"strings"

	florest_constants "github.com/jabong/floRest/src/common/constants"
)

func genError(field, msg string) florest_constants.AppError {
	tmp := florest_constants.AppError{
		Code:             appconstant.InvalidDataErrorCode,
		Message:          DATA_VALIDATION_FAIL,
		DeveloperMessage: field + " => " + msg,
	}
	return tmp
}

// ValidateV2 -> Validation function, returns AppErrors , bool.
// Supported tags => required:"true", accepts:"[1,2,3]", validate:"true"
// Currently works fully only for flat structs. Complex structs may fail or cause panic.
// Required input is unmarshalled struct.
// For not required fields but need validation, use validate="true" tag. Use only for stucts and slices.
// TODO Pointers must be handled.
func ValidateV2(container interface{}) (florest_constants.AppErrors, bool) {
	stype := reflect.TypeOf(container)
	sval := reflect.ValueOf(container)
	listError := new(florest_constants.AppErrors)
	flag := true
	for i := 0; i < stype.NumField(); i++ {
		field := stype.Field(i)
		val := sval.Field(i)
		jtag := field.Tag.Get("json")
		name := strings.Split(jtag, ",")[0]
		// First check is required.
		tag := field.Tag.Get("required")
		validate := field.Tag.Get("validate")
		if tag == "true" || validate == "true" {
			switch val.Kind() {
			// This case for recursive struct check is untested.
			// Please use caution.
			case reflect.Struct:
				appErr, ok := ValidateV2(val.Interface())
				flag = ok
				listError.Errors = append(listError.Errors, appErr.Errors...)
				break

			// Slices also call ValidateV2 recursively.
			// Case has been tested.
			case reflect.Array, reflect.Slice:
				for k := 0; k < val.Len(); k++ {
					appErr, ok := ValidateV2(val.Index(k).Interface())
					flag = ok
					listError.Errors = append(listError.Errors, appErr.Errors...)
				}
				break

			case reflect.String:
				data := strings.Trim(val.String(), " ")
				if data == "" {
					flag = false
					listError.Errors = append(listError.Errors, genError(name, "Cannot be empty"))
				}
				break

			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				data := val.Int()
				if data == 0 {
					flag = false
					listError.Errors = append(listError.Errors, genError(name, "Cannot be zero or empty."))
				}
				break

			default:
				break
			}
		}

		// Second check is accepts
		acceptTag := field.Tag.Get("accepts")
		if acceptTag != "" {
			switch val.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				acceptTag = strings.Replace(acceptTag, "[", "", -1)
				acceptTag = strings.Replace(acceptTag, "]", "", -1)
				acceptTag = strings.Trim(acceptTag, " ")
				acceptLs := strings.Split(acceptTag, ",")
				num := val.Int()
				data := strconv.FormatInt(num, 10)
				found := false
				for _, v := range acceptLs {
					if data == v {
						found = true
						break
					}
				}
				if !found {
					flag = false
					listError.Errors = append(listError.Errors, genError(name, "Accpeted values are "+acceptTag))
				}
				break

			default:
				break
			}
		}
	}
	return *listError, flag
}

package common

import (
	categoryService "amenities/services/categories"
	factory "common/ResourceFactory"
	"common/appconfig"
	"common/notification"
	"common/notification/datadog"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/jabong/floRest/src/common/config"
	"github.com/jabong/floRest/src/common/utils/logger"
	validator "gopkg.in/go-playground/validator.v8"
)

// Initialized from PUT API definition
var AttributesInfo map[string]interface{}

// Called from PUT API definition
func ReadAttributeFile() {
	attrInfoFile, err := ioutil.ReadFile(ATTRIBUTES_INFO_FILE)
	if err != nil {
		panic(fmt.Sprintf("Error loading Attributes Info file %s \n %s", ATTRIBUTES_INFO_FILE, err))
	}
	err = json.Unmarshal(attrInfoFile, &AttributesInfo)
	if err != nil {
		panic(fmt.Sprintf("Incorrect Json in %s \n %s", ATTRIBUTES_INFO_FILE, err))
	}
}

// Called from PUT API definition
func ReadTestAttributeFile() {

	attrInfoFile, err := ioutil.ReadFile("/tmp/" + ATTRIBUTES_INFO_FILE)
	if err != nil {
		panic(fmt.Sprintf("Error loading Attributes Info file %s \n %s", ATTRIBUTES_INFO_FILE, err))
	}
	err = json.Unmarshal(attrInfoFile, &AttributesInfo)
	if err != nil {
		panic(fmt.Sprintf("Incorrect Json in %s \n %s", ATTRIBUTES_INFO_FILE, err))
	}
}

var NotFoundErr error = errors.New("Data Not Found.")

var HttpClient = &http.Client{
	Transport: &http.Transport{
		MaxIdleConnsPerHost: 20,
	}}

// Get shipment value
func GetShipmentById(id int) string {
	data := map[int]string{
		1: "Own Warehouse",
		2: "Dropshipping",
		3: "Cross docking",
		4: "Marketplace",
	}
	return data[id]
}

// Get All supported EXPANSE
func GetAllExpanse() []string {
	return []string{
		EXPANSE_CATALOG,
		EXPANSE_XLARGE,
		EXPANSE_LARGE,
		EXPANSE_MEDIUM,
		EXPANSE_SMALL,
		EXPANSE_XSMALL,
		EXPANSE_SOLR,
		EXPANSE_MEMCACHE,
		EXPANSE_PROMOTION,
	}
}

// Get All supported Visibility
func GetAllVisibility() []string {
	return []string{
		VISIBILITY_PDP,
		VISIBILITY_MSKU,
		VISIBILITY_DOOS,
		VISIBILITY_NONE,
	}
}

// Convert time to Mysql format
func ToMySqlTime(t *time.Time) (formatted string) {
	if t == nil {
		return
	}

	return t.Format(FORMAT_MYSQL_TIME)
}

// Convert time to Mysql format, can also return null
func ToMySqlTimeNull(t *time.Time) (formatted *string) {
	if t == nil {
		return nil
	}

	retTime := t.Format(FORMAT_MYSQL_TIME)
	return &retTime
}

// Convert From Mysql format
func FromMysqlTime(s string, localize bool) (*time.Time, error) {
	if s == "" {
		return nil, nil
	}
	var t time.Time
	var err error
	if localize {
		t, err = time.ParseInLocation(FORMAT_MYSQL_TIME, s, time.Local)
	} else {
		t, err = time.Parse(FORMAT_MYSQL_TIME, s)
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func PrepareErrorMessages(errs validator.ValidationErrors) []string {
	var msgs []string
	for _, err := range errs {
		var msg string
		switch err.Tag {
		case "required":
			msg = err.Field + ": Is Required."
		default:
			msg = err.Field + ": Validation failed."
		}
		msgs = append(msgs, msg)
	}
	return msgs
}

// Pad a string to 2 characters on right
func RightPad2Len(s string, padStr string, overallLen int) string {
	var padCountInt int
	padCountInt = 1 + ((overallLen - len(padStr)) / len(padStr))
	var retStr = s + strings.Repeat(padStr, padCountInt)
	return retStr[:overallLen]
}

// Reverse a string
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// Remove unwanted characters from Image name.
func SanitizeImageName(input string) string {
	illegalChars := regexp.MustCompile(`[^[a-zA-Z0-9\s]]*`)
	output := illegalChars.ReplaceAllString(input, "-")
	return output
}

func AppendIfMissingInt(slice []int, i int) []int {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
}

func Round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

func GetConfig() *appconfig.AppConfig {
	return config.ApplicationConfig.(*appconfig.AppConfig)
}

func PrintInfo(s string) {
	fmt.Println(s)
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
		tags := []string{"product-error", "panic"}
		notification.SendNotification(title, text, tags, datadog.ERROR)
	}
}

func GetTyByName(name string) (int, error) {
	sql := `SELECT
  		id_catalog_ty
			FROM catalog_ty
			WHERE name = '` + name + `'`
	driver, err := factory.GetMySqlDriver(PRODUCT_COLLECTION)
	if err != nil {
		logger.Error("Cannot initiate mysql" + err.Error())
		return 0, err
	}
	result, sqlerr := driver.Query(sql)
	if sqlerr != nil {
		logger.Error("Cannot initiate mysql" + sqlerr.DeveloperMessage)
		return 0, errors.New(sqlerr.DeveloperMessage)
	}
	var idTy int
	defer result.Close()
	for result.Next() {
		result.Scan(&idTy)
	}
	return idTy, nil
}

func GetTyByCategory(categories []int) ([]int, error) {
	catStr := ""
	for _, val := range categories {
		catStr = fmt.Sprintf("%s%d,", catStr, val)
	}
	catStr = catStr[0 : len(catStr)-1]
	sql := `SELECT
  		fk_catalog_ty
			FROM catalog_category_ty
			WHERE fk_catalog_category in (` + catStr + `)`
	driver, err := factory.GetMySqlDriver(PRODUCT_COLLECTION)
	if err != nil {
		logger.Error("Cannot initiate mysql" + err.Error())
		return nil, err
	}
	result, sqlerr := driver.Query(sql)
	if sqlerr != nil {
		logger.Error("Cannot initiate mysql" + sqlerr.DeveloperMessage)
		return nil, errors.New(sqlerr.DeveloperMessage)
	}
	var (
		fkTy  int
		TyArr []int
	)
	defer result.Close()
	for result.Next() {
		result.Scan(&fkTy)
		TyArr = append(TyArr, fkTy)
	}
	return TyArr, nil
}

func SetSpecialPriceDates(s *ProductSimple) (*time.Time, *time.Time) {
	tmpSpFrom := s.SpecialFromDate
	tmpSpTo := s.SpecialToDate
	defFrom, _ := time.Parse(time.RFC3339, DEFAULT_FROM_DATE)
	defTo, _ := time.Parse(time.RFC3339, DEFAULT_TO_DATE)

	if tmpSpFrom == nil && tmpSpTo == nil {
		return tmpSpFrom, tmpSpTo
	}
	if tmpSpFrom == nil {
		tmpSpFrom = &defFrom
	}
	if tmpSpTo == nil {
		tmpSpTo = &defTo
	}
	return tmpSpFrom, tmpSpTo
}

func FetchCategoryTree(leafCategory int) []int {
	var result []int
	var currentId int = leafCategory
	for currentId > 0 {
		parentId := func(catId int) int {
			catInfo := categoryService.ById(catId)
			return catInfo.Parent
		}(currentId)
		result = append(result, currentId)
		currentId = parentId
	}
	return result
}

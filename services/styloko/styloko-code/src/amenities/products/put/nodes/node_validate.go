package nodes

import (
	proUtil "amenities/products/common"
	put "amenities/products/put"
	"amenities/products/put/updaters"
	"common/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
)

type ValidateUpdate struct {
	id string
}

func (vu *ValidateUpdate) SetID(id string) {
	vu.id = id
}

func (vu ValidateUpdate) GetID() (string, error) {
	return vu.id, nil
}

func (vu ValidateUpdate) Name() string {
	return "ValidateUpdateRequest"
}

func (vu ValidateUpdate) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	logger.Info("Enter Validate")

	//ADD NODE PROFILER
	profiler := logger.NewProfiler()
	logger.StartProfile(profiler, proUtil.PUT_VALIDATE_NODE)
	defer logger.EndProfile(profiler, proUtil.PUT_VALIDATE_NODE)

	//SET DEBUG MESSAGE
	io.ExecContext.SetDebugMsg(proUtil.DEBUG_KEY_NODE, "Validate")

	httpReq, err := utils.GetRequestFromIO(io)
	if err != nil {
		return io, &constants.AppError{Code: constants.ResourceErrorCode,
			Message:          "cannot read request data",
			DeveloperMessage: err.Error(),
		}
	}
	httpReqOrig := httpReq.OriginalRequest
	data, _ := ioutil.ReadAll(httpReqOrig.Body)
	logger.Warning(string(data))

	//Get and Set update header in IO data
	updateType := httpReqOrig.Header.Get(proUtil.HEADER_UPDATE_TYPE)
	io.IOData.Set(proUtil.HEADER_UPDATE_TYPE, strings.ToLower(updateType))
	//Set Debug log
	io.ExecContext.SetDebugMsg("update type", updateType)
	logger.Info(updateType)

	//Preapre IO data
	ioDataArr := []*put.ProIOData{}
	var requestData interface{}
	err = func() error {
		vu.setCustomDatadogMetrics(io, updateType)
		switch strings.ToLower(updateType) {
		case proUtil.UPDATE_TYPE_MYSQL:
			requestData = &[]*updaters.MySqlUpdate{}

		case proUtil.UPDATE_TYPE_VIDEO:
			requestData = &[]*updaters.VideoUpdate{}

		case proUtil.UPDATE_TYPE_VIDEO_STATUS:
			requestData = &[]*updaters.VideoStatusUpdate{}

		case proUtil.UPDATE_TYPE_IMAGEDEL:
			requestData = &updaters.ImageDel{}

		case proUtil.UPDATE_TYPE_IMAGEADD:
			requestData = &[]*updaters.ImageAdd{}

		case proUtil.UPDATE_TYPE_NODE:
			requestData = &[]*updaters.Node{}

		case proUtil.UPDATE_TYPE_PRODUCT:
			requestData = &[]*updaters.ProductUpdate{}

		case proUtil.UPDATE_TYPE_SHIPMENT:
			requestData = &[]*updaters.ProductShipmentUpdate{}

		case proUtil.UPDATE_TYPE_PRICE:
			requestData = &[]*proUtil.M{}

		case proUtil.UPDATE_TYPE_CACHE:
			requestData = &[]*updaters.CacheInvalidate{}

		case proUtil.UPDATE_TYPE_PRODUCT_ATTRIBUTE:
			requestData = &[]*updaters.ProductAttributeUpdate{}

		case proUtil.UPDATE_TYPE_PRODUCT_STATUS:
			requestData = &[]*updaters.ProductStatusUpdate{}

		case proUtil.UPDATE_TYPE_JABONG_DISCOUNT:
			requestData = &[]*updaters.JabongDiscount{}

		default:
			vu.setCustomDatadogMetrics(io, "")
			return errors.New("Update Type Not Supported")
		}
		err := vu.UnMarshal(data, requestData)
		if err != nil {
			return errors.New("Data Unmarshalling failed: " + err.Error())
		}
		ioDataArr, err = vu.PrepareIOData(requestData, updateType)
		if err != nil {
			return err
		}
		return nil
	}()
	if err != nil {
		logger.Error(fmt.Errorf("(vu ValidateUpdate)#Execute: %s", err.Error()))
		return io, &constants.AppError{Code: constants.ResourceErrorCode,
			Message:          "Parsing Request Data Failed",
			DeveloperMessage: err.Error(),
		}
	}
	io.IOData.Set(proUtil.IODATA, ioDataArr)
	logger.Info("Exit Validate")
	return io, nil
}

// Unmarshal Request data
func (vu ValidateUpdate) UnMarshal(data []byte, requestData interface{}) error {
	err := json.Unmarshal(data, requestData)
	if err != nil {
		return fmt.Errorf("(vu ValidateUpdate)#UnMarshal failed: %s", err.Error())
	}
	return nil
}

// PrepareIOData prepares IO Data for usage in upcoming nodes
func (vu ValidateUpdate) PrepareIOData(
	requestData interface{}, reqType string) ([]*put.ProIOData, error) {

	ioDataArr := []*put.ProIOData{}
	reqTypeLower := strings.ToLower(reqType)
	var ok bool

	switch reqTypeLower {
	case proUtil.UPDATE_TYPE_MYSQL:
		var reqData *[]*updaters.MySqlUpdate
		if reqData, ok = requestData.(*[]*updaters.MySqlUpdate); !ok {
			return ioDataArr, errors.New(
				"(vu ValidateUpdate)#PrepareIOData UPDATE_TYPE_MYSQL: Assertion failed")
		}
		for i, v := range *reqData {
			iodata := vu.validate(v, strconv.Itoa(i))
			ioDataArr = append(ioDataArr, iodata)
		}
	case proUtil.UPDATE_TYPE_VIDEO:
		var reqData *[]*updaters.VideoUpdate
		if reqData, ok = requestData.(*[]*updaters.VideoUpdate); !ok {
			return ioDataArr, errors.New(
				"(vu ValidateUpdate)#PrepareIOData UPDATE_TYPE_VIDEO: Assertion failed")
		}
		for i, v := range *reqData {
			iodata := vu.validate(v, strconv.Itoa(i))
			ioDataArr = append(ioDataArr, iodata)
		}
	case proUtil.UPDATE_TYPE_VIDEO_STATUS:
		var reqData *[]*updaters.VideoStatusUpdate
		if reqData, ok = requestData.(*[]*updaters.VideoStatusUpdate); !ok {
			return ioDataArr, errors.New(
				"(vu ValidateUpdate)#PrepareIOData UPDATE_TYPE_VIDEO_STATUS: Assertion failed")
		}
		for i, v := range *reqData {
			iodata := vu.validate(v, strconv.Itoa(i))
			ioDataArr = append(ioDataArr, iodata)
		}

	case proUtil.UPDATE_TYPE_IMAGEDEL:
		var reqData *updaters.ImageDel
		if reqData, ok = requestData.(*updaters.ImageDel); !ok {
			return ioDataArr, errors.New(
				"(vu ValidateUpdate)#PrepareIOData UPDATE_TYPE_IMAGEDEL: Assertion failed")
		}
		iodata := vu.validate(reqData, "Image del")
		ioDataArr = append(ioDataArr, iodata)

	case proUtil.UPDATE_TYPE_IMAGEADD:
		var reqData *[]*updaters.ImageAdd
		if reqData, ok = requestData.(*[]*updaters.ImageAdd); !ok {
			return ioDataArr, errors.New(
				"(vu ValidateUpdate)#PrepareIOData UPDATE_TYPE_IMAGEADD: Assertion failed")
		}
		for i, v := range *reqData {
			iodata := vu.validate(v, strconv.Itoa(i))
			ioDataArr = append(ioDataArr, iodata)
		}

	case proUtil.UPDATE_TYPE_PRODUCT:
		var reqData *[]*updaters.ProductUpdate
		if reqData, ok = requestData.(*[]*updaters.ProductUpdate); !ok {
			return ioDataArr, errors.New(
				"(vu ValidateUpdate)#PrepareIOData UPDATE_TYPE_PRODUCT: Assertion failed")
		}
		for i, v := range *reqData {
			iodata := vu.validate(v, strconv.Itoa(i))
			ioDataArr = append(ioDataArr, iodata)
		}

	case proUtil.UPDATE_TYPE_SHIPMENT:
		var reqData *[]*updaters.ProductShipmentUpdate
		if reqData, ok = requestData.(*[]*updaters.ProductShipmentUpdate); !ok {
			return ioDataArr, errors.New(
				"(vu ValidateUpdate)#PrepareIOData UPDATE_TYPE_SHIPMENT: Assertion failed")
		}
		for i, v := range *reqData {
			iodata := vu.validate(v, strconv.Itoa(i))
			ioDataArr = append(ioDataArr, iodata)
		}

	case proUtil.UPDATE_TYPE_CACHE:
		var reqData *[]*updaters.CacheInvalidate
		if reqData, ok = requestData.(*[]*updaters.CacheInvalidate); !ok {
			return ioDataArr, errors.New(
				"(vu ValidateUpdate)#PrepareIOData UPDATE_TYPE_CACHE: Assertion failed")
		}
		for i, v := range *reqData {
			iodata := vu.validate(v, strconv.Itoa(i))
			ioDataArr = append(ioDataArr, iodata)
		}

	case proUtil.UPDATE_TYPE_NODE:
		var reqData *[]*updaters.Node
		if reqData, ok = requestData.(*[]*updaters.Node); !ok {
			return ioDataArr, errors.New(
				"(vu ValidateUpdate)#PrepareIOData UPDATE_TYPE_NODE: Assertion failed")
		}
		for i, v := range *reqData {
			iodata := vu.validate(v, strconv.Itoa(i))
			ioDataArr = append(ioDataArr, iodata)
		}

	case proUtil.UPDATE_TYPE_PRICE:
		reqData := requestData.(*[]*proUtil.M)
		for i, v := range *reqData {
			pu, err := vu.MapReqData2PriceUpdate(*v)
			if err != nil {
				logger.Error(fmt.Sprintf("(vu ValidateUpdate)#PrepareIOData UPDATE_TYPE_PRICE: %s", err.Error()))
				iodata := &put.ProIOData{}
				iodata.SetFailure(constants.AppError{
					Code:             constants.IncorrectDataErrorCode,
					Message:          err.Error(),
					DeveloperMessage: err.Error(),
				})
				ioDataArr = append(ioDataArr, iodata)
				continue
			}
			iodata := vu.validate(pu, strconv.Itoa(i))
			ioDataArr = append(ioDataArr, iodata)
		}

	case proUtil.UPDATE_TYPE_PRODUCT_ATTRIBUTE:
		var reqData *[]*updaters.ProductAttributeUpdate
		if reqData, ok = requestData.(*[]*updaters.ProductAttributeUpdate); !ok {
			return ioDataArr, errors.New(
				"(vu ValidateUpdate)#PrepareIOData UPDATE_TYPE_PRODUCT_ATTRIBUTE: Assertion failed")
		}
		for i, v := range *reqData {
			iodata := vu.validate(v, strconv.Itoa(i))
			ioDataArr = append(ioDataArr, iodata)
		}

	case proUtil.UPDATE_TYPE_PRODUCT_STATUS:
		var reqData *[]*updaters.ProductStatusUpdate
		if reqData, ok = requestData.(*[]*updaters.ProductStatusUpdate); !ok {
			return ioDataArr, errors.New(
				"(vu ValidateUpdate)#PrepareIOData UPDATE_TYPE_PRODUCT_STATUS: Assertion failed")
		}
		for i, v := range *reqData {
			iodata := vu.validate(v, strconv.Itoa(i))
			ioDataArr = append(ioDataArr, iodata)
		}

	case proUtil.UPDATE_TYPE_JABONG_DISCOUNT:
		var reqData *[]*updaters.JabongDiscount
		if reqData, ok = requestData.(*[]*updaters.JabongDiscount); !ok {
			return ioDataArr, errors.New(
				"(vu ValidateUpdate)#PrepareIOData UPDATE_TYPE_JABONG_DISCOUNT: Assertion failed")
		}
		for i, v := range *reqData {
			iodata := vu.validate(v, strconv.Itoa(i))
			ioDataArr = append(ioDataArr, iodata)
		}
	}
	return ioDataArr, nil
}

//
// validate struct implementing ProductUpdater
//
func (vu ValidateUpdate) validate(v put.ProductUpdater, identifier string) *put.ProIOData {
	iodata := &put.ProIOData{}
	errs := v.Validate()
	if errs != nil {
		appError := constants.AppError{
			Code:             constants.IncorrectDataErrorCode,
			Message:          "Validation Failed : " + identifier,
			DeveloperMessage: strings.Join(errs, ";"),
		}
		iodata.SetFailure(appError)
	}
	iodata.SetReqData(v)
	return iodata
}

// Specialized to be used in price update, converts priceUpdate request
// to a productUpdater format for validation and further processing.
func (vu ValidateUpdate) MapReqData2PriceUpdate(v proUtil.M) (*updaters.PriceUpdate, error) {
	pu := &updaters.PriceUpdate{}
	if simple, ok := v["simpleId"]; ok {
		pu.SimpleId, _ = utils.GetInt(simple)
	}

	if price, ok := v["price"]; ok {
		pu.Price, _ = utils.GetFloat(price)
	}
	if sprice, ok := v["specialPrice"]; ok {
		pu.SpecialPrice = proUtil.FloatNull{}
		pu.SpecialPrice.Isset = true
		spriceF, _ := utils.GetFloat(sprice)
		if spriceF >= 1 {
			pu.SpecialPrice.Value = &spriceF
		} else {
			pu.SpecialPrice.Value = nil
		}
	} else {
		pu.SpecialPrice.Isset = false
	}
	if sFromDate, ok := v["specialFromDate"]; ok {
		pu.SpecialFromDate = proUtil.TimeNull{}
		pu.SpecialFromDate.Isset = true
		if sFromDate == nil {
			pu.SpecialFromDate.Value = nil
		} else {
			sFromDate, _ := sFromDate.(string)
			value, err := proUtil.FromMysqlTime(sFromDate, true)
			if err != nil {
				return pu, errors.New("Cannot Parse Special From Date")
			}
			pu.SpecialFromDate.Value = value
		}
	} else {
		pu.SpecialFromDate.Isset = false
	}
	if sToDate, ok := v["specialToDate"]; ok {
		pu.SpecialToDate = proUtil.TimeNull{}
		pu.SpecialToDate.Isset = true

		if sToDate == nil {
			pu.SpecialToDate.Value = nil
		} else {
			sToDate, _ := sToDate.(string)
			value, err := proUtil.FromMysqlTime(sToDate, true)
			if err != nil {
				return pu, errors.New("Cannot Parse Special To Date")
			}
			*value = (*value).Add(time.Duration(proUtil.TO_DATE_DIFF) * time.Second)
			pu.SpecialToDate.Value = value
		}
	} else {
		pu.SpecialToDate.Isset = false
	}
	return pu, nil
}

func (vu ValidateUpdate) setCustomDatadogMetrics(io workflow.WorkFlowData, ty string) {
	if ty == "" {
		io.ExecContext.Set(constants.MONITOR_CUSTOM_METRIC_PREFIX, "")
		return
	}
	metricName := "_CUSTOM_PRODUCTS_PUT_" + ty + "_"
	io.ExecContext.Set(constants.MONITOR_CUSTOM_METRIC_PREFIX, metricName)
}

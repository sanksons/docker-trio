package put

import (
	proUtil "amenities/products/common"
	"common/appconfig"

	"github.com/jabong/floRest/src/common/constants"
	validator "gopkg.in/go-playground/validator.v8"
)

//Define Global variables
var (
	CacheMngr     proUtil.CacheManager
	Validate      *validator.Validate
	Conf          *appconfig.AppConfig
	DbAdapterName string
)

//
// IOData to be passed from node to node for PUT request
//
type ProIOData struct {
	Product *proUtil.Product
	ReqData ProductUpdater
	Error   *constants.AppError
	Status  string
}

//Set Request Data
func (data *ProIOData) SetReqData(d ProductUpdater) {
	data.ReqData = d
	return
}

//Set success info
func (data *ProIOData) SetSuccess(p *proUtil.Product) {
	data.Status = proUtil.STATUS_SUCCESS
	data.Product = p
	data.Error = nil
	return
}

//Set failure status and error
func (data *ProIOData) SetFailure(err constants.AppError) {
	data.Status = proUtil.STATUS_FAILURE
	data.Error = &err
	return
}

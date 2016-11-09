// Update-Type: MySql
package updaters

import (
	proUtil "amenities/products/common"
	put "amenities/products/put"
	"common/constants"
	taskPool "common/pool/tasker"
	"fmt"
	"github.com/jabong/floRest/src/common/utils/logger"
	validator "gopkg.in/go-playground/validator.v8"
)

type MySqlUpdate struct {
	ConfigId int `json:"configId" validate:"required"`
}

func (mu *MySqlUpdate) Validate() []string {
	errs := put.Validate.Struct(mu)
	if errs != nil {
		validationErrors := errs.(validator.ValidationErrors)
		msgs := proUtil.PrepareErrorMessages(validationErrors)
		return msgs
	}
	return nil
}

func (mu *MySqlUpdate) Update() (proUtil.Product, error) {

	p, err := proUtil.GetAdapter(put.DbAdapterName).GetById(mu.ConfigId)
	if err != nil {
		return p, fmt.Errorf("(mu *MySqlUpdate)#Update(): %s", err.Error())
	}
	//Cancel all previous jobs first.
	err = taskPool.CancelTasks(constants.PRODUCT_RESOURCE_NAME, p.SeqId)
	if err != nil {
		logger.Error(fmt.Errorf("(mu *MySqlUpdate)#Update():%s", err.Error()))
		return p, fmt.Errorf("(mu *MySqlUpdate)#Update():%s", err.Error())
	}
	//Add JOB for syncing to mysql
	//pick price from mysql master
	for index, _ := range p.Simples {
		p.Simples[index].SetPriceDetailsFrmMysql(true)
	}
	taskPool.AddProductSyncJob(p.SeqId, proUtil.UPDATE_TYPE_PRODUCT, p)
	return p, err
}

func (mu *MySqlUpdate) InvalidateCache() error {
	p, err := proUtil.GetAdapter(put.DbAdapterName).GetById(mu.ConfigId)
	if err != nil {
		return fmt.Errorf("(mu *MySqlUpdate)#InvalidateCache(): %s", err.Error())
	}
	go func() {
		defer proUtil.RecoverHandler("MySqlUpdate Invalidate Cache")
		put.CacheMngr.DeleteById([]int{p.SeqId}, true)
	}()
	return nil
}

func (mu *MySqlUpdate) Publish() error {
	p, err := proUtil.GetAdapter(put.DbAdapterName).GetById(mu.ConfigId)
	if err != nil {
		return err
	}
	go func() {
		defer proUtil.RecoverHandler("MySqlUpdate#Publish")
		p.Publish("", true)
	}()
	return nil
}

func (mu *MySqlUpdate) Response(p *proUtil.Product) interface{} {
	response := make(map[string]interface{})
	response["configId"] = p.SeqId
	return response
}

//
// Acquire Lock
//
func (mu *MySqlUpdate) Lock() bool {
	return true
}

//
// Release Lock
//
func (mu *MySqlUpdate) UnLock() bool {
	return true
}

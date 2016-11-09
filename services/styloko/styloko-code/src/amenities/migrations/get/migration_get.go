package get

import (
	"amenities/migrations/common"
	"amenities/migrations/common/attribute"
	"amenities/migrations/common/brand"
	"amenities/migrations/common/category"
	"amenities/migrations/common/filter"
	"amenities/migrations/common/product"
	"amenities/migrations/common/sizechart"
	"amenities/migrations/common/util"
	proUtil "amenities/products/common"
	"common/ResourceFactory"
	"common/appconstant"
	"common/utils"
	"fmt"
	"strconv"
	"strings"

	florest_constants "github.com/jabong/floRest/src/common/constants"
	workflow "github.com/jabong/floRest/src/common/orchestrator"
	"github.com/jabong/floRest/src/common/utils/logger"
)

// MigrationGet -> struct for node based data
type MigrationGet struct {
	id string
}

// SetID -> set ID for current node from orchestrator
func (mg *MigrationGet) SetID(id string) {
	mg.id = id
}

// GetID -> returns current node ID to orchestrator
func (mg MigrationGet) GetID() (id string, err error) {
	return mg.id, nil
}

// Name -> Returns node name to orchestrator
func (mg MigrationGet) Name() string {
	return "MigrationGet"
}

// Execute -> Executes the current workflow
func (mg MigrationGet) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	defer common.RecoverHandler("MigrationAPI")
	qParams, ok := utils.GetQueryParams(io, "key")
	if !ok || len(qParams) < 1 {
		return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Missing key", DeveloperMessage: "Please provide a valid key."}
	}
	switch strings.ToLower(qParams) {
	case common.AttributeSets:
		attribute.StartAttributeSetMigration()
		err := mg.setSequence(util.AttributeSets)
		if err != nil && err.Error() != "" {
			return io, mg.genError(util.AttributeSets, err)
		}
		io.IOData.Set(florest_constants.RESULT, "AttributeSet migration finished")
		return io, nil

	case common.Attributes:
		attribute.StartAttributeMigration()
		err := mg.setSequence(util.Attributes)
		if err != nil && err.Error() != "" {
			return io, mg.genError(util.Attributes, err)
		}
		attribute.StartMapping()
		io.IOData.Set(florest_constants.RESULT, "Attribute migration finished")
		return io, nil

	case common.AttributesById:
		q, ok := utils.GetQueryParams(io, "id")
		if !ok || len(q) < 1 {
			return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Missing id", DeveloperMessage: "Please provide a valid product ID"}
		}
		pid, err := strconv.Atoi(q)
		if err != nil {
			return io, mg.genError(util.Attributes, err)
		}
		attribute.MigrateSingleAttribute(pid)
		err = mg.setSequence(util.Attributes)
		if err != nil && err.Error() != "" {
			return io, mg.genError(util.Attributes, err)
		}
		io.IOData.Set(florest_constants.RESULT, "Attributes by ID migration finished")
		return io, nil

	case common.Brands:
		brand.StartBrandMigration()
		err := mg.setSequence(util.Brands)
		if err != nil && err.Error() != "" {
			return io, mg.genError(util.Brands, err)
		}
		io.IOData.Set(florest_constants.RESULT, "Brands migration finished")
		return io, nil

	case common.Filters:
		filter.StartFilterMigration()
		err := mg.setSequence(util.Filters)
		if err != nil && err.Error() != "" {
			return io, mg.genError(util.Filters, err)
		}
		io.IOData.Set(florest_constants.RESULT, "Filters migration finished")
		return io, nil

	case common.Categories:
		category.StartCategoryMigration()
		err := mg.setSequence(util.Categories)
		if err != nil && err.Error() != "" {
			return io, mg.genError(util.Categories, err)
		}
		io.IOData.Set(florest_constants.RESULT, "Category migration finished")
		return io, nil

	case common.SizeCharts:
		err := sizechart.StartSizeChartMigrationPartial()
		if err != nil {
			return io, mg.genError("sizecharts", err)
		}
		err = mg.setSequence(util.SizeCharts)
		if err != nil && err.Error() != "" {
			return io, mg.genError(util.SizeCharts, err)
		}
		io.IOData.Set(florest_constants.RESULT, "Sizechart migration finished")
		return io, nil

	case common.TaxClass:
		product.MigrateTaxClass()
		err := mg.setSequence(util.TaxClass)
		if err != nil && err.Error() != "" {
			return io, mg.genError(util.TaxClass, err)
		}
		io.IOData.Set(florest_constants.RESULT, "Tax Class migration finished")
		return io, nil

	// Product and related migration cases below
	case common.Products:
		flag, err := common.GetFlagFromRedis(common.Products)
		if err != nil {
			return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Redis connection cannot be made", DeveloperMessage: "Redis connection cannot be made to get flag"}
		}
		if flag {
			return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Already running.", DeveloperMessage: "Another migration with the same key is already running."}
		}
		err = common.SetRedisFlag(common.Products, "true")
		if err != nil {
			return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Redis connection cannot be made", DeveloperMessage: "Redis connection cannot be made to set flag"}
		}
		go func() {
			defer common.RecoverHandler(common.Products)
			product.StartActiveMigration(false)
			product.StartInActiveMigration()
			fmt.Println("Done")
			common.DeleteRedisFlag(common.Products)
			err := mg.setSequence(util.Products)
			if err != nil && err.Error() != "" {
				logger.Error(err.Error())
			}
			mg.setSequence(util.Simples)
			mg.setSequence(util.ProductImages)
			mg.setSequence(util.ProductVideos)
		}()
		io.IOData.Set(florest_constants.RESULT, "Product migration started in background")
		return io, nil

	case common.ProductsActive:
		flag, err := common.GetFlagFromRedis(common.ProductsActive)
		if err != nil {
			return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Redis connection cannot be made", DeveloperMessage: "Redis connection cannot be made to get flag"}
		}
		if flag {
			return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Already running.", DeveloperMessage: "Another migration with the same key is already running."}
		}
		err = common.SetRedisFlag(common.ProductsActive, "true")
		if err != nil {
			return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Redis connection cannot be made", DeveloperMessage: "Redis connection cannot be made to set flag"}
		}
		go func() {
			defer common.RecoverHandler(common.ProductsActive)
			product.StartActiveMigration(false)
			common.DeleteRedisFlag(common.ProductsActive)
		}()
		io.IOData.Set(florest_constants.RESULT, "Active product migration started in background")
		return io, nil

	case common.ProductsDrop:
		flag, err := common.GetFlagFromRedis(common.ProductsDrop)
		if err != nil {
			return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Redis connection cannot be made", DeveloperMessage: "Redis connection cannot be made to get flag"}
		}
		if flag {
			return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Already running.", DeveloperMessage: "Another migration with the same key is already running."}
		}
		err = common.SetRedisFlag(common.ProductsDrop, "true")
		if err != nil {
			return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Redis connection cannot be made", DeveloperMessage: "Redis connection cannot be made to set flag"}
		}
		go func() {
			defer common.RecoverHandler(common.ProductsDrop)
			product.StartActiveMigration(true)
			product.StartInActiveMigration()
			err := mg.setSequence(util.Products)
			if err != nil && err.Error() != "" {
				logger.Error(err.Error())
			}
			mg.setSequence(util.Simples)
			mg.setSequence(util.ProductImages)
			mg.setSequence(util.ProductVideos)
			common.DeleteRedisFlag(common.ProductsDrop)
		}()
		io.IOData.Set(florest_constants.RESULT, "Drop product migration started in backgound")
		return io, nil

	case common.ProductsInactive:
		flag, err := common.GetFlagFromRedis(common.ProductsInactive)
		if err != nil {
			return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Redis connection cannot be made", DeveloperMessage: "Redis connection cannot be made to get flag"}
		}
		if flag {
			return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Already running.", DeveloperMessage: "Another migration with the same key is already running."}
		}
		err = common.SetRedisFlag(common.ProductsInactive, "true")
		if err != nil {
			return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Redis connection cannot be made", DeveloperMessage: "Redis connection cannot be made to set flag"}
		}
		go func() {
			defer common.RecoverHandler(common.ProductsInactive)
			product.StartInActiveMigration()
			common.DeleteRedisFlag(common.ProductsInactive)
		}()
		io.IOData.Set(florest_constants.RESULT, "Inactive product migration started in background")
		return io, nil

	case common.ProductGroups:
		flag, err := common.GetFlagFromRedis(common.ProductGroups)
		if err != nil {
			return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Redis connection cannot be made", DeveloperMessage: "Redis connection cannot be made to get flag"}
		}
		if flag {
			return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Already running.", DeveloperMessage: "Another migration with the same key is already running."}
		}
		err = common.SetRedisFlag(common.ProductGroups, "true")
		if err != nil {
			return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Redis connection cannot be made", DeveloperMessage: "Redis connection cannot be made to set flag"}
		}
		go func() {
			defer common.RecoverHandler(common.ProductGroups)
			product.MigrateProductGroup()
			err := mg.setSequence(util.ProductGroups)
			if err != nil && err.Error() != "" {
				logger.Error(err.Error)
			}
			common.DeleteRedisFlag(common.ProductGroups)
		}()
		io.IOData.Set(florest_constants.RESULT, "Product groups migration started in background.")
		return io, nil

	case common.ProductsIndex:
		flag, err := common.GetFlagFromRedis(common.ProductsIndex)
		if err != nil {
			return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Redis connection cannot be made", DeveloperMessage: "Redis connection cannot be made to get flag"}
		}
		if flag {
			return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Already running.", DeveloperMessage: "Another migration with the same key is already running."}
		}
		err = common.SetRedisFlag(common.ProductsIndex, "true")
		if err != nil {
			return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Redis connection cannot be made", DeveloperMessage: "Redis connection cannot be made to set flag"}
		}
		go func() {
			defer common.RecoverHandler(common.ProductsIndex)
			product.ReCreateIndexes()
			common.DeleteRedisFlag(common.ProductsIndex)
		}()
		io.IOData.Set(florest_constants.RESULT, "Product Re-Indexing started in background")
		return io, nil

	case common.ProductsById:
		q, ok := utils.GetQueryParams(io, "id")
		if !ok || len(q) < 1 {
			return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Missing id", DeveloperMessage: "Please provide a valid product ID"}
		}
		pid, _ := strconv.Atoi(q)
		product.MigrateSingleProduct(pid)
		io.IOData.Set(florest_constants.RESULT, "Product by ID migration finished")
		return io, nil

	case common.ProductsBySeller:
		q, ok := utils.GetQueryParams(io, "id")
		if !ok || len(q) < 1 {
			return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Missing id", DeveloperMessage: "Please provide a seller ID"}
		}
		sellerid, _ := strconv.Atoi(q)
		product.StartSellerMigration(sellerid)
		io.IOData.Set(florest_constants.RESULT, "Product by Seller ID migration finished.")
		return io, nil

	case common.ProductsByBrand:
		q, ok := utils.GetQueryParams(io, "id")
		if !ok || len(q) < 1 {
			return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Missing id", DeveloperMessage: "Please provide a Brand ID"}
		}
		brandId, _ := strconv.Atoi(q)
		product.StartMigrationByBrand(brandId)
		io.IOData.Set(florest_constants.RESULT, "Product by Brand ID migration finished.")
		return io, nil

	case common.ProductsByPromotion:
		q, ok := utils.GetQueryParams(io, "id")
		if !ok || len(q) < 1 {
			return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Missing id", DeveloperMessage: "Please provide a Promotion ID"}
		}
		promoId, _ := strconv.Atoi(q)
		product.StartMigrationByPromotion(promoId)
		io.IOData.Set(florest_constants.RESULT, "Product by Promotion ID migration finished.")
		return io, nil

	case common.ProductsSizeChartById:
		q, ok := utils.GetQueryParams(io, "id")
		if !ok || len(q) < 1 {
			return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Missing id", DeveloperMessage: "Please provide a Product ID"}
		}
		prodId, _ := strconv.Atoi(q)
		err := sizechart.SingleProductSizeChartUpdate(prodId)
		if err != nil {
			return io, &florest_constants.AppError{
				Code:             appconstant.FailedToCreateErrorCode,
				Message:          "Could not update product with sizechart",
				DeveloperMessage: err.Error(),
			}
		}
		io.IOData.Set(florest_constants.RESULT, "Sizechart By Product ID migration finished.")
		return io, nil

	case common.SizeChartMapping:
		err := sizechart.StartSKUSizechartMapping()
		if err != nil && err.Error() != "" {
			return io, mg.genError(common.SizeChartMapping, err)
		}
		io.IOData.Set(florest_constants.RESULT, "Sizechart Mapping migration finished")
		return io, nil

	case common.ProductsSizeChart:
		flag, err := common.GetFlagFromRedis(common.ProductsSizeChart)
		if err != nil {
			return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Redis connection cannot be made", DeveloperMessage: "Redis connection cannot be made to get flag"}
		}
		if flag {
			return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Already running.", DeveloperMessage: "Another migration with the same key is already running."}
		}
		err = common.SetRedisFlag(common.ProductsSizeChart, "true")
		if err != nil {
			return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Redis connection cannot be made", DeveloperMessage: "Redis connection cannot be made to set flag"}
		}
		go func() {
			defer common.RecoverHandler(common.ProductsSizeChart)
			err := sizechart.WriteSizeChartToProduct()
			if err != nil {
				logger.Error(err.Error())
			}
			common.DeleteRedisFlag(common.ProductsSizeChart)
		}()
		io.IOData.Set(florest_constants.RESULT, "Product Sizechart migration finished")
		return io, nil
	case common.ResetCounter:
		err := mg.setSequence(
			proUtil.DUMMY_IMAGES,
			util.Products,
			util.Simples,
			util.ProductImages,
			util.ProductVideos,
			util.AttributeSets,
			util.Attributes,
			util.Brands,
			util.Categories,
			util.Filters,
			util.ProductGroups,
			util.TaxClass,
			util.SizeCharts,
			util.PrePack,
		)
		if err != nil {
			logger.Error(fmt.Sprintf("Reset counters failed: %s", err.Error()))
			return io, &florest_constants.AppError{
				Code:             appconstant.BadRequestCode,
				Message:          "Reset Counters Failed.",
				DeveloperMessage: err.Error(),
			}
		}
		io.IOData.Set(florest_constants.RESULT, "Reset counters finished")
		return io, nil
	default:
		return io, &florest_constants.AppError{Code: appconstant.BadRequestCode, Message: "Invalid key", DeveloperMessage: "Valid keys are: " + common.GetKeys()}
	}
}

func (mg MigrationGet) setSequence(args ...string) error {
	mgoSession := ResourceFactory.GetMongoSession("Sequencing")
	var err error
	defer mgoSession.Close()
	for _, v := range args {
		var counter int
		switch v {
		case
			proUtil.DUMMY_IMAGES,
			util.Brands,
			util.SizeCharts,
			util.ProductGroups,
			util.Products,
			util.Simples,
			util.ProductImages,
			util.ProductVideos:
			counter, err = common.GetSeqCounter(v)
			if err != nil {
				return err
			}
		case
			util.PrePack:
			counter, err = common.GetSeqCounterMysql(v)
			if err != nil {
				fmt.Println(err)
				return err
			}
		}
		err = mgoSession.SetCollectionInCounter(v, counter)
		if err != nil && err.Error() != "" {
			return err
		}
	}
	return nil
}

func (mg MigrationGet) genError(name string, err error) *florest_constants.AppError {
	return &florest_constants.AppError{Code: appconstant.ServiceFailureCode, Message: name + " Migration Failed", DeveloperMessage: err.Error()}
}

package common

import ()

const ColumnsOnMouseHover = "2,3"
const SizeChartCacheKey = "sizechart"

const (
	// Worker Pool
	SizeChartCreate     = "SizeChartCreate"
	SizeChartProdCreate = "SizeChartProdCreate"
	PoolSize            = 1
	QueueSize           = 4
	RetryCount          = 2
	WaitTime            = 500
	SizeChartWorker     = "CreateSizeChartWorker"
)

// Mongo Constant
const (
	// MongoSession Name
	SizeChartAPI           = "SizeChartApi"
	SizeChartGetAPI        = "SizeChartGetAPI"
	SizeChartSer           = "SizeChartService"
	SizeChartCollec        = "sizecharts"
	BrandCollection        = "brands"
	CategoryCollection     = "categories"
	ProductsCollection     = "products"
	SizeChartMappingCollec = "skuSizechartMapping"
)

// Orchestrator IO data constants
const (
	SizeChartBrandWise = "SizeChartBrandWise"
	SizeChartInput     = "SizeChartInput"
	FailedSkus         = "FailedSkus"
	SavedSizeChCollec  = "SavedSizeChCollection"
	SuccessSkus        = "successSku"
	SizeChartResource  = "SIZECHART"
)

// API error Messages
const (
	CreationFailed   = "Data insertion failed"
	FailedValidation = "Data validation failed"
	InvalidJson      = "Invalid input Json"
	SkuNotUpdated    = "SizeChart not updated to some skus"
	NoSkusFound      = "Skus doesn't exist for sizechart"
)

// Request
const (
	SizeChartLevelSku      = 0
	SizeChartLevelSpecific = 1
	SizeChartLevelGeneric  = 2
	SizeChartHeader        = "SizeChart-Type"
	BrandBrick             = "brand-brick"
	SKU                    = "sku"
	CacheControlHeader     = "No-Cache"
)

// Profiler Constants
const (
	SizeChartCreation      = "SizeChartCreation"
	GetSizeChartByProd     = "GetSizeChartByProd"
	SizeChartMapping       = "Mapping sizechart to product"
	SizeChartCreationMongo = "Creating New Size Chart in MongoDb"
	GetCategoryNameById    = "Getting category Name by Id from Mongodb"
	GetBrandNameById       = "Getiing Brand Name by Id from MongoDb"
	CheckBrandInSystem     = "Check If brand exists in the system"
	UpdateDefaultSetting   = "Update default setting for sizechart"
)

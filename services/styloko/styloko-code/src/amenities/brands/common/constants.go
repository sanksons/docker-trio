package common

const (
	//Constants pertaining to Brand Get API
	BRAND_GET       = "_BRAND_GET_ONE_"
	BRANDS_GET_ALL  = "_BRAND_GET_ALL_"
	BRANDS          = "brand"
	BRAND_API       = "BRANDS"
	BRAND_DATA      = "brandData"
	BRAND_CREATE    = "brandCreate"
	BRAND_SEARCH    = "brandSearch"
	BRAND_OPERATION = "brandOperation"

	//Brand Status constants
	ACTIVE   = "active"
	INACTIVE = "inactive"
	DELETED  = "deleted"
	ALL      = "all"
	STATUS   = "status"

	//Constants pertaining to Brand Post API
	BRAND_INSERT        = "insert"
	FAILURE             = "failure"
	FAILURE_DATA        = "failureData"
	VALIDATE_BRAND_POST = "validateBrandPost"

	//Constants pertaining to Brand Put API
	BRAND_UPDATE_DATA = "brandUpdateData"
	BRAND_UPDATE      = "brandUpdate"
	FAILURE_FLAG      = "failureFlag"
	VALIDATE_PUT      = "validatePut"

	BRAND_UPDATE_POOL_SIZE   = 2
	BRAND_UPDATE_QUEUE_SIZE  = 40
	BRAND_UPDATE_RETRY_COUNT = 3
	BRAND_UPDATE_WAIT_TIME   = 100
	BRAND_CREATE_POOL_SIZE   = 2
	BRAND_CREATE_QUEUE_SIZE  = 40
	BRAND_CREATE_RETRY_COUNT = 3
	BRAND_CREATE_WAIT_TIME   = 100

	MONGO_GET_BRAND_SEARCH_COUNT = "MONGO_GET_BRAND_SEARCH_COUNT"
	MONGO_SEARCH                 = "MONGO_SEARCH"

	CUSTOM_BRAND_GET_ALL = "_CUSTOM_BRAND_GET_ALL_"
	CUSTOM_BRAND_GET_ONE = "_CUSTOM_BRAND_GET_ONE_"
)

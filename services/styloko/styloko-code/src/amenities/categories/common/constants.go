package common

const (

	// Error strings
	NO_DATA              = "No Data"
	INVALID_DATA         = "Data cannot be unmarshalled. Invalid"
	CANNOT_BE_ZERO       = "Value cannot be 0"
	CANNOT_BE_EMPTY      = "Value cannot be empty"
	INVALID_PATH         = "Invalid Path"
	DATA_VALIDATION_FAIL = "Data validation failure."

	// Profiler constants
	GET_ALL             = "_CATEGORY_GET_ALL_NODE_"
	GET_ONE             = "_CATEGORY_GET_ONE_NODE_"
	GET_QUERY           = "_CATEGORY_GET_QUERY_NODE_"
	GET_PATH_DECISION   = "_CATEGORY_GET_PATH_DECISION_"
	GET_QUERY_DECISION  = "_CATEGORY_GET_QUERY_DECISION_"
	PUT                 = "_CATEGORY_PUT_NODE_"
	PUT_VALIDATION_NODE = "_CATEGORY_PUT_VALID_NODE_"
	PUT_CHECK_NODE      = "_CATEGORY_PUT_CHECK_NODE_"
	POST                = "_CATEGORY_POST_CREATE_NODE_"

	// Error or Data.
	CATEGORY_ERROR        = ""
	CATEGORY_VALID_DATA   = ""
	CATEGORY_PATH_PARAMS  = ""
	CATEGORY_QUERY_PARAMS = ""
	CATEGORY_ID           = ""

	// category status
	ACTIVE   = "active"
	INACTIVE = "inactive"
	DELETED  = "deleted"
	ALL      = "all"

	// Worker Constants per API
	CATEGORY_UPDATE_POOL_SIZE   = 5
	CATEGORY_UPDATE_QUEUE_SIZE  = 1000
	CATEGORY_UPDATE_RETRY_COUNT = 3
	CATEGORY_UPDATE_WAIT_TIME   = 500

	CATEGORY_CREATE_POOL_SIZE   = 5
	CATEGORY_CREATE_QUEUE_SIZE  = 1000
	CATEGORY_CREATE_RETRY_COUNT = 3
	CATEGORY_CREATE_WAIT_TIME   = 500

	// Datadog Metric Names
	CUSTOM_CATEGORY_GET_ALL = "_CUSTOM_CATEGORY_GET_ALL_"
	CUSTOM_CATEGORY_GET_ONE = "_CUSTOM_CATEGORY_GET_ONE_"

	//key for blitz
	CATEGORIES = "category"
)

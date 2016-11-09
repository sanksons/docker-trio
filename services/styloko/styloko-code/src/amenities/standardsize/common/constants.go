package common

const (

	// Error strings
	NO_DATA              = "No Data"
	INVALID_DATA         = "Invalid Data, cannot be unmarshalled"
	CANNOT_BE_ZERO       = "Value cannot be 0"
	CANNOT_BE_EMPTY      = "Value cannot be empty"
	INVALID_PATH         = "Invalid Path"
	DATA_VALIDATION_FAIL = "Data validation failure"

	// Validation Error strings
	ERROR_GENERIC         = "Unsuccessful DB insertion"
	ERROR_ATTRIBUTE_SET   = "Invalid Attribute Set"
	ERROR_LEAF_CATEGORY   = "Invalid Leaf Category Id"
	ERROR_BRAND           = "Invalid Brand"
	ERROR_INCOMPLETE_DATA = "Attribute Set, Leaf Category, Standard Size are mandatory"
	VALIDATION_FAILURE_1  = "This standard size already present for this Attribute Set and Leaf Category"
	VALIDATION_FAILURE_2  = "Attribute Set and Leaf Category mapping must be present before adding brand size mapping"
	VALIDATION_FAILURE_3  = "This standard size must have a mapping with this Attribute Set and Leaf Category before adding brand size mapping"
	VALIDATION_FAILURE_4  = "This standard size already present for this data combination"
	VALIDATION_FAILURE_5  = "Another standard size already present for this data combination"

	// Pagination Params
	DEFAULT_LIMIT   = 0
	DEFAULT_PAGE_NO = 1

	// Workflow data
	STANDARDSIZE_ERROR          = ""
	STANDARDSIZE_VALID_DATA     = ""
	STANDARDSIZE_GET_ERROR      = ""
	STANDARDSIZE_GET_VALID_DATA = ""
	STANDARDSIZE_PATH_PARAMS    = ""
	STANDARDSIZE_QUERY_PARAMS   = ""

	// Mongo Params
	STANDARDSIZE_API             = "STANDARDSIZE"
	STANDARDSIZE_SEARCH          = "StandardSizeSearch"
	STANDARDSIZE_CREATE          = "StandardSizeCreate"
	STANDARDSIZE_UPDATE          = "StandardSizeUpdate"
	STANDARDSIZE_COLLECTION      = "standardsize"
	STANDARDSIZEERROR_COLLECTION = "standardsizeerror"
	STANDARDSIZEERROR_CREATE     = "StandardSizeErrorCreate"
	STANDARDSIZEERROR_SEARCH     = "StandardSizeErrorSearch"

	// Profiling Params
	POST_CREATE        = "STANDARDSIZE_POST_CREATE"
	POST_ERROR         = "STANDARDSIZE_POST_ERROR"
	POST_VALIDATE      = "STANDARDSIZE_POST_VALIDATE"
	PUT_UPDATE         = "STANDARDSIZE_PUT_UPDATE"
	PUT_ERROR          = "STANDARDSIZE_PUT_ERROR"
	GET_ERROR          = "STANDARDSIZE_GET_ERROR"
	PUT_VALIDATE       = "STANDARDSIZE_PUT_VALIDATE"
	GET_VALIDATE       = "STANDARDSIZE_GET_VALIDATE"
	GET_PATH_DECISION  = "STANDARDSIZE_GET_PATH_DECISION"
	GET_QUERY_DECISION = "STANDARDSIZE_QUERY_PATH_DECISION"
	GET_ID             = "STANDARDSIZE_GET_ID"
	GET_SEARCH         = "STANDARDSIZE_GET_SEARCH"
	GET_ALL            = "STANDARDSIZE_GET_ALL"
	MONGO_CREATE       = "STANDARDSIZE_MONGO_CREATE"
	MONGO_UPDATE       = "STANDARDSIZE_MONGO_UPDATE"
	MONGO_GET_ID       = "STANDARDSIZE_MONGO_GET_ID"
	MONGO_GET_SEARCH   = "STANDARDSIZE_MONGO_GET_SEARCH"
	MONGO_GET_ALL      = "STANDARDSIZE_MONGO_GET_ALL"
)

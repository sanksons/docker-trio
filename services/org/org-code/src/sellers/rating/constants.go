package rating

import ()

const (
	RATING                           = "RATING"
	UPLOAD_RATING                    = "UPLOAD_RATING"
	UPLOAD_RATING_MONGO              = "UPLOAD_RATING_MONGO"
	PRODUCT_INAVLIDATION             = "PRODUCT_INAVLIDATION"
	PRODUCT_INAVLIDATION_POOL_SIZE   = 5
	PRODUCT_INAVLIDATION_QUEUE_SIZE  = 1000
	PRODUCT_INAVLIDATION_WAIT_TIME   = 100
	PRODUCT_INAVLIDATION_RETRY_COUNT = 5
)
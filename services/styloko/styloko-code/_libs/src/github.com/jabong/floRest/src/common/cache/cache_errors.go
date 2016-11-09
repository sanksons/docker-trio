package cache

import (
	"errors"
	"fmt"
)

var (
	// ErrCacheMiss means that a Get failed because the item wasn't present.
	ErrCacheMiss = errors.New("cache miss")

	// ErrUnsupportedOperation means that the operation is not supported
	ErrUnsupportedOperation = errors.New("unsupported operation")

	// ErrCASConflict means that a CompareAndSwap call failed due to the
	// cached value being modified between the Get and the CompareAndSwap.
	// If the cached value was simply evicted rather than replaced,
	// ErrNotStored will be returned instead.
	ErrCASConflict = errors.New("compare-and-swap conflict")

	// ErrNotStored means that a write operation failed in cache
	ErrNotStored = errors.New("item not stored")

	// ErrServer means that a server error occurred.
	ErrServerError = errors.New("server error")

	// ErrNoStats means that no statistics were available.
	ErrNoStats = errors.New("no statistics available")

	// ErrMalformedKey is returned when an invalid key is used.
	// Keys must be at maximum 250 bytes long, ASCII, and not
	// contain whitespace or control characters.
	ErrMalformedKey = errors.New("key is too long or contains invalid characters")

	// ErrNoServers is returned when no servers are configured or available.
	ErrNoServers = errors.New("no servers configured or available")

	// ErrBadMagic is returned when the magic number in a response is not valid.
	ErrBadMagic = errors.New("bad magic number in response")

	// ErrBadIncrDec is returned when performing a incr/decr on non-numeric values.
	ErrBadIncrDec = errors.New("incr or decr on non-numeric value")

	// ErrSerializationFailed is returned when serialization/deserialization failed
	// for a value
	ErrSerializationFailed = errors.New("serialization failed")

	// ErrCompressionFailed is returned when compression / decompression of a value
	// failed
	ErrCompressionFailed = errors.New("compression failed")

	// ErrNotSupportedFormat is returned when some not supported format is used
	// to set a cache value
	ErrNotSupportedFormat = errors.New("not supported format")

	// ErrInvalidExpireTime is returned when some invalid cache expire time is specified
	ErrInvalidExpireTime = errors.New("invalid cache expire time")

	// ErrDumpFailed is returned when dumping all keys & its values stored in a cache store
	// failed
	ErrDumpFailed = errors.New("cache dump failed")

	// ErrCacheNameSpaceNotFound not found is returned when the cache key is not found in
	// the requested namespace (namespace is like a bucket in couchbase, etc)
	ErrCacheNameSpaceNotFound = errors.New("cache namespace not found")

	// ErrCacheDeleteFailed is returned when cache deletion failed
	ErrCacheDeleteFailed = errors.New("cache namespace not found")

	// ErrInvalidCacheResponse is returned when the cache response cannot be decoded
	ErrInvalidCacheResponse = errors.New("invalid cache response")

	// ErrUnknown is returned when some unexpected error occurs
	ErrUnknown = errors.New("unknown Error")

	// ErrCacheDisabled is returned when this cache is disabled from config
	ErrCacheDisabled = errors.New("cache is disabled")
	
	ErrInvalidData = errors.New("Invalid data Supplied")
)

const (
	// ccEntityNotFoundErrCode is the error code returned when the request entity key is not
	// found in central cache
	ccEntityNotFoundErrCode int = 1000

	// ccBucketNotFoundErrCode is the error code returned when the requested bucket is not found
	// in central cache
	ccBucketNotFoundErrCode int = 1001

	// ccEntityFetchFailedErrCode is the error code when the requested entity fetch failed due to
	// some internal server error in central cache
	ccEntityFetchFailedErrCode int = 1002

	// ccEntityCreationFailedErrCode is the error code returned when setting an entity in central
	// cache failed
	ccEntityCreationFailedErrCode int = 1004

	// ccEntityInvalidTtlErrCode is the error code returned when an invalid ttl (cache expiry time)
	// is specified for a key in central cache
	ccEntityInvalidTtlErrCode int = 1006

	// ccEntityDeletionFailedErrCode is the error code returned when the deletion of a key from
	// central cache failed
	ccEntityDeletionFailedErrCode int = 1007

	// ccEntityKeyInvalidErrCode is the error code returned when an invalid key name is specified for
	// an item to store in central cache
	ccEntityKeyInvalidErrCode int = 1008

	// ccBucketIdentifierInvalidErrCode is the error code returned from central cache when an invalid bucket name is
	// specified
	ccBucketIdentifierInvalidErrCode int = 1009
)

//getCacheErrorType returns a cache error
func getCacheErrorType(cacheType string, e interface{}) error {
	switch cacheType {
	case CentralCache:
		errCode, ok := e.(int)
		if !ok {
			return errors.New("Wrong data type passed for " + cacheType)
		}
		switch errCode {
		case ccEntityNotFoundErrCode:
			return ErrCacheMiss
		case ccBucketNotFoundErrCode:
			return ErrCacheNameSpaceNotFound
		case ccEntityFetchFailedErrCode:
			return ErrServerError
		case ccEntityCreationFailedErrCode:
			return ErrNotStored
		case ccEntityInvalidTtlErrCode:
			return ErrInvalidExpireTime
		case ccEntityDeletionFailedErrCode:
			return ErrCacheDeleteFailed
		case ccEntityKeyInvalidErrCode:
			return ErrMalformedKey
		case ccBucketIdentifierInvalidErrCode:
			return ErrCacheNameSpaceNotFound
		default:
			return errors.New(fmt.Sprintf("%s Error Code %s Not handled", CentralCache, errCode))
		}
	}
	//Ideally the code should not reach here
	return errors.New(fmt.Sprintf("Unknown Cache Type %s", cacheType))
}

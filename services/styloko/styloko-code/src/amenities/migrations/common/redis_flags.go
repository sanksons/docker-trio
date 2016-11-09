package common

import (
	"common/ResourceFactory"

	"github.com/jabong/floRest/src/common/utils/logger"
)

// GetFlagFromRedis returns bool for provided key from redis hMap
func GetFlagFromRedis(key string) (bool, error) {
	redisAdapter, err := ResourceFactory.GetRedisDriver(MIGRATION_DRIVER)
	if err != nil {
		logger.Error("Cannot acquire redis driver. ", err.Error())
		return false, err
	}
	adapterResponse := redisAdapter.HGetAllMap(MIGRATION_HASH_MAP)
	flags, err := adapterResponse.Result()
	if _, ok := flags[key]; ok {
		return true, nil
	}
	return false, nil
}

// SetRedisFlag sets redis flag for solr
func SetRedisFlag(key, value string) error {
	redisAdapter, err := ResourceFactory.GetRedisDriver(MIGRATION_DRIVER)
	if err != nil {
		logger.Error("Cannot acquire redis driver. ", err.Error())
		return err
	}
	redisAdapter.HSetNX(MIGRATION_HASH_MAP, key, value)
	return nil
}

// DeleteRedisFlag deletes the key from redis hMap
func DeleteRedisFlag(key string) error {
	redisAdapter, err := ResourceFactory.GetRedisDriver(MIGRATION_DRIVER)
	if err != nil {
		logger.Error("Cannot acquire redis driver. ", err.Error())
		return err
	}
	redisAdapter.HDel(MIGRATION_HASH_MAP, key)
	return nil
}

// DeleteRedisKey deletes any key provided. Will delete full map
func DeleteRedisKey(key string) error {
	redisAdapter, err := ResourceFactory.GetRedisDriver(MIGRATION_DRIVER)
	if err != nil {
		logger.Error("Cannot acquire redis driver. ", err.Error())
		return err
	}
	redisAdapter.Del(key)
	return nil
}

// ClearAllRedisFlags deletes all redis flags for Migrations
func ClearAllRedisFlags() error {
	logger.Info("Clearing all redis flags")
	err := DeleteRedisKey(MIGRATION_HASH_MAP)
	if err != nil {
		return err
	}
	err = DeleteRedisKey(PRODUCT_BOOTSTRAP_FLAGS)
	if err != nil {
		return err
	}
	DeleteRedisKey("{styloko}_tasker_main_product")
	DeleteRedisKey("{styloko}_tasker_processing")
	return nil
}

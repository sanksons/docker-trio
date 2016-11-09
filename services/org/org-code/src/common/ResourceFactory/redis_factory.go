package ResourceFactory

import (
	"common/appconfig"
	"errors"
	"github.com/jabong/floRest/src/common/config"
	"gopkg.in/redis.v3"
	"strings"
)

const DEFAULT_REDIS_DRIVER = "DEFAULT_POOL"

var RedisMap = map[string]*redis.ClusterClient{}

func GetDefaultDriver() (*redis.ClusterClient, error) {
	return GetRedisDriver(DEFAULT_REDIS_DRIVER)
}

func GetRedisDriver(adapterName string) (*redis.ClusterClient, error) {
	if _, ok := RedisMap[adapterName]; !ok {
		pool, err := InitRedisPool()
		if err != nil {
			return pool, err
		}
		RedisMap[adapterName] = pool
	}
	return RedisMap[adapterName], nil
}

func InitRedisPool() (*redis.ClusterClient, error) {
	config, err := GetRedisConfig()
	if err != nil {
		return nil, errors.New("Unable to get Redis config")
	}
	client := redis.NewClusterClient(&config)
	_, err = client.Ping().Result()
	if err != nil {
		return nil, errors.New("cannot connect to Redis")
	}
	return client, err
}

func GetRedisConfig() (redis.ClusterOptions, error) {
	conf := config.ApplicationConfig.(*appconfig.AppConfig)
	hosts := strings.Split(conf.Redis.Host, ",")
	options := redis.ClusterOptions{}
	options.Addrs = hosts
	options.Password = ""
	options.PoolSize = conf.Redis.PoolSize
	return options, nil
}

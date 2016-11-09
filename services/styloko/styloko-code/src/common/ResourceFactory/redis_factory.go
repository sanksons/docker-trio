package ResourceFactory

import (
	"common/appconfig"
	"errors"
	"github.com/jabong/floRest/src/common/config"
	"gopkg.in/redis.v3"
	"strings"
)

const DEFAULT_REDIS_DRIVER = "DEFAULT_POOL"
const (
	REDIS_CONFIG_STOCK   = "stock"
	REDIS_CONFIG_STYLOKO = "styloko"
)

var RedisMap = map[string]map[string]*redis.ClusterClient{
	REDIS_CONFIG_STOCK:   map[string]*redis.ClusterClient{},
	REDIS_CONFIG_STYLOKO: map[string]*redis.ClusterClient{},
}

func GetDefaultDriver() (*redis.ClusterClient, error) {
	return GetRedisDriver(DEFAULT_REDIS_DRIVER)
}

func GetRedisDriver(adapterName string) (*redis.ClusterClient, error) {
	if _, ok := RedisMap[REDIS_CONFIG_STYLOKO][adapterName]; !ok {
		pool, err := InitRedisPool(REDIS_CONFIG_STYLOKO)
		if err != nil {
			return pool, err
		}
		RedisMap[REDIS_CONFIG_STYLOKO][adapterName] = pool
	}
	return RedisMap[REDIS_CONFIG_STYLOKO][adapterName], nil
}

func GetStockRedisDriver(adapterName string) (*redis.ClusterClient, error) {
	if _, ok := RedisMap[REDIS_CONFIG_STOCK][adapterName]; !ok {
		pool, err := InitRedisPool(REDIS_CONFIG_STOCK)
		if err != nil {
			return pool, err
		}
		RedisMap[REDIS_CONFIG_STOCK][adapterName] = pool
	}
	return RedisMap[REDIS_CONFIG_STOCK][adapterName], nil
}

func InitRedisPool(typ string) (*redis.ClusterClient, error) {
	config, err := GetRedisConfig(typ)
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

func GetRedisConfig(typ string) (redis.ClusterOptions, error) {
	conf := config.ApplicationConfig.(*appconfig.AppConfig)
	options := redis.ClusterOptions{}
	var hostList string
	var poolSize int
	if typ == REDIS_CONFIG_STOCK {
		hostList = conf.Redis.Stock.Host
		poolSize = conf.Redis.Stock.PoolSize
	} else {
		hostList = conf.Redis.Styloko.Host
		poolSize = conf.Redis.Styloko.PoolSize
	}
	hosts := strings.Split(hostList, ",")
	options.Addrs = hosts
	options.Password = ""
	options.PoolSize = poolSize
	return options, nil
}

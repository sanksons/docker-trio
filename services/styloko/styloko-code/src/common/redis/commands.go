package redis

import (
	factory "common/ResourceFactory"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gopkg.in/redis.v3"
)

func GetDriver() (Cluster, error) {
	driver, err := factory.GetDefaultDriver()
	if err != nil {
		return Cluster{}, err
	}
	return Cluster{Conn: driver}, nil
}

var NotFoundErr error = errors.New("Redis: Data Not Found.")

type Error struct {
	Err error
}

func (e *Error) IsNotFound() bool {
	if e.Err == NotFoundErr {
		return true
	}
	return false
}

func (e *Error) Error() string {
	return e.Err.Error()
}

type Cluster struct {
	Conn *redis.ClusterClient
}

func (client Cluster) Ping() error {
	if client.Conn == nil {
		return fmt.Errorf("Redis connection not initialized")
	}
	result := client.Conn.Ping()
	return result.Err()
}

func (client Cluster) GetSetBool(key string, val interface{}) (bool, *Error) {
	result := client.Conn.GetSet(key, val)
	r, err := result.Result()
	if err != nil && err != redis.Nil {
		//its error
		return false, &Error{Err: err}
	}
	if strings.ToLower(r) == "1" || strings.ToLower(r) == "true" {
		return true, nil
	}
	return false, nil
}

func (client Cluster) GetSetInt(key string, val interface{}) (int, *Error) {
	result := client.Conn.GetSet(key, val)
	r, err := result.Result()
	if err != nil && err != redis.Nil {
		//its error
		return 0, &Error{Err: err}
	}
	if err == redis.Nil {
		return 0, nil
	}
	intr, err := strconv.Atoi(r)
	if err != nil {
		return 0, &Error{Err: err}
	}
	return intr, nil
}

func (client Cluster) GetSetTime(key string, val interface{}) (*time.Time, *Error) {
	result := client.Conn.GetSet(key, val)
	r, err := result.Result()
	if err != nil && err != redis.Nil {
		//its error
		return nil, &Error{Err: err}
	}
	if err == redis.Nil {
		return nil, nil
	}
	intR, err := strconv.ParseInt(r, 10, 64)
	if err != nil {
		//its error
		return nil, &Error{Err: err}
	}
	t := time.Unix(intR, 0)
	return &t, nil
}

func (client Cluster) Set(key string, val interface{}) *Error {
	result := client.Conn.Set(key, val, 0)
	err := result.Err()
	if err != nil && err != redis.Nil {
		//its error
		return &Error{Err: err}
	}
	return nil
}

func (client Cluster) Get(key string) (string, *Error) {
	result := client.Conn.Get(key)
	val, err := result.Result()
	if err != nil && err != redis.Nil {
		//its error
		return "", &Error{Err: err}
	}
	if err == redis.Nil {
		return "", &Error{Err: NotFoundErr}
	}
	return val, nil
}

func (client Cluster) SetNX(key string, val interface{}, ttl time.Duration) (bool, *Error) {
	result := client.Conn.SetNX(key, val, ttl)
	err := result.Err()
	if err != nil {
		//its error
		return false, &Error{Err: err}
	}
	return result.Val(), nil
}

func (client Cluster) Del(key string) *Error {
	result := client.Conn.Del(key)
	err := result.Err()
	if err != nil && err != redis.Nil {
		//its error
		return &Error{Err: err}
	}
	return nil
}

func (client Cluster) SADD(key string, val ...string) (bool, *Error) {
	result := client.Conn.SAdd(key, val...)
	_, err := result.Result()
	if err != nil {
		//its error
		return false, &Error{Err: err}
	}
	return true, nil
}

func (client Cluster) SREM(key string, member ...string) (bool, *Error) {
	result := client.Conn.SRem(key, member...)
	_, err := result.Result()
	if err != nil && err != redis.Nil {
		//its error
		return false, &Error{Err: err}
	}
	if err == redis.Nil {
		return true, &Error{Err: NotFoundErr}
	}
	return true, nil
}

func (client Cluster) SMembers(key string) ([]string, *Error) {
	result := client.Conn.SMembers(key)
	strarr, err := result.Result()
	if err != nil {
		//its error
		return nil, &Error{Err: err}
	}
	return strarr, nil
}

func (client Cluster) EVAL(script string, keys []string, args []string) (interface{}, error) {
	result := client.Conn.Eval(script, keys, args)
	err := result.Err()
	if err != nil {
		return "", &Error{Err: err}
	}
	return result.Val(), nil
}

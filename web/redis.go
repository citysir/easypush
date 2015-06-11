package main

import (
	"github.com/xuyu/goredis"
	"time"
)

func CreateRedisClient(addr string, db int, password string, timeout, poolSize int) (*goredis.Redis, error) {
	config := &goredis.DialConfig{
		Network:  "tcp",
		Address:  addr,
		Database: db,
		Password: password,
		Timeout:  10 * time.Second,
		MaxIdle:  poolSize,
	}
	return goredis.Dial(config)
}

func CloseRedisClient(redis *goredis.Redis) {
	if redis != nil {
		redis.ClosePool()
	}
}

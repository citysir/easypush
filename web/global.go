package main

import (
	"fmt"
	"github.com/xuyu/goredis"
	"sync"
)

type GlobalContext struct {
	lock  *sync.Mutex
	Redis *goredis.Redis
}

var Global = &GlobalContext{lock: new(sync.Mutex)}

func InitGlobal() {
	Global.lock.Lock()
	defer Global.lock.Unlock()
	initRedis()
}

func initRedis() {
	redis, err := CreateRedisClient(Conf.RedisAddr, Conf.RedisDb, Conf.RedisPassword, Conf.RedisTimeout, Conf.RedisPoolSize)
	if err != nil {
		panic(fmt.Sprintf("failed CreateRedisClient %s, %v", Conf.RedisAddr, err))
	}
	Global.Redis = redis
}

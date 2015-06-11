package main

type Config struct {
	WebBind  string
	PushBind string
	StatBind string
	PerfBind string

	RpcBind string

	NodeRpcPort int

	NodeHostMap map[int]string

	NodeSlotMap map[int][]int //SlotCount=4096

	TokenTimeout   int
	MessageTimeout int

	RedisAddr     string
	RedisDb       int
	RedisPassword string
	RedisTimeout  int
	RedisPoolSize int

	MinClientVersion uint64

	AccessKeys map[string]string
}

var Conf = &Config{
	WebBind:  ":80",
	PushBind: ":8081",
	StatBind: ":8082",
	PerfBind: ":8083",

	RpcBind: ":8099",

	NodeRpcPort: 9099,

	NodeHostMap: map[int]string{1: "127.0.0.1"},

	NodeSlotMap: map[int][]int{1: []int{1, 4096}},

	TokenTimeout:   60 * 60 * 24 * 7,
	MessageTimeout: 60 * 60 * 24 * 7,

	RedisAddr:     "10.20.216.113:6379",
	RedisDb:       0,
	RedisPassword: "",
	RedisTimeout:  10,
	RedisPoolSize: 500,

	MinClientVersion: uint64(201506031950),

	AccessKeys: map[string]string{"bRVYoA5Y70m3dHYk": "rP#S,9O4kl]GwjOD"},
}

package main

type Config struct {
	ChannelBucketCount       int
	ChannelMessageBufferSize int

	ZookeeperNodeId int
	NodeSlotMap     map[int][]int //SlotCount=4096

	NodeBind string
	StatBind string
	PerfBind string

	RpcBind string

	WebRpcAddr string
}

var Conf = &Config{
	ChannelBucketCount:       256,
	ChannelMessageBufferSize: 64,

	ZookeeperNodeId: 1,
	NodeSlotMap:     map[int][]int{1: []int{1, 4096}},

	NodeBind: "127.0.0.1:8080",
	StatBind: "127.0.0.1:9091",
	PerfBind: "127.0.0.1:9093",

	RpcBind: "127.0.0.1:9099",

	WebRpcAddr: "127.0.0.1:8099",
}

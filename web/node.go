package main

import (
	"github.com/citysir/easypush/hash"
)

var nodeHash = hash.NewNodeHash(Conf.NodeSlotMap)

func FindNodeHost(key string) string {
	nodeId := nodeHash.Hash(key)
	return Conf.NodeHostMap[nodeId]
}

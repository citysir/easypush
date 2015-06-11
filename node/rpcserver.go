package main

import (
	log "github.com/log4go"
	"net"
	"net/rpc"
)

type RPC int

func (t *RPC) PushMessages(pushMessageArgs *PushMessageArgs, failedKeys *([]string)) error {
	onPushMessages(pushMessageArgs.Keys, pushMessageArgs.Message, func(failedKey string) {
		*failedKeys = append(*failedKeys, failedKey)
	})
	return nil
}

func rpcListen(bind string) {
	rpcServer := rpc.NewServer()
	rpcServer.Register(new(RPC))
	rpcServer.HandleHTTP("/foo", "/bar")

	l, err := net.Listen("tcp", bind)
	if err != nil {
		log.Error("net.Listen(\"%s\") error(%v)", bind, err)
		panic(err)
	}
	rpcServer.Accept(l)
}

func StartRpc() {
	log.Info("Start rpc listen addr: %s", Conf.RpcBind)
	go rpcListen(Conf.RpcBind)
}

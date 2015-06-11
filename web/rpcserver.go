package main

import (
	"errors"
	log "github.com/log4go"
	"net"
	"net/rpc"
)

type RPC int

func (t *RPC) CheckTokenData(arg *TokenDataArg, valid *bool) error {
	tokenData, err := GetTokenData(arg.Token)
	if err != nil {
		*valid = false
		return err
	}
	if tokenData == nil {
		*valid = false
		return errors.New("token invalid or timeout")
	}
	*valid = (arg.User == tokenData.User)
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
	log.Info("Start rpc listen addr %s", Conf.RpcBind)
	go rpcListen(Conf.RpcBind)
}

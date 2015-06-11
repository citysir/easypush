package main

import (
	"net"
	"net/rpc"
)

type WebRpcClient struct {
	conn   net.Conn
	client *rpc.Client
}

func NewWebRpcClient(rpcAddr string) (*WebRpcClient, error) {
	client := new(WebRpcClient)
	err := client.init(rpcAddr)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (n *WebRpcClient) init(rpcAddr string) error {
	address, err := net.ResolveTCPAddr("tcp", rpcAddr)
	if err != nil {
		return err
	}
	n.conn, err = net.DialTCP("tcp", nil, address)
	if err != nil {
		return err
	}
	n.client = rpc.NewClient(n.conn)
	return nil
}

func (n *WebRpcClient) Close() {
	n.conn.Close()
	n.client.Close()
}

func (n *WebRpcClient) CallCheckTokenData(arg *TokenDataArg) (bool, error) {
	valid := false
	err := n.client.Call("RPC.CheckTokenData", arg, &valid)
	if err != nil {
		return false, err
	}
	return valid, nil
}

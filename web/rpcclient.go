package main

import (
	"net"
	"net/rpc"
)

type NodeRpcClient struct {
	conn   net.Conn
	client *rpc.Client
}

func NewNodeRpcClient(rpcAddr string) (*NodeRpcClient, error) {
	client := new(NodeRpcClient)
	err := client.init(rpcAddr)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (n *NodeRpcClient) init(rpcAddr string) error {
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

func (n *NodeRpcClient) Close() {
	n.conn.Close()
	n.client.Close()
}

func (n *NodeRpcClient) CallPushMessages(args *PushMessageArgs) ([]string, error) {
	failedKeys := []string{}
	err := n.client.Call("RPC.PushMessages", args, &failedKeys)
	if err != nil {
		return failedKeys, err
	}
	return failedKeys, nil
}

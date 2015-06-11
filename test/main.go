package main

import (
	"log"
	"net"
	"net/rpc"
	"time"
)

type Args struct {
	A, B int
}

type Arith int

func (t *Arith) Multiply(args *Args, reply *([]string)) error {
	*reply = append(*reply, "test")
	return nil
}

func (t *Arith) Abc(args *Args, reply *int) error {
	*reply = args.A + args.B
	return nil
}

func main() {
	newServer := rpc.NewServer()
	newServer.Register(new(Arith))

	l, e := net.Listen("tcp", "127.0.0.1:1234") // any available address
	if e != nil {
		log.Fatalf("net.Listen tcp :0: %v", e)
	}

	go newServer.Accept(l)
	newServer.HandleHTTP("/foo", "/bar")
	time.Sleep(2 * time.Second)

	address, err := net.ResolveTCPAddr("tcp", "127.0.0.1:1234")
	if err != nil {
		panic(err)
	}
	conn, _ := net.DialTCP("tcp", nil, address)
	defer conn.Close()

	client := rpc.NewClient(conn)
	defer client.Close()

	args := &Args{7, 8}
	var reply int
	err = client.Call("Arith.Abc", &args, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	log.Println(reply)
}

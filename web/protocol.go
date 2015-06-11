package main

import (
	"encoding/json"
)

const (
	MESSAGE_TYPE_PUSH = 1
	MESSAGE_TYPE_ACK  = 2
)

type Message struct {
	Type int16            `json:"t"`
	Id   int64            `json:"id"`
	Data *json.RawMessage `json:"data"`
}

type PushMessageArgs struct {
	Keys    []string
	Message *Message
}

type TokenDataArg struct {
	Token string
	User  string
	Did   string
}

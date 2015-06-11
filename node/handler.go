package main

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	log "github.com/log4go"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  512,
	WriteBufferSize: 1024,
}

func validate(w http.ResponseWriter, r *http.Request) (string, string, error) {
	token := r.FormValue("token")
	if token == "" {
		return "", "", errors.New("empty token")
	}

	user := r.FormValue("user")
	if user == "" {
		return "", "", errors.New("empty user")
	}

	did := r.FormValue("did")
	if did == "" {
		return "", "", errors.New("empty did")
	}

	log.Debug("token=%s, user=%s, did=%s", token, user, did)

	rpcClient, err := NewWebRpcClient(Conf.WebRpcAddr)
	if err != nil {
		return "", "", err
	}

	arg := &TokenDataArg{Token: token, User: user, Did: did}
	valid, err := rpcClient.CallCheckTokenData(arg)
	if err != nil {
		return "", "", err
	}

	if !valid {
		return "", "", fmt.Errorf("invalid token %s, %s, %s", token, user, did)
	}

	log.Debug("valid token %s %s %s", token, user, did)
	return user, did, nil
}

// serverWs handles websocket requests from the peer.
func serveWs(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	key, did, err := validate(w, r)
	if err != nil {
		http.Error(w, "Failed to validate", 405)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err)
		return
	}
	channel := &Channel{did: did, output: make(chan []byte, Conf.ChannelMessageBufferSize), ws: ws, messages: NewMessageList()}
	hub.register <- &RegisterInfo{key: key, channel: channel}
	go channel.writePump()
	channel.readPump(key)
}

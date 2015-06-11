package main

import (
	log "github.com/log4go"
)

// Hub maintains the set of active conns and broadcasts messages to the
// conns.
type Hub struct {
	// Registered conns.
	channelList *ChannelList

	// Register requests from the conns.
	register chan *RegisterInfo

	// Unregister requests from conns.
	unregister chan *RegisterInfo
}

type RegisterInfo struct {
	key     string
	channel *Channel
}

func (h *Hub) run() {
	for {
		select {
		case c := <-h.register:
			h.channelList.Add(c.key, c.channel)
			log.Debug("Registered %s, %s", c.key, c.channel.did)
			// 可以记录到redis中
		case c := <-h.unregister:
			h.channelList.Remove(c.key)
			log.Debug("UnRegistered %s, %s", c.key, c.channel.did)
			// 可以删除到redis记录
		}
	}
}

func onPushMessage(key string, message *Message) bool {
	channel := hub.channelList.Get(key)
	if channel != nil {
		channel.Push(message)
		return true
	}
	return false
}

func onPushMessages(keys []string, message *Message, f func(string)) {
	for _, key := range keys {
		if !onPushMessage(key, message) {
			f(key)
		}
	}
}

var hub = Hub{
	channelList: NewChannelList(),
	register:    make(chan *RegisterInfo),
	unregister:  make(chan *RegisterInfo),
}

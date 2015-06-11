package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	log "github.com/log4go"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 30 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 1024
)

// Conn is an middleman between the websocket Conn and the hub.
type Channel struct {
	// 连接设备Id
	did string
	// The websocket Conn.
	ws       *websocket.Conn
	messages *MessageList
	// Buffered channel of outbound messages.
	output chan []byte
}

func (c *Channel) Push(message *Message) {
	c.messages.PushBack(message)
	if c.messages.Len() == 1 { //没有消息积压
		data, _ := json.Marshal(message)
		c.output <- data
	} else if c.messages.Len() == Conf.ChannelMessageBufferSize {
		c.Close()
	}
}

func (c *Channel) Len() int {
	return c.messages.Len()
}

func (c *Channel) Close() error {
	if c.ws != nil {
		close(c.output)
		c.messages = nil
		return c.ws.Close()
	}
	return nil
}

// readPump pumps messages from the websocket Conn to the hub.
func (c *Channel) readPump(key string) {
	defer func() {
		hub.unregister <- &RegisterInfo{key, c}
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, data, err := c.ws.ReadMessage()
		if err != nil {
			break
		}

		message := Message{}
		err = json.Unmarshal(data, &message)
		if err != nil {
			break
		}
		c.handleInputMessage(&message)
	}
}

func (c *Channel) handleInputMessage(message *Message) {
	log.Debug("handleInputMessage type=%d, id=%d", message.Type, message.Id)
	if MESSAGE_TYPE_ACK == message.Type { //客户端要回一个ACK确认已经收到消息
		front := c.messages.Front()
		if front != nil && front.Value.Id == message.Id {
			c.messages.Remove(front)
		}
		if c.messages.Len() > 0 { //还有消息积压
			data, _ := json.Marshal(c.messages.Front().Value)
			c.output <- data
		}
	}
}

func (c *Channel) write(messageType int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(messageType, payload)
}

func (c *Channel) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case data, ok := <-c.output:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, data); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

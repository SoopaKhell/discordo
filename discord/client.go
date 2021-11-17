package discord

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
)

const gatewayURL = "wss://gateway.discord.gg/?v=9&encoding=json"

type Client struct {
	wsConn *websocket.Conn
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Connect() error {
	var err error
	c.wsConn, _, err = websocket.DefaultDialer.Dial(gatewayURL, nil)
	if err != nil {
		return err
	}

	var p DiscordGatewayPayload
	if err = c.wsConn.ReadJSON(&p); err != nil {
		return err
	}

	c.onMessage(p)

	return nil
}

func (c *Client) onMessage(p DiscordGatewayPayload) error {
	switch p.Opcode {
	case 10: // Hello
		var h DiscordHello
		if err := json.Unmarshal(p.EventData, &h); err != nil {
			return err
		}

		go c.sendHeartbeat(h.HeartbeatInterval)
	}

	return nil
}

func (c *Client) sendHeartbeat(heartbeatInterval time.Duration) error {
	t := time.NewTicker(heartbeatInterval * time.Millisecond)

	for range t.C {
		if err := c.wsConn.WriteJSON(struct{ op int }{1}); err != nil {
			return err
		}
	}

	return nil
}

package discord

import (
	"encoding/json"
	"time"
)

type DiscordGatewayPayload struct {
	Opcode         int             `json:"op"`
	EventName      string          `json:"t"`
	EventData      json.RawMessage `json:"d"`
	SequenceNumber int64           `json:"s"`
}

type DiscordHello struct {
	HeartbeatInterval time.Duration `json:"heartbeat_interval"`
}

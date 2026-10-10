package game

// GameEvent represents something that happened in the game
type GameEvent struct {
	Type   string         `json:"type"`
	Data   map[string]any `json:"data"`
	Player int            `json:"player"` // which player this is relevant to, -1 for both
}

// ActionMessage is a player action received via WebSocket
type ActionMessage struct {
	Action    string         `json:"action"`
	Data      map[string]any `json:"data"`
	RequestID string         `json:"request_id,omitempty"`
}

// EventCallback is called when events occur (to send to clients)
type EventCallback func(event GameEvent, targetPlayer int)

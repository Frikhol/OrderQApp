package types

import "encoding/json"

type WSMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type OrderUpdatePayload struct {
	OrderID string                 `json:"order_id"`
	Fields  map[string]interface{} `json:"fields"`
}

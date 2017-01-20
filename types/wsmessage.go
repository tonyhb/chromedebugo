package types

type WSMessage struct {
	Method string                 `json:"method"`
	Params map[string]interface{} `json:"params"`
}

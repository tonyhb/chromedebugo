package chromedebugo

import (
	"encoding/json"
	"fmt"
)

type Command struct {
	Method string                 `json:"method"`
	Params map[string]interface{} `json:"params"`
}

type Result struct {
	ID     int                    `json:"id"`
	Result map[string]interface{} `json:"result"`
}

type Error struct {
	ErrorDetail ErrorDetail `json:"error"`
	ID          int         `json:"id"`
}

func (e Error) Error() string {
	return fmt.Sprintf(
		"request '%d' failed with code '%d': %s",
		e.ID,
		e.ErrorDetail.Code,
		e.ErrorDetail.Message,
	)
}

type ErrorDetail struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type commandWrapper struct {
	ID int `json:"id"`
	Command
}

func (c commandWrapper) MarshalJSON() ([]byte, error) {
	data := map[string]interface{}{
		"id":     c.ID,
		"method": c.Command.Method,
		"params": c.Command.Params,
	}
	return json.Marshal(data)
}

type Version struct {
	Browser         string
	ProtocolVersion string `json:"Protocol-Version"`
	UserAgent       string `json:"User-Agent"`
	V8Version       string `json:"V8-Version"`
	WebkitVersion   string `json:"Webkit-Version"`
}

type Info struct {
	Description          string `json:"description"`
	DevtoolsFrontendURL  string `json:"devtoolsFrontendUrl"`
	Id                   string `json:"id"` // todo: UUID package
	Title                string `json:"title"`
	Type                 string `json:"type"`
	URL                  string `json:"url"`
	WebsocketDebuggerURL string `json:"websocketDebuggerURL"`
}

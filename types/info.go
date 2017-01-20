package types

import "net/url"

type Info struct {
	Description          string  `json:"description"`
	DevtoolsFrontendURL  url.URL `json:"devtoolsFrontendUrl"`
	Id                   string  `json:"id"` // todo: UUID package
	Title                string  `json:"title"`
	Type                 string  `json:"type"`
	URL                  url.URL `json:"url"`
	WebsocketDebuggerURL url.URL `json:"websocketDebuggerURL"`
}

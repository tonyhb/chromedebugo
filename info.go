package chromedebugo

type Info struct {
	Description          string `json:"description"`
	DevtoolsFrontendURL  string `json:"devtoolsFrontendUrl"`
	Id                   string `json:"id"` // todo: UUID package
	Title                string `json:"title"`
	Type                 string `json:"type"`
	URL                  string `json:"url"`
	WebsocketDebuggerURL string `json:"websocketDebuggerURL"`
}

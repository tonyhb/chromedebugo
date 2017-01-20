package types

type Version struct {
	Browser         string
	ProtocolVersion string `json:"Protocol-Version"`
	UserAgent       string `json:"User-Agent"`
	V8Version       string `json:"V8-Version"`
	WebkitVersion   string `json:"Webkit-Version"`
}

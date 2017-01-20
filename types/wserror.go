package types

type WSError struct {
	Error DebuggerError `json:"error"`
	ID    int           `json:"id"`
}

type DebuggerError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

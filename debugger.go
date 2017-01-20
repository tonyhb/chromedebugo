package chromedebugo

import (
	"encoding/json"
	"net/http"
	"path"

	"github.com/gorilla/websocket"
	"github.com/tonyhb/chromedebugo/types"
)

type chromeDebugger struct {
	ws   *websocket.Conn
	host string
}

func New(host string) (Debugger, error) {
	// First we must get the websocket URL of the u

	conn, resp, err := new(websocket.Dialer).Dial(host, nil)
	if err != nil {
		return nil, err
	}

	return chromeDebugger{
		host: host,
		ws:   conn,
	}, nil
}

func (cd chromeDebugger) Version() (types.Version, error) {
	return version(cd.host)
}

func (cd chromeDebugger) Info() ([]types.Info, error) {
	return info(cd.host)
}

func version(host string) (types.Version, error) {
	vers := types.Version{}
	data, err := get(path.Join(host, "version"), &vers)
	if err != nil {
		return vers, err
	}
	return data.(types.Version), nil
}

func info(host string) ([]types.Info, error) {
	info := []types.Info{}
	data, err := get(path.Join(host, "version"), &info)
	if err != nil {
		return info, err
	}
	return data.([]types.Info), nil
}

func get(path string, data interface{}) (interface{}, error) {
	resp, err := http.Get(path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(data)
	return data, err
}

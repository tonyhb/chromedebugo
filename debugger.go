package chromedebugo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type debugger struct {
	conn *websocket.Conn
	host string

	errChan chan Error
	resChan chan Result
	cmdChan chan Command

	// id incremnets with each command we send to the debugger
	id   int
	lock sync.Mutex
}

func newBaseDebugger(host string) (*debugger, error) {
	// First we must get the websocket URL of the host
	info, err := info(host)
	if err != nil {
		return nil, err
	}
	if len(info) != 1 {
		return nil, fmt.Errorf("error getting chrome info: unexpected number of responses from /info")
	}

	conn, _, err := new(websocket.Dialer).Dial(info[0].WebsocketDebuggerURL, nil)
	if err != nil {
		return nil, err
	}

	// Remote debugging is async, and there are three classes of messages
	// that can be sent back.  They are:
	// - Errors, from failed commands sent to the debugger
	// - Results, from successful commands sent to the debugger
	// - Commands, wihch notify clients of commands created by the remote
	//   debugger
	errChan := make(chan Error)
	cmdChan := make(chan Command)
	resChan := make(chan Result)

	return &debugger{
		host: host,
		conn: conn,

		errChan: errChan,
		cmdChan: cmdChan,
		resChan: resChan,

		id:   1,
		lock: sync.Mutex{},
	}, nil
}

func version(host string) (Version, error) {
	resp, err := http.Get(host + "/json/version")
	if err != nil {
		return Version{}, err
	}
	defer resp.Body.Close()
	data := Version{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	return data, err
}

func info(host string) ([]Info, error) {
	resp, err := http.Get(host + "/json/list")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data := []Info{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	return data, err
}

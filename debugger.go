package chromedebugo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type chromeDebugger struct {
	conn *websocket.Conn
	host string

	errChan chan Error
	resChan chan Result
	cmdChan chan Command

	// id incremnets with each command we send to the debugger
	id   int
	lock sync.Mutex
}

func New(host string) (Debugger, error) {
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
	//
	// You can listen on any or all of these, or ignore them entirely.
	errChan := make(chan Error)
	cmdChan := make(chan Command)
	resChan := make(chan Result)

	go func() {
		for {
			_, data, err := conn.ReadMessage()
			if err != nil {
				return
			}
			resp, err := decodeResponse(data)
			if err != nil {
				return
			}
			switch resp.(type) {
			case Error:
				errChan <- resp.(Error)
			case Result:
				resChan <- resp.(Result)
			case Command:
				cmdChan <- resp.(Command)
			}
		}
	}()

	return &chromeDebugger{
		host: host,
		conn: conn,

		errChan: errChan,
		cmdChan: cmdChan,
		resChan: resChan,

		id:   1,
		lock: sync.Mutex{},
	}, nil
}

// Version returns the chrome version inforamation from /json/version
func (cd chromeDebugger) Version() (Version, error) {
	return version(cd.host)
}

// Info returns a slice of browser contexts from /json/list
func (cd chromeDebugger) Info() ([]Info, error) {
	return info(cd.host)
}

func (cd *chromeDebugger) Send(msg Command) (int, error) {
	cd.lock.Lock()
	defer cd.lock.Unlock()
	defer func() {
		cd.id++
	}()

	wrapper := commandWrapper{
		ID:      cd.id,
		Command: msg,
	}

	if err := cd.conn.WriteJSON(wrapper); err != nil {
		return 0, fmt.Errorf("error sending command to chrome: %s", err)
	}

	return wrapper.ID, nil
}

func (cd chromeDebugger) ErrorChan() chan Error {
	return cd.errChan
}

func (cd chromeDebugger) ResultChan() chan Result {
	return cd.resChan
}

func (cd chromeDebugger) CommandChan() chan Command {
	return cd.cmdChan
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

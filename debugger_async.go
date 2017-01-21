package chromedebugo

import "fmt"

type asyncDebugger struct {
	*debugger
}

func NewAsync(host string) (*asyncDebugger, error) {
	base, err := newBaseDebugger(host)
	if err != nil {
		return nil, err
	}

	// Asynchronous debugging is simple: set up a goroutine which listens
	// for all incoming messages from the websocket connection and dispatch
	// them to the relevant channels.
	//
	// We never block for incoming calls and only communicate this way.
	go func() {
		for {
			_, data, err := base.conn.ReadMessage()
			if err != nil {
				return
			}
			resp, err := decodeResponse(data)
			if err != nil {
				return
			}
			switch resp.(type) {
			case Error:
				base.errChan <- resp.(Error)
			case Result:
				base.resChan <- resp.(Result)
			case Command:
				base.cmdChan <- resp.(Command)
			}
		}
	}()

	return &asyncDebugger{base}, nil
}

// Version returns the chrome version inforamation from /json/version
func (ad asyncDebugger) Version() (Version, error) {
	return version(ad.host)
}

// Info returns a slice of browser contexts from /json/list
func (ad asyncDebugger) Info() ([]Info, error) {
	return info(ad.host)
}

func (ad *asyncDebugger) Send(cmd Command) (int, error) {
	ad.lock.Lock()
	defer ad.lock.Unlock()

	wrapper := commandWrapper{
		ID:      ad.id,
		Command: cmd,
	}

	if err := ad.conn.WriteJSON(wrapper); err != nil {
		return 0, fmt.Errorf("error sending command to chrome: %s", err)
	}

	ad.id++
	return wrapper.ID, nil
}

func (ad asyncDebugger) ErrorChan() chan Error {
	return ad.errChan
}

func (ad asyncDebugger) ResultChan() chan Result {
	return ad.resChan
}

func (ad asyncDebugger) CommandChan() chan Command {
	return ad.cmdChan
}

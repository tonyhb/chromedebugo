package chromedebugo

import (
	"fmt"
	"sync"
)

type syncDebugger struct {
	*debugger
	// lock ensures that we can only send or batch on one goroutine at a
	// time
	lock *sync.Mutex

	// outstanding is a waitgroup which stores the remaining number of
	// commands sent without a response
	outstanding *sync.WaitGroup

	// responses stores a map of all command responses keyed by their ID
	responses map[int]interface{}
}

func NewSync(host string) (*syncDebugger, error) {
	base, err := newBaseDebugger(host)
	if err != nil {
		return nil, err
	}

	debugger := &syncDebugger{
		debugger:    base,
		lock:        &sync.Mutex{},
		outstanding: &sync.WaitGroup{},
		responses:   map[int]interface{}{},
	}

	go func() {
		for {
			_, data, err := debugger.conn.ReadMessage()
			if err != nil {
				return
			}
			resp, err := decodeResponse(data)
			if err != nil {
				return
			}

			switch resp.(type) {
			case Error:
				debugger.responses[resp.(Error).ID] = resp
				debugger.outstanding.Done()
				base.errChan <- resp.(Error)
			case Result:
				debugger.responses[resp.(Result).ID] = resp
				debugger.outstanding.Done()
				base.resChan <- resp.(Result)
			case Command:
				base.cmdChan <- resp.(Command)
			}
		}
	}()

	return debugger, nil
}

// Version returns the chrome version inforamation from /json/version
func (sd syncDebugger) Version() (Version, error) {
	return version(sd.host)
}

// Info returns a slice of browser contexts from /json/list
func (sd syncDebugger) Info() ([]Info, error) {
	return info(sd.host)
}

func (sd syncDebugger) Send(cmd Command) (interface{}, error) {
	sd.lock.Lock()
	defer sd.lock.Unlock()

	wrapper := commandWrapper{
		ID:      sd.id,
		Command: cmd,
	}
	if err := sd.conn.WriteJSON(wrapper); err != nil {
		return nil, fmt.Errorf("error sending command to chrome: %s", err)
	}

	sd.id++

	sd.outstanding.Add(1)
	sd.outstanding.Wait()
	return sd.responses[wrapper.ID], nil
}

func (sd syncDebugger) Batch(commands []Command) ([]interface{}, error) {
	sd.lock.Lock()
	defer sd.lock.Unlock()

	startId := sd.id

	for _, cmd := range commands {
		wrapper := commandWrapper{
			ID:      sd.id,
			Command: cmd,
		}
		if err := sd.conn.WriteJSON(wrapper); err != nil {
			return nil, fmt.Errorf("error sending command to chrome: %s", err)
		}
		sd.id++
		sd.outstanding.Add(1)
	}

	sd.outstanding.Wait()

	responses := make([]interface{}, len(commands), len(commands))
	for i := 0; i < len(commands); i++ {
		idx := startId + i
		responses[i] = sd.responses[idx]
	}

	return responses, nil
}

func (sd syncDebugger) ErrorChan() chan Error {
	return sd.errChan
}

func (sd syncDebugger) ResultChan() chan Result {
	return sd.resChan
}

func (sd syncDebugger) CommandChan() chan Command {
	return sd.cmdChan
}

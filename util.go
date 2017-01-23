package chromedebugo

import (
	"encoding/json"
	"fmt"
)

func decodeResponse(data []byte, commands map[int]Command) (interface{}, error) {
	// data can be a JSON marshaled Result, Error
	// or Command.
	root := map[string]interface{}{}
	if err := json.Unmarshal(data, &root); err != nil {
		return nil, err
	}

	if _, ok := root["error"]; ok {
		errMsg := Error{}
		if err := json.Unmarshal(data, &errMsg); err != nil {
			return nil, err
		}

		cmd, ok := commands[errMsg.ID]
		if ok {
			errMsg.Request = &cmd
		}

		return errMsg, nil
	}

	if _, ok := root["id"]; ok {
		result := Result{}
		if err := json.Unmarshal(data, &result); err != nil {
			return nil, err
		}

		cmd, ok := commands[result.ID]
		if ok {
			result.Request = &cmd
		}

		return result, nil
	}

	if _, ok := root["method"]; ok {
		cmd := Command{}
		if err := json.Unmarshal(data, &cmd); err != nil {
			return nil, err
		}
		return cmd, nil
	}

	return nil, fmt.Errorf("unknown response: %s", data)
}

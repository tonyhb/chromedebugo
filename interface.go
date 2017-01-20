package chromedebugo

import "github.com/tonyhb/chromedebugo/types"

type Debugger interface {
	Version() (types.Version, error)
	Info() ([]types.Info, error)
}

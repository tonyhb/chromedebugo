package chromedebugo

type Debugger interface {
	Version() (types.Version, error)
	Info() ([]types.Info, error)
}

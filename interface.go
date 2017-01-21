package chromedebugo

type Debugger interface {
	Version() (Version, error)
	Info() ([]Info, error)

	Send(Command) (int, error)

	ErrorChan() chan Error
	ResultChan() chan Result
	CommandChan() chan Command
}

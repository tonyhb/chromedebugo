package chromedebugo

type AsyncDebugger interface {
	Version() (Version, error)
	Info() ([]Info, error)

	Send(Command) (int, error)

	ErrorChan() chan Error
	ResultChan() chan Result
	CommandChan() chan Command
}

type SyncDebugger interface {
	Version() (Version, error)
	Info() ([]Info, error)

	// Send dispatches a command to headless chrome and blocks until chrome
	// sends us a Result or Error
	Send(Command) (Result, error)

	// Batch dispatches mutiple commands in order to chrome and blocks until
	// all responses for commands have been receivevd.
	//
	// Note that the responses back are held in a slice of interfaces; they
	// may be either a Result or Error type.
	//
	// The err type returned will be non-nil if there was an error sending
	// any of the batched commands to chrome.
	Batch([]Command) ([]interface{}, error)

	ErrorChan() chan Error
	ResultChan() chan Result
	CommandChan() chan Command
}

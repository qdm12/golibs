package command

//go:generate mockgen -destination=mock_$GOPACKAGE/$GOFILE . Commander

// Commander contains methods to run and start shell commands.
type Commander interface {
	// Run runs a command in a blocking manner, returning its output and an error if it failed
	Run(cmd Cmd) (output string, err error)
	// Start launches a command and streams stdout and stderr to channels.
	// All the channels returned should be closed when an error,
	// nil or not, is received in the waitError channel.
	// The channels should NOT be closed if an error is returned directly
	// with err, as they will already be closed internally by the function.
	Start(cmd Cmd) (stdoutLines, stderrLines chan string,
		waitError chan error, err error)
}

type commander struct{}

func NewCommander() Commander {
	return &commander{}
}

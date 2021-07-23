package command

var _ RunStarter = (*Cmder)(nil)

//go:generate mockgen -destination=mock_$GOPACKAGE/$GOFILE . RunStarter

type RunStarter interface {
	Runner
	Starter
}

// Cmder handles running subprograms synchronously and asynchronously.
type Cmder struct{}

func NewCmder() *Cmder {
	return &Cmder{}
}

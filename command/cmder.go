package command

type RunStarter interface {
	Runner
	Starter
}

// Cmder handles running subprograms synchronously and asynchronously.
type Cmder struct{}

func NewCmder() *Cmder {
	return &Cmder{}
}

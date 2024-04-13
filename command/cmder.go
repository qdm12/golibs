package command

// Cmder handles running subprograms synchronously and asynchronously.
type Cmder struct{}

func NewCmder() *Cmder {
	return &Cmder{}
}

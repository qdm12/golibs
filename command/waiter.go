package command

import (
	"sync"
)

// Waiter waits for multiple wait functions to finish
//go:generate mockgen -destination=mock_command/waiter.go . Waiter
type Waiter interface {
	// Adds a wait function to the waiter in its own blocking goroutine
	Add(waitFunction func() error)
	// WaitForAll waits for all wait functions to complete and return any eventual error from them
	WaitForAll() (errors []error)
}

type waiter struct {
	n      uint
	errors chan error
	sync.Mutex
}

func NewWaiter() Waiter {
	return &waiter{
		errors: make(chan error),
	}
}

func (w *waiter) Add(waitFunction func() error) {
	w.Lock()
	w.n++
	w.Unlock()
	go func() {
		w.errors <- waitFunction()
	}()
}

func (w *waiter) WaitForAll() (errors []error) {
	w.Lock()
	for w.n > 0 {
		w.Unlock()
		if err := <-w.errors; err != nil {
			errors = append(errors, err)
		}
		w.Lock()
		w.n--
	}
	w.Unlock()
	return errors
}

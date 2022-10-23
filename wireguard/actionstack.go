package kernelspace

import (
	"errors"
	"sync"
)

var (
	// ErrAlreadyRun error is throwed with panic on Push or Run invocation
	// after Run was invoked at least once
	ErrAlreadyRun = errors.New("actions already fired")
)

// Action represents stackable action.
type Action func()

// ActionStack is a stack of actions which are executed in reverse order
// they were added
type ActionStack struct {
	stack []Action
	fired bool
	mu    sync.Mutex
}

// NewActionStack creates empty ActionStack
func NewActionStack() *ActionStack {
	return new(ActionStack)
}

// Run executes pushed actions in reverse order (FILO)
func (a *ActionStack) Run() {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.fired {
		panic(ErrAlreadyRun)
	}

	for i := len(a.stack) - 1; i >= 0; i-- {
		a.stack[i]()
	}
	a.stack = nil
	a.fired = true
}

// Push adds new action on top of stack. Last added action will be executed
// by Run() first.
func (a *ActionStack) Push(elems ...Action) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.fired {
		panic(ErrAlreadyRun)
	}

	a.stack = append(a.stack, elems...)
}

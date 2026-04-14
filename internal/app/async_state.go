package app

// AsyncState tracks async operation state.
// Accessor methods (End, IsAborted, ClearAborted, IsActive, Abort) enforce
// Tell-Don't-Ask for operationAborted state transitions and provide a
// stable API surface for future cross-package use.
type AsyncState struct {
	operationActive  bool
	operationAborted bool
	exitAllowed      bool
}

func NewAsyncState() *AsyncState {
	return &AsyncState{
		operationActive:  false,
		operationAborted: false,
		exitAllowed:      false,
	}
}

func (as *AsyncState) Start() {
	as.operationActive = true
	as.operationAborted = false
	as.exitAllowed = false
}

func (as *AsyncState) End() {
	as.operationActive = false
}

func (as *AsyncState) Abort() {
	as.operationAborted = true
}

func (as *AsyncState) ClearAborted() {
	as.operationAborted = false
}

func (as *AsyncState) IsActive() bool {
	return as.operationActive
}

func (as *AsyncState) IsAborted() bool {
	return as.operationAborted
}

func (as *AsyncState) CanExit() bool {
	return as.exitAllowed
}

func (as *AsyncState) SetExitAllowed(allowed bool) {
	as.exitAllowed = allowed
}

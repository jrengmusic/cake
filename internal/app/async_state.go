package app

// AsyncState tracks async operation state
type AsyncState struct {
	OperationActive  bool
	OperationAborted bool
	ExitAllowed      bool
}

func NewAsyncState() *AsyncState {
	return &AsyncState{
		OperationActive:  false,
		OperationAborted: false,
		ExitAllowed:      false,
	}
}

func (as *AsyncState) Start() {
	as.OperationActive = true
	as.OperationAborted = false
	as.ExitAllowed = false
}

func (as *AsyncState) End() {
	as.OperationActive = false
}

func (as *AsyncState) Abort() {
	as.OperationAborted = true
}

func (as *AsyncState) ClearAborted() {
	as.OperationAborted = false
}

func (as *AsyncState) IsActive() bool {
	return as.OperationActive
}

func (as *AsyncState) IsAborted() bool {
	return as.OperationAborted
}

func (as *AsyncState) CanExit() bool {
	return as.ExitAllowed
}

func (as *AsyncState) SetExitAllowed(allowed bool) {
	as.ExitAllowed = allowed
}

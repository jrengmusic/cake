package app

// AsyncState tracks async operation state
type AsyncState struct {
	operationActive  bool
	operationAborted bool
}

func NewAsyncState() *AsyncState {
	return &AsyncState{
		operationActive:  false,
		operationAborted: false,
	}
}


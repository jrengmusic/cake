package app

import "time"

const (
	// CacheRefreshInterval is the UI refresh rate during async operations
	CacheRefreshInterval = 100 * time.Millisecond

	// SpinnerTickInterval is the independent spinner animation tick rate (decoupled from console refresh)
	SpinnerTickInterval = 80 * time.Millisecond

	MaxAutoScanInterval   = 60
	MinAutoScanInterval   = 1
	AutoScanIntervalStep  = 10
	DefaultTerminalWidth  = 80
	DefaultTerminalHeight = 24
	MenuBannerSplitPct    = 65
	IdleScanThreshold     = 30 * time.Second
)

package ui

// OpType identifies the async operation currently in progress.
// Lives in ui to avoid circular imports — app imports ui, not the reverse.
type OpType int

const (
	OpNone       OpType = iota
	OpBuild
	OpGenerate
	OpClean
	OpCleanAll
	OpRegenerate
)

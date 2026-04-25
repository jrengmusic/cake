//go:build !windows

package utils

import (
	"fmt"
	"os/exec"
)

// ProcessTree is a no-op on non-Windows platforms. Unix process groups handle
// tree termination via SIGTERM to the group leader, but CAKE currently only
// ships Windows support for this abort path.
type ProcessTree struct{}

// StartProcessTree starts cmd with default stdlib behavior. Tree kill is not
// implemented for non-Windows — context cancellation terminates only the
// direct child.
func StartProcessTree(cmd *exec.Cmd) (*ProcessTree, error) {
	if startErr := cmd.Start(); startErr != nil {
		return nil, fmt.Errorf("cmd.Start failed: %w", startErr)
	}
	return &ProcessTree{}, nil
}

// Close is a no-op on non-Windows.
func (pt *ProcessTree) Close() {}

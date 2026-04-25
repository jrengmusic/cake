//go:build windows

package utils

import (
	"fmt"
	"os/exec"
	"unsafe"

	"golang.org/x/sys/windows"
)

// ProcessTree represents a subprocess and its descendants, bound to a Windows Job Object
// so that closing the job handle terminates the whole tree. Closing the job is the only
// way to guarantee ninja.exe / cl.exe grandchildren are killed when CAKE aborts a build.
type ProcessTree struct {
	jobHandle windows.Handle
}

// StartProcessTree starts cmd, then binds the resulting process (and all future descendants)
// to a newly-created Job Object configured with JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE.
// Callers MUST call Close on the returned ProcessTree when the subprocess has finished
// (normally or via abort) to release the handle and, if still running, terminate the tree.
func StartProcessTree(cmd *exec.Cmd) (*ProcessTree, error) {
	jobHandle, jobErr := windows.CreateJobObject(nil, nil)
	if jobErr != nil {
		return nil, fmt.Errorf("CreateJobObject failed: %w", jobErr)
	}

	info := windows.JOBOBJECT_EXTENDED_LIMIT_INFORMATION{
		BasicLimitInformation: windows.JOBOBJECT_BASIC_LIMIT_INFORMATION{
			LimitFlags: windows.JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE,
		},
	}

	_, setErr := windows.SetInformationJobObject(
		jobHandle,
		windows.JobObjectExtendedLimitInformation,
		uintptr(unsafe.Pointer(&info)),
		uint32(unsafe.Sizeof(info)),
	)
	if setErr != nil {
		// discard: rollback cleanup — original error is what the caller sees
		_ = windows.CloseHandle(jobHandle)
		return nil, fmt.Errorf("SetInformationJobObject failed: %w", setErr)
	}

	if startErr := cmd.Start(); startErr != nil {
		// discard: rollback cleanup — original error is what the caller sees
		_ = windows.CloseHandle(jobHandle)
		return nil, fmt.Errorf("cmd.Start failed: %w", startErr)
	}

	openedHandle, openErr := windows.OpenProcess(
		windows.PROCESS_SET_QUOTA|windows.PROCESS_TERMINATE,
		false,
		uint32(cmd.Process.Pid),
	)
	if openErr != nil {
		// discard: rollback cleanup — original error is what the caller sees
		_ = windows.CloseHandle(jobHandle)
		return nil, fmt.Errorf("OpenProcess failed: %w", openErr)
	}

	if assignErr := windows.AssignProcessToJobObject(jobHandle, openedHandle); assignErr != nil {
		// discard: rollback cleanup — original error is what the caller sees
		_ = windows.CloseHandle(openedHandle)
		// discard: rollback cleanup — original error is what the caller sees
		_ = windows.CloseHandle(jobHandle)
		return nil, fmt.Errorf("AssignProcessToJobObject failed: %w", assignErr)
	}

	// discard: process handle no longer needed — job retains the binding
	_ = windows.CloseHandle(openedHandle)

	return &ProcessTree{jobHandle: jobHandle}, nil
}

// Close releases the job handle. Because the job is configured with
// JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE, closing the last handle also terminates
// every process still in the job. Safe to call multiple times.
func (pt *ProcessTree) Close() {
	if pt.jobHandle != 0 {
		// discard: best-effort teardown — process tree kill has already been commanded
		// by JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE; nothing actionable remains if close fails
		_ = windows.CloseHandle(pt.jobHandle)
		pt.jobHandle = 0
	}
}
